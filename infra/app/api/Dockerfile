FROM golang:1.20.12-alpine3.18 as builder
RUN apk --no-cache add build-base
WORKDIR /go/src/github.com/redreamerlab/api
COPY . /go/src/github.com/redreamerlab/api
RUN CGO_ENABLED=1 go build -o ./bin/app ./cmd/api/main.go

FROM alpine:3.15.7
RUN apk add --no-cache ca-certificates libc6-compat libgcc libstdc++
WORKDIR /
COPY --from=builder /go/src/github.com/redreamerlab/api/bin/app /app/app
ENTRYPOINT [ "/app/app" ]
