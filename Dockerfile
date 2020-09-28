FROM scratch

COPY config.toml config.toml
COPY balanceapp balanceapp
COPY VERSION VERSION
EXPOSE 8080
ENTRYPOINT ["./balanceapp"]
