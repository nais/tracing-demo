'use strict';
const { ExpressInstrumentation } = require('@opentelemetry/instrumentation-express');
const { HttpInstrumentation } = require('@opentelemetry/instrumentation-http');
const opentelemetry = require("@opentelemetry/sdk-node");
const api = require('@opentelemetry/api');
const { getNodeAutoInstrumentations } = require('@opentelemetry/auto-instrumentations-node');
const { OTLPTraceExporter } = require('@opentelemetry/exporter-trace-otlp-grpc');
const { OTLPMetricExporter } = require("@opentelemetry/exporter-metrics-otlp-grpc");
const { Resource } = require('@opentelemetry/resources');
const { SemanticResourceAttributes } = require('@opentelemetry/semantic-conventions');


// The defaults rul works fine for this demo
const traceOtlpExporter = new OTLPTraceExporter();

// The default url works fine for this demo
const metricOtlpExporter = new OTLPMetricExporter();

const sdk = new opentelemetry.NodeSDK({
 resource: new Resource({
    [ SemanticResourceAttributes.SERVICE_NAME ]: "my-cool-service-name",
    [ SemanticResourceAttributes.SERVICE_NAMESPACE ]: "Nais-namespace",
    [ SemanticResourceAttributes.SERVICE_VERSION ]: "1.0",
 }),
  traceExporter: traceOtlpExporter,
  metricExporter: metricOtlpExporter,
  autoDetectResources: true,
    instrumentations: [new HttpInstrumentation(), new ExpressInstrumentation()],
});

sdk.start()
