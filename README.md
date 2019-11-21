# coolsearch
A demo application that can aggregate search results from Google, DuckDuckGo and Wikipedia for a given search keyword
### Design 

Design diagram
<img src="https://github.com/kodefoundry/coolsearch/blob/master/docs/coolsearch-design.png" width="600" height="430">

### Screen Shots
Search Page
<img src="https://github.com/kodefoundry/coolsearch/blob/master/docs/Coolsearch.png" width="600" height="430">

Search Listing
<img src="https://github.com/kodefoundry/coolsearch/blob/master/docs/Search%20Listing.png" width="600" height="430">

Search Pending
<img src="https://github.com/kodefoundry/coolsearch/blob/master/docs/Search%20Pending.png" width="600" height="430">

Search Results1
<img src="https://github.com/kodefoundry/coolsearch/blob/master/docs/Search%20Result1.png" width="600" height="430">

Search Results2
<img src="https://github.com/kodefoundry/coolsearch/blob/master/docs/Search%20Result2.png" width="600" height="430">


#### Non docker setup.

Note I have tested this way on a Mac Book Pro with 16Gb RAM.

There are total 9 micro services as listed below
1. coolsearch-api               - golang
2. ddg-search-service           - Python
3. google-search-service        - Python
4. wiki-search-service          - Python
5. search-aggregator            - golang
6. search-persistance           - golang
7. coolsearch-result-api        - golang
8. coolsearch-result-notifier   - Nodejs
9. coolsearch-app-server        - JS, Nodejs, React

##### Dependencies

The following dependencies should be resolved before running the start.sh script

librdkafka  - This is a dependency for confluentinc/confluent-kafka-go.v1/kafka

Run the following to resolve this dependency on host machine. (Should run on Mac and linux)

Run this only on Ubuntu
apk -U add ca-certificates
apk update && apk upgrade && apk add git bash build-base sudo

Linux/Mac
git clone https://github.com/edenhill/librdkafka.git && cd librdkafka && ./configure --prefix /usr && make -j && make install

BeautifulSoup4 - This is a python dependency

Run pip/pip3 install beautifulsoup4

To fix Nodejs dependencies for services 8,9 run 
npm install

Golang Dependencies - To resolve these need to run 'go get' on below.

go get github.com/go-chi/chi
go get gopkg.in/confluentinc/confluent-kafka-go.v1/kafka
go get github.com/go-redis/redis
go get go.mongodb.org/mongo-driver/mongo
go get gopkg.in/mgo.v2/bson
go get github.com/go-chi/cors

##### How to startup all the servics

First startup the infrastructure containers using docker compose by running below command

docker-compose up

Running the Micro services - (Run the above steps first to resolve dependencies)

Run the nodejs services by running 'npm start'
Run the Python services by running 'python3 <python file name>' (Python 3 is required)
Run the Golang services by running 'go run <go file name>'

##### TODO

Create a startup.sh & shutdown.sh file that can startup and shutdown the services.
Complete Dockerization and create docker-compose.yml to run all services together.










