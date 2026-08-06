export const config = {
  host: process.env.YJS_WORKER_HOST ?? "0.0.0.0",
  port: (() => {
    const portString = process.env.YJS_WORKER_PORT;
    if (portString === undefined || portString === "") {
      return 8787;
    }

    const port = Number(portString);
    if (!Number.isInteger(port) || port < 1 || port > 65535) {
      throw new Error("YJS_WORKER_PORT must be an integer between 1 and 65535");
    }

    return port;
  })(),
  telemetry: {
    serviceName: process.env.OTEL_SERVICE_NAME ?? "notezy-yjs-worker",
    serviceVersion: process.env.OTEL_SERVICE_VERSION ?? "0.1.0",
    deploymentEnvironment:
      process.env.OTEL_DEPLOYMENT_ENVIRONMENT ??
      process.env.NODE_ENV ??
      "development",
    serviceInstanceId:
      process.env.OTEL_SERVICE_INSTANCE_ID ?? process.env.HOSTNAME ?? "unknown",
    otlpEndpoint: (
      process.env.OTEL_EXPORTER_OTLP_ENDPOINT ??
      "http://notezy-otel-collector:4318"
    ).replace(/\/$/, ""),
  },
};
