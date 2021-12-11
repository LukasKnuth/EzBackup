# Run the builder on the native architecture of the building computer
FROM --platform=$BUILDPLATFORM golang:1.17 AS builder
ENV LANG=C.UTF-8
# Build truely static executables, no C dependencies.
ENV CGO_ENABLED=0

# Download/Copy sources
RUN git clone https://github.com/restic/restic /restic
COPY . /build/

# Download dependencies
RUN cd /restic && go mod download
RUN cd /build && go mod download

# Architecture to build the executable for (set by "buildx")
# This is where the build process splits, everything before this is cached/executed once!
ARG TARGETOS TARGETARCH

# Build EzBackup
WORKDIR /build
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -ldflags="-w -s" -o /out/EzBackup

# Build restic
WORKDIR /restic
RUN go run build.go --goos $TARGETOS --goarch $TARGETARCH -o /out/restic

FROM gcr.io/distroless/base-debian11
COPY --from=builder /out/ /app/
ENTRYPOINT ["/app/EzBackup"]
CMD ["help"]