# Development
FROM golang:1.14-alpine AS builder
WORKDIR /go/src/github.com/tidepool-org/kafka-data-import
RUN apk add gcc musl-dev
RUN adduser -D tidepool && \
    chown -R tidepool /go/src/github.com/tidepool-org/kafka-data-import
USER tidepool
COPY --chown=tidepool . .
RUN ./build.sh
CMD ["./dist/kafka-data-import"]
