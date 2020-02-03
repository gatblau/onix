import UserProfile from "data/userProfile";
import axios from "axios";

export const ACTIONS = {
  SET_USER: "SET_USER"
};

const storedUser = sessionStorage.getItem("USER");
const initialState:{user: UserProfile} = { user : storedUser ? JSON.parse(storedUser) as UserProfile : new UserProfile()};

export const REDUCER = (state = initialState, action) => {
  if (action.type === ACTIONS.SET_USER) {
    let user = action.user.clone();
    sessionStorage.setItem("USER", JSON.stringify(user.toJson()));
    return {user: user};
  }
  return state;
};

