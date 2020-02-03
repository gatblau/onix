import React from "react";
import { BrowserRouter as Router, Redirect, Route } from "react-router-dom";
import { useSelector } from "react-redux";
import Login from "./auth/Login";
import Console from "./console";
import "./app.scss";
import axios from "axios";

const App: React.FunctionComponent = (props) => {
  const user = useSelector(store => store.AuthReducer.user);

  if (user.token === undefined || user.token === "" || user.token === null) {
    return (
      <Router>
        <Route path="/login" exact component={Login}/>
        <Redirect to={"/login"}/>
      </Router>);
  } else {

    // set up token injection for all http requests
    axios.interceptors.request.use(function (config) {

      console.log("+++++++++++++++++Configuring");

      config.headers.Authorization = `Basic ${user.token}`;
      return config;
    }, function (err) {
      return Promise.reject(err);
    });

    return (
      <Router>
        <Route path="/console" exact component={Console}/>
        <Redirect to={"/console"}/>
      </Router>);
  }
};
export default App;
