import React, { MouseEvent, useState } from "react";
import "@patternfly/react-core/dist/styles/base.css";
import "assets/fonts.css";
import brandImg from "assets/logo-text_64.png";
import { BackgroundImageSrc, ListItem, LoginFooterItem, LoginForm, LoginPage } from "@patternfly/react-core";
import bg_image from "assets/images/bg_4k.jpg";
import { withRouter } from "react-router-dom";
import { useDispatch } from "react-redux";
import UserProfile from "../store/userProfile";
import { ACTIONS } from "./authReducer";

const Login: React.ComponentClass<{}> = withRouter((props) => {
  const dispatch = useDispatch();
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
      let user = new UserProfile();
      user.token = "TODO - read from back end";
      dispatch({type: ACTIONS.SET_USER, user: user});
      props.history.push("/");
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
