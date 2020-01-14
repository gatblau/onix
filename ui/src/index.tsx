import React from "react";
import ReactDOM from "react-dom";
import App from "./App";
import { createStore } from "redux";
import { Provider } from "react-redux";
import rootReducer from "./store/rootReducer";
import { REDUCER as MetaModelReducer } from "console/MetaModel/store/metamodelReducer";

if (process.env.NODE_ENV !== "production") {
  // tslint:disable-next-line
  const axe = require("react-axe"); // eslint-disable-line
  axe(React, ReactDOM, 1000);
}

const store = createStore(rootReducer);

console.log("!!!!!!!!!!!!", store.getState());

const app = (
  <Provider store={store}>
    <App/>
  </Provider>
);

ReactDOM.render(app, document.getElementById("root") as HTMLElement);
