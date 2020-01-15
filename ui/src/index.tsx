import React from "react";
import ReactDOM from "react-dom";
import App from "./app/App";
import { createStore } from "redux";
import { Provider } from "react-redux";
import rootReducer from "./data/rootReducer";

if (process.env.NODE_ENV !== "production") {
  // tslint:disable-next-line
  const axe = require("react-axe");
  axe(React, ReactDOM, 1000);
}
const store = createStore(rootReducer);
const app = (
  <Provider store={store}>
    <App/>
  </Provider>
);

ReactDOM.render(app, document.getElementById("root") as HTMLElement);
