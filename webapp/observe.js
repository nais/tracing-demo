import {faro} from "@grafana/faro-web-sdk";

// get OTel trace and context APIs
const {trace, context} = faro.api.getOTEL();

const tracer = trace.getTracer("default");

// Generate a Fibonacci number, and then generate a new Fibonacci number from the result.
const traceFibonacci = () => {

    const span = tracer.startSpan("Fibonacci button clicked");
    const number = parseInt(document.getElementById("fibonacciNumber").value)
    console.log(number);

    span.setAttribute("input", number)
    context.with(trace.setSpan(context.active(), span), () => {

        console.info("Generating the first Fibonacci number")
        fetch("/api/", {method: "POST", body: JSON.stringify({Number: number})})
            .then((response) => response.json())
            .then((data) => {
                span.setAttribute("result1", data.Number)
                console.info("Generating the second Fibonacci number")
                fetch("/api/", {
                    method: "POST",
                    body: JSON.stringify({Number: number}),
                })
                    .then((response) => response.json()
                        .then((data) => {
                            span.setAttribute("result2", data.Number)
                            faro.api.pushMeasurement({
                                type: "custom",
                                values: {
                                    fibonacci_result: data.Number
                                },
                            });
                            document.getElementById("fibonacciResult").innerText = "Result: " + data.Number;
                            console.info("The final results are in", data.Number);
                        }))
                    .catch((e) => {
                        faro.api.pushLog([`got an error: ${e}`]);
                    })
            })
            .catch((e) => {
                faro.api.pushLog([`got an error: ${e}`]);
            });

        faro.api.pushLog(["nais tracing says hello"]);

        span.end();
    });

}

const throwButton = document.getElementById("throwButton");
throwButton.addEventListener("click", function () {
    throw new Error("An exception was thrown!");
});

const consoleErrorButton = document.getElementById("consoleErrorButton");
consoleErrorButton.addEventListener("click", function () {
    console.error("An error was logged to the console!");
});

const faroEventButton = document.getElementById("faroEventButton");
faroEventButton.addEventListener("click", function () {
    faro.api.pushEvent("buttonClicked", {buttonId: "faroEventButton"});
});

const fibonacciButton = document.getElementById("fibonacciButton");
fibonacciButton.addEventListener("click", function () {
    traceFibonacci()
});
