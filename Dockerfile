FROM golang:1.18-alpine

WORKDIR /pblb
COPY . .

RUN apk add git
RUN go build -v ./...
RUN go install -v ./...
RUN mkdir -p /etc/pblb/
COPY config.yaml /etc/pblb/

EXPOSE 2839:2839

CMD [ "pblb", "run"]