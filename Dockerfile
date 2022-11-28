FROM golang:1.18-alpine as buildbase

RUN apk add git build-base

WORKDIR /go/src/github.com/Swapica/swapica-svc
COPY vendor .
COPY . .

RUN GOOS=linux go build  -o /usr/local/bin/swapica-svc /go/src/github.com/Swapica/swapica-svc


FROM alpine:3.9

COPY --from=buildbase /usr/local/bin/swapica-svc /usr/local/bin/swapica-svc
RUN apk add --no-cache ca-certificates

ENTRYPOINT ["swapica-svc"]
