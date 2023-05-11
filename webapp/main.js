import { getWebInstrumentations, initializeFaro } from '@grafana/faro-web-sdk';
import { TracingInstrumentation } from '@grafana/faro-web-tracing';
import nais from './nais.js';

const faro = initializeFaro({
    url: nais.telemetryCollectorURL,
    app: {
        name: 'tracing-demo',
        version: '1.0.0',
    },
    instrumentations: [...getWebInstrumentations(), new TracingInstrumentation()],
});
