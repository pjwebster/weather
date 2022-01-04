FROM registry.access.redhat.com/ubi8/go-toolset:1.16.12-2 AS builder
RUN go build .

FROM registry.access.redhat.com/ubi8/ubi-minimal AS runner
COPY --from=builder weather .
ENTRYPOINT ["./weather"]
