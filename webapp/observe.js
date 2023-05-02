import {faro} from '@grafana/faro-web-sdk';

// get OTel trace and context APIs
const {trace, context} = faro.api.getOTEL();

const tracer = trace.getTracer('default');
const span = tracer.startSpan("demo frontend span");

context.with(trace.setSpan(context.active(), span), () => {
    faro.api.pushMeasurement({
        type: 'custom',
        values: {
            nais_tracing_answer: 42,
        },
    });

    fetch("/api/", {method: "POST", body: JSON.stringify({"Number": 5})})
        .then((response) => response.json())
        .then((data) => {
            fetch("/api/", {method: "POST", body: JSON.stringify({"Number": data.Number})}).then((response) => {
                console.log(response)
            })
        })
        .catch((e) => {
            faro.api.pushLog([`got an error: ${e}`]);
        });

    faro.api.pushLog(['nais tracing says hello']);

    span.end();
});
