import { getWebInstrumentations, initializeFaro } from '@grafana/faro-web-sdk';
import { TracingInstrumentation } from '@grafana/faro-web-tracing';

const faro = initializeFaro({
    url: import.meta.env.VITE_TELEMETRY_ENDPOINT,
    app: {
        name: 'tracing-demo',
        version: '1.0.0',
    },
    instrumentations: [...getWebInstrumentations(), new TracingInstrumentation()],
});
