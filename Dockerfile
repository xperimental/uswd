FROM golang:1 AS builder

ENV CGO_ENABLED=0

WORKDIR /go/src/github.com/xperimental/uswd
COPY . .

RUN go test ./...
RUN go install -v -tags netgo -ldflags "-w" ./cmd/uswd-server

FROM busybox
COPY --from=builder /go/bin/uswd-server /bin/uswd-server
RUN mkdir /data

EXPOSE 8080
ENTRYPOINT [ "/bin/uswd-server" ]
CMD [ "--base", "/data" ]
