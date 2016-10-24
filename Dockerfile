FROM golang:1.7

MAINTAINER Denis Isaev <d.isaev@corp.mail.ru>

RUN apt-get update && apt-get install -y --force-yes git

WORKDIR /usr/src/app

ARG OSM_FILE=planet.osm

COPY $OSM_FILE data.asm

WORKDIR gohome/src
RUN git clone https://github.com/jirfag/osm2elastic.git
ENV GOPATH=/usr/src/app/gohome
WORKDIR osm2elastic
RUN go get gopkg.in/olivere/elastic.v2
RUN go get github.com/jirfag/osm2elastic

CMD go run osm2elastic.go --osm-file /usr/src/app/data.asm --elastic-addr $ELASTIC_PORT_9200_TCP_ADDR:$ELASTIC_PORT_9200_TCP_PORT
# sudo docker build -t osm2elastic .
# sudo docker run --name osm2elastic --link elasticsearch:elastic osm2elastic
