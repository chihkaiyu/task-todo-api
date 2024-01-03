
FROM golang:1.20.12-alpine AS swagger-builder

RUN apk add --no-cache g++ musl-dev libstdc++ git

# To build the swagger by the same version as the compiler to build app binary
# build swagger command and using ldflags to replace the swagger version
# ref: https://github.com/go-swagger/go-swagger/blob/master/.circleci/config.yml#L51
ENV SWAGGO_VERSION v1.8.10
RUN go install github.com/swaggo/swag/cmd/swag@${SWAGGO_VERSION}

FROM golang:1.20.12-alpine

RUN apk add --no-cache g++ musl-dev libstdc++

COPY --from=swagger-builder /go/bin/swag /go/bin/swag
