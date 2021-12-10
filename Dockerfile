FROM golang:1.17 AS builder
ENV LANG=C.UTF-8

WORKDIR /restic
RUN git clone https://github.com/restic/restic /restic
RUN go run build.go -o /out/restic

WORKDIR /build
COPY . /build/
RUN go build -ldflags="-w -s" -o /out/EzBackup

FROM gcr.io/distroless/base-debian11
COPY --from=builder /out/ /app/
ENTRYPOINT ["/app/EzBackup"]
CMD ["help"]