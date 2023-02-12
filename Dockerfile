### Description: Dockerfile for gocd-prometheus-exporter
FROM alpine:3.16

COPY gocd-cli /

# Starting
ENTRYPOINT [ "/gocd-cli" ]