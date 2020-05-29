# Development
FROM golang:1.13.10-alpine AS development
WORKDIR /go/src/github.com/tidepool-org/lift-data-import
RUN adduser -D tidepool && \
    chown -R tidepool /go/src/github.com/tidepool-org/lift-data-import
USER tidepool
COPY --chown=tidepool . .
RUN ./build.sh
CMD ["./dist/lift-data-import"]
