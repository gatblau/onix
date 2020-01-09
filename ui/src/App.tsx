import React from "react";
import { BrowserRouter as Router, Redirect, Route } from "react-router-dom";
import Login from "./Login/Login";
import Console from "./console";
import "./app.scss"

const App: React.FunctionComponent = (props) => {
  return (
    <Router>
      <Route path="/login" exact component={Login}/>
      <Redirect to={"/login"}/>
      <Route path="/" exact component={Console}/>
    </Router>
  );
};

export default App;
