FROM golang:latest

WORKDIR /go/src/tg-notifier
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

VOLUME "./data"

ENTRYPOINT ["tg-notifier"]
EXPOSE 8080