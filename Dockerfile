FROM --platform=$BUILDPLATFORM golang:1.26 as builder
ARG TARGETOS
ARG TARGETARCH
COPY . /src

WORKDIR /src

# Build statically linked file and strip debug information
RUN go get
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="-extldflags=-static -s -w" -o horaco_exporter

FROM scratch
COPY --from=builder /src/horaco_exporter /horaco_exporter

ENTRYPOINT ["/horaco_exporter"]
