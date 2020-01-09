import React from "react";
import ReactDOM from "react-dom";
import App from "./App";

if (process.env.NODE_ENV !== "production") {
  // tslint:disable-next-line
  const axe = require("react-axe"); // eslint-disable-line
  axe(React, ReactDOM, 1000);
}

const app = (
  <App/>
);

ReactDOM.render(app, document.getElementById("root") as HTMLElement);
