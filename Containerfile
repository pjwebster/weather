FROM registry.access.redhat.com/ubi8/go-toolset:1.16.12-2 AS builder
ADD . /builder
USER 0
RUN chown -R 1001:0 /builder
USER 1001
WORKDIR /builder
RUN go build .

FROM registry.access.redhat.com/ubi8/ubi-minimal AS runner
RUN mkdir /app
WORKDIR /app

COPY --from=builder /builder/weather .
COPY ./dist ./dist

EXPOSE 8090
ENTRYPOINT ["/app/weather"]
