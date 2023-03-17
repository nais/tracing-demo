import { getWebInstrumentations, initializeFaro } from '@grafana/faro-web-sdk';
import { TracingInstrumentation } from '@grafana/faro-web-tracing';

const faro = initializeFaro({
    url: 'https://frontend-metric-collector.dev.dev-nais.cloud.nais.io/collect',
    app: {
        name: 'tracing-demo',
        version: '1.0.0',
    },
    instrumentations: [...getWebInstrumentations(), new TracingInstrumentation()],
});
