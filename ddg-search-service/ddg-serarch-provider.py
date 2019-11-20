from confluent_kafka import KafkaError, Consumer, Producer
import json
import requests
import itertools
from bs4 import BeautifulSoup
import urllib.parse

kafkaConsumer = Consumer({
    'bootstrap.servers': 'localhost:9092',
    'group.id': 'ddg-search-provider',
    'auto.offset.reset': 'earliest'
})

kafkaProducer = Producer({'bootstrap.servers': 'localhost:9091'})

kafkaConsumer.subscribe(['search'])

def delivery_report(err, msg):
    """ Called once for each message produced to indicate delivery result.
        Triggered by poll() or flush(). """
    if err is not None:
        print('Message delivery failed: {}'.format(err))
    else:
        print('Message delivered to {} [{}]'.format(msg.topic(), msg.partition()))

URL = 'https://duckduckgo.com/html/?q='
header = "{\"authority\": \"duckduckgo.com\", \"method\": \"POST\", \"path\": \"/html/\", \"scheme\": \"https\", \"accept\": \"text/html,application/xml,image/webp,image/apng,*/*\", \"accept-language\": \"en-US,en;q=0.9\",\"cache-control\": \"max-age=0\", \"content-length\": \"28\", \"content-type\": \"application/json\",\"cookie\": \"ax=v195-2\", \"origin\": \"https://duckduckgo.com\", \"referer\": \"https://duckduckgo.com/\", \"sec-fetch-mode\": \"navigate\", \"sec-fetch-site\": \"same-origin\",\"sec-fetch-user\": \"?1\", \"upgrade-insecure-requests\": \"1\", \"user-agent\": \"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.97 Safari/537.36\"}"

print("Listening for new search messages to process...")
while True:
    msg = kafkaConsumer.poll(1.0)

    if msg is None:
        continue
    if msg.error():
        print("Consumer error {}".format(msg.error()))
    # message format sample
    # {"searchword": "Shah%20Rukh%20Khan", "uuid":"9b7c489a-798d-4dee-a3a6-934a90fd8e99"}
    print("New message received: ", msg.value())
    try:
        s = json.loads(msg.value().decode('utf-8'))
    except Exception as e:
        print("Malformed search query, query should format {'searchword':'search string', 'uuid':'uuid string'}")
        continue
    searchkey = s["searchword"]
    uuid = s["uuid"]

    if searchkey is None:
        continue

    params = {'q': searchkey,'kl': 'us-en',}
    
    attempt = 0
    while attempt in range(0,3):
        try:
            print("Running Duck Duck Go search for '{}'".format(urllib.parse.unquote_plus(searchkey)))
            
            r = requests.post(URL + searchkey, data=params, headers=json.loads(header))
            soup = BeautifulSoup(r.content, 'html5lib')

            top10 = {}
            for element in list(itertools.islice(soup.findAll('a', attrs={'class':'result__a'}),10)):
                link = element['href']
                top10[link] = {'title':element.get_text(), 'link': link}

            for element in list(itertools.islice(soup.findAll('a', attrs={'class':'result__snippet'}),10)):
                link = element['href']
                top10[link]['snippet'] = element.get_text()

            msg = {}
            msg['searchengine']="ddg"
            msg['uuid'] = uuid
            msg['keyword'] = searchkey
            msg['results'] = json.dumps(list(top10.values()))

            # Post to Kafka
            en_msg = json.dumps(msg).encode('utf-8')
            print("Posting result to 'search-result' topic")
            kafkaProducer.produce('searchengine-result', en_msg, callback=delivery_report)
            break
        except Exception as e:
            print("Attempt-",attempt,": Exception occured while making DDG query", e.with_traceback())
            break
    # Fetch any available callbacks from previous posts
    kafkaProducer.flush()

kafkaConsumer.close()