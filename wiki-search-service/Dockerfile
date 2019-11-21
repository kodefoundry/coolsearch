FROM python:alpine

RUN apk --no-cache add build-base git gcc pkgconfig
# Install librdkafka that is needed for confluentinc/confluent-kafka-go.v1/kafka
RUN apk -U add ca-certificates
RUN apk update && apk upgrade && apk add git bash build-base sudo
RUN git clone https://github.com/edenhill/librdkafka.git && cd librdkafka && ./configure --prefix /usr && make -j && make install
RUN rm -rf librdkafka

COPY . /usr/bin/wiki-search-service
WORKDIR /usr/bin/wiki-search-service

RUN pip install confluent-kafka urllib3 requests wikipedia

CMD [ "python", "wiki-search-provider.py" ]