import React from "react";
import { BrowserRouter as Router, Redirect, Route } from "react-router-dom";
import Login from "./auth/Login";
import Console from "./console";
import "./app.scss"
import { useSelector } from "react-redux";

const App: React.FunctionComponent = (props) => {
  const user = useSelector(state => state.AuthReducer.user);

  console.log("$$$$$$$$$$$$$$$$$$$$$$$$$$$$", user);

  return (
    <Router>
      <Route path="/login" exact component={Login}/>
      {user.token === undefined && <Redirect to={"/login"}/>}
      <Route path="/" exact component={Console}/>
    </Router>
  );
};

export default App;
