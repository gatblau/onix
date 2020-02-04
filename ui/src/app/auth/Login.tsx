import React, { MouseEvent, useState } from "react";
import { Redirect, withRouter } from "react-router-dom";
import "@patternfly/react-core/dist/styles/base.css";
import "assets/fonts.css";
import brandImg from "assets/logo-text_64.png";
import { BackgroundImageSrc, ListItem, LoginFooterItem, LoginForm, LoginPage } from "@patternfly/react-core";
import bg_image from "assets/images/bg_4k.jpg";
import { useDispatch } from "react-redux";
import UserProfile from "data/userProfile";
import { ACTIONS } from "./authReducer";
import axios from "axios";

const Login: React.ComponentClass<{}> = withRouter((props) => {
  const dispatch = useDispatch();
  const [isLoggedIn, setLoggedIn] = useState(false);
  const [isError, setIsError] = useState(false);
  const [authToken, setAuthToken] = useState("");
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const isValidPassword = () => {
    return password !== "";
  };
  const isValidUsername = () => {
    return username !== "";
  };
  const loginHandler = (event: MouseEvent) => {
    event.preventDefault();
    if (isValidUsername() && isValidPassword()) {

      // constructs a basic authentication token (base64 encoded)
      setAuthToken(btoa(`${username}:${password}`));

      // use the /user api to test if the username and password is correct
      const url = "/api/user";
      axios.get(url, {
          headers: {
            "Authorization": `Basic ${authToken}`
          }
        }
      ).then(response => {
          setLoggedIn(true);
        }
      ).catch(error => console.error(error));
    }
  };
  const listItem = (
    <>
      <ListItem>
        <LoginFooterItem href="#">Terms of Use </LoginFooterItem>
      </ListItem>
      <ListItem>
        <LoginFooterItem href="#">Help</LoginFooterItem>
      </ListItem>
      <ListItem>
        <LoginFooterItem href="#">Privacy Policy</LoginFooterItem>
      </ListItem>
    </>
  );
  const images = {
    [BackgroundImageSrc.lg]: bg_image,
    [BackgroundImageSrc.sm]: bg_image,
    [BackgroundImageSrc.sm2x]: bg_image,
    [BackgroundImageSrc.xs]: bg_image,
    [BackgroundImageSrc.xs2x]: bg_image
  };
  const loginForm = (
    <LoginForm
      // showHelperText={this.state.showHelperText}
      // helperText={helperText}
      usernameLabel="Username"
      usernameValue={username}
      onChangeUsername={setUsername}
      isValidUsername={isValidUsername()}
      passwordLabel="Password"
      passwordValue={password}
      onChangePassword={setPassword}
      isValidPassword={isValidPassword()}
      // isRememberMeChecked={this.state.isRememberMeChecked}
      onLoginButtonClick={loginHandler}
    />
  );

  if (isLoggedIn) {
    let user = new UserProfile();
    user.token = authToken;
    dispatch({type: ACTIONS.SET_USER, user: user});
    return <Redirect to="/"/>;
  }

  return (
    <LoginPage
      brandImgSrc={brandImg}
      brandImgAlt="Onix logo"
      backgroundImgSrc={images}
      backgroundImgAlt="Images"
      footerListItems={listItem}
      textContent="Onix Management Console"
      loginTitle="Log in to your account"
      loginSubtitle="Please use your Local Onix credentials"
    >
      {loginForm}
    </LoginPage>
  );
});

export default Login;
