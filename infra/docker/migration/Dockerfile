FROM golang:1.19.13 as builder
WORKDIR /
ENV SQL_MIGRATE v1.0.0
RUN CGO_ENABLED=0 go install github.com/rubenv/sql-migrate/...@${SQL_MIGRATE}

FROM alpine:3.16.2
RUN apk --no-cache add ca-certificates && mkdir -p /app/databases

COPY --from=builder /go/bin/sql-migrate /bin/sql-migrate
COPY entrypoint.sh /app/entrypoint.sh

WORKDIR /app

ENTRYPOINT [ "/bin/sh", "/app/entrypoint.sh" ]
