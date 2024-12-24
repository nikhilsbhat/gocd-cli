### Description: Dockerfile for gocd-prometheus-exporter
FROM alpine:3.21

COPY gocd-cli /

# Starting
ENTRYPOINT [ "/gocd-cli" ]