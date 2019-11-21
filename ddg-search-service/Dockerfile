FROM python:alpine

RUN apk --no-cache add build-base git gcc pkgconfig
# Install librdkafka that is needed for confluentinc/confluent-kafka-go.v1/kafka
RUN apk -U add ca-certificates
RUN apk update && apk upgrade && apk add git bash build-base sudo
RUN git clone https://github.com/edenhill/librdkafka.git && cd librdkafka && ./configure --prefix /usr && make -j && make install
RUN rm -rf librdkafka

COPY . /usr/bin/ddg-search-service
WORKDIR /usr/bin/ddg-search-service

RUN pip install beautifulsoup4 confluent-kafka urllib3 requests

EXPOSE 9091

CMD [ "python", "ddg-search-provider.py" ]