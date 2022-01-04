FROM golang:1.16 AS builder
RUN go build .

FROM registry.access.redhat.com/ubi8/ubi-minimal AS runner
COPY --from=builder weather .
ENTRYPOINT ["./weather"]
