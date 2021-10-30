FROM gcr.io/distroless/static
COPY . /
ENTRYPOINT ["/firewall_sync"]
CMD ["-c", "config.yml"]