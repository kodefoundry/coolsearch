from confluent_kafka import KafkaError, Consumer, Producer
import json
import requests
import itertools
import urllib.parse

kafkaConsumer = Consumer({
    'bootstrap.servers': 'localhost:9092',
    'group.id': 'googl-search-provider',
    'auto.offset.reset': 'earliest'
})

kafkaProducer = Producer({'bootstrap.servers': 'localhost:9092'})

kafkaConsumer.subscribe(['search'])

def delivery_report(err, msg):
    """ Called once for each message produced to indicate delivery result.
        Triggered by poll() or flush(). """
    if err is not None:
        print('Message delivery failed: {}'.format(err))
    else:
        print('Message delivered to {} [{}]'.format(msg.topic(), msg.partition()))


URL = "https://www.googleapis.com/customsearch/v1?key=AIzaSyC9jWQ_P03Ad_XpavhF3qwFTeG_tZR4p4c&cx=007658124463301196615:cdxn5t6cl2e&q="
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

    while range(0,3):
        try:
            print("Running Google search for '{}'".format(urllib.parse.unquote_plus(searchkey)))
            r = requests.get(url = URL + searchkey)
            data = r.json()
            top10 = list(itertools.islice(data["items"], 10))
            r = list(map(lambda v:{'title':v['title'], 'link':v['link'], 'snippet':v['htmlSnippet']}, top10))
            msg = {}
            msg['searchengine']="google"
            msg['uuid'] = uuid
            msg['keyword'] = searchkey
            msg['results'] = json.dumps(r)
            # Post to Kafka
            en_msg = json.dumps(msg).encode('utf-8')
            print("Posting result to 'search-result' topic")
            kafkaProducer.produce('searchengine-result', en_msg, callback=delivery_report)
            break
        except Exception as e:
            print("Exception occured while making google query", e)
            break
    # Fetch any available callbacks from previous posts
    kafkaProducer.flush()

kafkaConsumer.close()
