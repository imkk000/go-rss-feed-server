FROM docker.io/golang:1.26-alpine3.23 as builder

WORKDIR /builder
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build --ldflags="-w -s -X main.addr=:80" -o main .

FROM docker.io/alpine:3.23

WORKDIR /opt/rss
COPY --from=builder /builder/main server
RUN chown nobody:nobody server

USER nobody
CMD ["/opt/rss/server"]
