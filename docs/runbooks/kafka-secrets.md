# Kafka TLS secrets

When Kafka TLS is enabled, provision the certificate material outside the
repository and mount it into Compose at `/run/secrets/kafka/`.

Expected files are:

```text
ca.pem       # Kafka CA certificate
client.crt   # optional mutual-TLS client certificate
client.key   # optional mutual-TLS client key
```

The production Compose file mounts the local `secrets/kafka/` directory
read-only. Keep that directory, including documentation and certificates,
outside Git. On deployment hosts, provision it through the host secret manager
or the deployment platform before starting Compose.
