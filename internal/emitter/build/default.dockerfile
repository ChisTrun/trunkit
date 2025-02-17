# Use the official Alpine Linux image as the base image
FROM registry.ugaming.io/marketplace/cicd/golang:1.18.2-builder-1.0 AS builder

#final stage
FROM alpine:3.15
RUN apk --no-cache add ca-certificates bash

COPY --from=builder /tools/grpc_health_probe /grpc_health_probe
COPY --from=builder /tools/wait-for-it.sh ./wait-for-it.sh
COPY --from=builder /tools/envoy-preflight  /envoy-preflight
COPY ./cmd/main /app/server
COPY configs /app/
EXPOSE 8080

CMD ["/app/server", "-c", "/app/config.yaml"]
