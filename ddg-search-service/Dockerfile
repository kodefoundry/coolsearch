FROM python:alpine

COPY . /usr/bin/ddg-search-service
WORKDIR /usr/bin/ddg-search-service

RUN pip install beautifulsoup4

EXPOSE 9091

CMD [ "python --version" ]

