# build stage
FROM golang:1.17-alpine
WORKDIR /counter/
COPY . /counter/
RUN go get -d -v ./...
RUN apk add gcc musl-dev sqlite
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /counter/counter

CMD ["/counter/counter"]
