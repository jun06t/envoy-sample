FROM envoyproxy/envoy:v1.24.1

ADD envoy.yaml /etc/envoy/envoy.yaml

CMD ["-c", "/etc/envoy/envoy.yaml", "--drain-time-s", "30", "--parent-shutdown-time-s", "60"]

