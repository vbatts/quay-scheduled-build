FROM golang as build
WORKDIR /go/src/github.com/vbatts/quay-scheduled-build
COPY . .
ENV GO111MODULE=on
RUN go version && \
    go env && \
    go install -tags netgo github.com/vbatts/quay-scheduled-build

FROM alpine
COPY --from=build /go/bin/quay-scheduled-build /quay-scheduled-build
RUN apk add ca-certificates
USER 1000
ENTRYPOINT ["/quay-scheduled-build"]
