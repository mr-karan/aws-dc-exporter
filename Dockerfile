FROM alpine:latest AS deploy
RUN apk --no-cache add ca-certificates
COPY aws-dc-exporter /
COPY config.sample.toml  /etc/aws-dc-exporter/config.toml
VOLUME ["/etc/aws-dc-exporter"]
CMD ["./aws-dc-exporter", "--config", "/etc/aws-dc-exporter/config.toml"]  
