from confluent_kafka import KafkaError, Consumer, Producer
import json
import requests
import itertools
import wikipedia
import urllib.parse
import os

try:
    kafkaHost = os.environ["KAFKA_HOST"]
except KeyError as err:
    kafkaHost = "localhost:9092"

c = Consumer({
    'bootstrap.servers': kafkaHost,
    'group.id': 'wiki-search-provider',
    'auto.offset.reset': 'earliest'
})

p = Producer({'bootstrap.servers': kafkaHost})

c.subscribe(['search'])

def delivery_report(err, msg):
    """ Called once for each message produced to indicate delivery result.
        Triggered by poll() or flush(). """
    if err is not None:
        print('Message delivery failed: {}'.format(err))
    else:
        print('Message delivered to {} [{}]'.format(msg.topic(), msg.partition()))

print("Listening for new search messages to process...")
while True:
    msg = c.poll(1.0)

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

    #searchkey = urllib.parse.unquote_plus(searchkey)

    while range(0,3):
        try:
            print("Running Wiki search for '{}'".format(urllib.parse.unquote_plus(searchkey)))
            page = wikipedia.page(searchkey)
            r = {}
            r["title"] = page.title
            r["link"] = page.url
            summary = page.summary
            st = 0
            for i in range(0,2):
                st = summary.find("\n",st,len(summary))
            r["snippet"] = summary[:st]
            msg = {}
            msg['searchengine']="wiki"
            msg['uuid'] = uuid
            msg['keyword'] = searchkey
            msg['results'] = json.dumps([r])
            # Post to Kafka
            en_msg = json.dumps(msg).encode('utf-8')
            print("Posting result to 'search-result' topic")
            p.produce('searchengine-result', en_msg, callback=delivery_report)
            break
        except Exception as e:
            print("Exception occured while making google query", e)
            break

    # Fetch any available callbacks from previous posts
    p.flush()

c.close()
