FROM golang:1.19.3-alpine3.15 as builder
RUN mkdir /build
WORKDIR /build
RUN apk add make tree bash
COPY go.* ./
RUN go mod download
COPY . /build/
ENV CGO_ENABLED=0
ENV GOOS=linux
RUN mkdir -p /build/bin && CGO_ENABLED=0 GOOS=linux go build -o /build/bin/ds-node ./cmd/ds-node
RUN mkdir -p /build/bin && CGO_ENABLED=0 GOOS=linux go build -o /build/bin/wait-for ./cmd/wait-for
RUN mkdir -p /build/bin && CGO_ENABLED=0 GOOS=linux go test -c ./test/integration && mv integration.test /build/bin/


FROM alpine:3.16
COPY --from=builder /build/bin/ds-node /ds-node
COPY --from=builder /build/bin/wait-for /wait-for
COPY --from=builder /build/bin/integration.test /integration.test

WORKDIR /
ENTRYPOINT [ "/ds-node" ]
CMD []
