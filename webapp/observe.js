import {faro} from '@grafana/faro-web-sdk';

// get OTel trace and context APIs
const {trace, context} = faro.api.getOTEL();

// const tracer = trace.getTracer('default');
// const span = tracer.startSpan('hello');

// context.with(trace.setSpan(context.active(), span), () => {
    // faro.api.pushMeasurement({
    //     type: 'custom',
    //     values: {
    //         nais_tracing_answer: 42,
    //     },
    // });
    //
    fetch("/api/", {method: "POST", body: JSON.stringify({"Number": 42})}).then((response) => {
        console.log(response)
    });

    // faro.api.pushLog(['nais tracing says hello']);

    // span.end();
// });

