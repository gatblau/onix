import React from "react";
import { BrowserRouter as Router, Redirect, Route } from "react-router-dom";
import { useSelector } from "react-redux";
import Login from "./auth/Login";
import Console from "./console";
import "./app.scss";

const App: React.FunctionComponent = (props) => {
  const user = useSelector(state => state.AuthReducer.user);

  if (user.token === undefined) {
    return (
      <Router>
        <Route path="/login" exact component={Login}/>
        <Redirect to={"/login"}/>
      </Router>);
  } else {
    return (
      <Router>
        <Route path="/console" exact component={Console}/>
        <Redirect to={"/console"}/>
      </Router>);
  }
};
export default App;
