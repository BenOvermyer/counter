# build stage
FROM golang:1.13 AS build-env
WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/counter

# final stage
FROM scratch
COPY --from=build-env /go/bin/counter /go/bin/counter
COPY --from=build-env /go/src/app/Berylium.ttf /go/bin/Berylium.ttf
EXPOSE 9776
CMD ["/go/bin/counter"]
