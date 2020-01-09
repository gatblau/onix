import React from "react";

export class UserProfile {
  private _username = "";
  private _firstName = "";
  private _lastName = "";

  get username(): string {
    return this._username;
  }

  set username(value: string) {
    this._username = value;
  }

  get firstName(): string {
    return this._firstName;
  }

  set firstName(value: string) {
    this._firstName = value;
  }

  get lastName(): string {
    return this._lastName;
  }

  set lastName(value: string) {
    this._lastName = value;
  }
}

export default React.createContext(
  {
    authenticated: false,
    userProfile: new UserProfile(),
    token: null
  });
