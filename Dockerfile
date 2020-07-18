FROM golang:1.14-buster AS builder

ARG TAGS

ENV CGO_ENABLED 0
ENV GOOS linux
ENV GOARCH amd64

WORKDIR /src
COPY src /src

RUN go get -d ./... && \
    go build -a -ldflags "-s -w" \
    ${TAGS:+-tags=$TAGS} \
    -o /default-backend ./...

FROM scratch

ARG UID=101

USER $UID

COPY rootfs /

COPY --from=builder /default-backend /

ENTRYPOINT ["/default-backend"]
