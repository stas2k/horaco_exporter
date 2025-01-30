FROM golang:1.23 as builder
COPY . /src

WORKDIR /src

# Build statically linked file and strip debug information
RUN go get
RUN CGO_ENABLED=0 go build -ldflags="-extldflags=-static -s -w" -o horaco_exporter

FROM scratch
COPY --from=builder /src/horaco_exporter /horaco_exporter

ENTRYPOINT ["/horaco_exporter"]
