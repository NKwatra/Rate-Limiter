const express = require("express");
const cors = require("cors");

const app = express();

app.use(cors());
app.get("/limited", (_, res) =>
  res.json({ message: "Rate limited end point" })
);
app.get("/unlimited", (_, res) =>
  res.json({ message: "Use as much as required. Enjoy!!" })
);
app.listen(8080, () => console.log("Server started successfully"));
