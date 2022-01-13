const express = require("express");
const path = require("path");

const app = express();

app.use("/static/js", express.static(path.resolve(__dirname, "js")));
app.use("/static/css", express.static(path.resolve(__dirname, "css")));

app.get("/*", (req, res) => {
  res.sendFile(path.resolve(__dirname, "index.html"));
});

app.listen(process.env.PORT || 9101, () => console.log("Server running..."));
