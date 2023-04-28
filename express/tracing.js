const OpenTelemetry = require("@opentelemetry/sdk-node");
const Resources = require("@opentelemetry/resources");
const SemanticConventions = require("@opentelemetry/semantic-conventions");
const InstrumentationHttp = require("@opentelemetry/instrumentation-http");
const InstrumentationExpress = require("@opentelemetry/instrumentation-express");
const ExporterTraceOtlpHttp = require("@opentelemetry/exporter-trace-otlp-http");
const TraceNode = require("@opentelemetry/sdk-trace-node")
const { getNodeAutoInstrumentations } = require('@opentelemetry/auto-instrumentations-node');

const sdk = new OpenTelemetry.NodeSDK({
  resource: new Resources.Resource({
    [SemanticConventions.SemanticResourceAttributes.SERVICE_NAME]: "my-service",
  }),
    traceExporter: new ExporterTraceOtlpHttp.OTLPTraceExporter({
        url: "http://127.0.0.1:4317"
   }),
  instrumentations: [
    new InstrumentationHttp.HttpInstrumentation(),
      new InstrumentationExpress.ExpressInstrumentation(),
      getNodeAutoInstrumentations()

  ],
});

console.log("sdk")
sdk.start()
