# Development
FROM golang:1.13.10-alpine AS development
WORKDIR /go/src/github.com/tidepool-org/liftDataImport
RUN adduser -D tidepool && \
    chown -R tidepool /go/src/github.com/tidepool-org/liftDataImport
USER tidepool
COPY --chown=tidepool . .
RUN ./build.sh
CMD ["./dist/liftDataImport"]
