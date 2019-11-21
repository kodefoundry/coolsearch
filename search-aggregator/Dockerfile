# Stage 1 Preparing the go builder
FROM golang:alpine as go-builder

RUN apk --no-cache add build-base git gcc pkgconfig

ADD . /go/src/github.com/kodefoundry/coolsearch/search-aggregator
WORKDIR /go/src/github.com/kodefoundry/coolsearch/search-aggregator

# Install librdkafka that is needed for confluentinc/confluent-kafka-go.v1/kafka
RUN apk -U add ca-certificates
RUN apk update && apk upgrade && apk add git bash build-base sudo
RUN git clone https://github.com/edenhill/librdkafka.git && cd librdkafka && ./configure --prefix /usr && make -j && make install
#RUN rm -rf librdkafka
RUN go get ./

RUN go build -o search-aggregator-service

# Stage 2, statically bound app image is ready. We copy it to a alpine image and run it

FROM alpine

WORKDIR /app

# Install librdkafka that is needed for confluentinc/confluent-kafka-go.v1/kafka
RUN apk -U add ca-certificates
RUN apk update && apk upgrade && apk add git bash build-base sudo
RUN git clone https://github.com/edenhill/librdkafka.git && cd librdkafka && ./configure --prefix /usr && make -j && make install
# Keep only the installed binaries and remove the code so that image size is small
RUN rm -rf librdkafka

COPY --from=go-builder /go/src/github.com/kodefoundry/coolsearch/search-aggregator /app

EXPOSE 8090

ENTRYPOINT ./search-aggregator-service