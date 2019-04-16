FROM golang as build
WORKDIR /go/src/github.com/vbatts/quay-scheduled-build
COPY . .
ENV GO111MODULE=on
RUN go version && \
    go env && \
    go install -tags netgo github.com/vbatts/quay-scheduled-build

FROM scratch
COPY --from=build /go/bin/quay-scheduled-build /quay-scheduled-build
ENTRYPOINT ["/quay-scheduled-build"]
