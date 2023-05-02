const express = require("express")


const app = express();
const port = 8080;

app.get("/api/", (req, res) => {
    console.log("wee /")
    res.send("Hello World!")

});

app.listen(port, () => {
    console.log(`Example app listening on port ${port}`);
});
