FROM registry.access.redhat.com/ubi8/go-toolset:1.16.12-2 AS builder
WORKDIR /builder
COPY . .
RUN go build .

FROM registry.access.redhat.com/ubi8/ubi-minimal AS runner
WORKDIR /app
COPY --from=builder /builder/weather .
COPY ./dist ./dist

EXPOSE 8090
ENTRYPOINT ["/app/weather"]
