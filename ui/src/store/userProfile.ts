export default class UserProfile {
  constructor(private _username?,
              private _firstName?,
              private _lastName?,
              private _token?: string) {
  }

  get token(): string {
    return <string>this._token;
  }

  set token(value: string) {
    this._token = value;
  }

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

  clone = () => {
    let newInstance = new UserProfile();

    newInstance.username = this.username;
    newInstance.firstName = this.firstName;
    newInstance.lastName = this.lastName;
    newInstance.token = this.token;

    return newInstance;
  }
}
