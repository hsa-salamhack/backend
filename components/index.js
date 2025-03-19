import express from "express";
import { wiki } from "./wikiparse.js";
const app = express();
const port = 5050;

app.get("/wiki/:lang/:term", async (req, res) => {
  let { lang, term } = req.params;
  res.json(await wiki(term, lang));
});

app.listen(port);
