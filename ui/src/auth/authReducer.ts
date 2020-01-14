import UserProfile from "../store/userProfile";

export const ACTIONS = {
  SET_USER: "SET_USER"
};

const initialState = {
  user: new UserProfile()
};

export const REDUCER = (state = initialState, action) => {
  switch (action) {
    case ACTIONS.SET_USER :
      return {user: <UserProfile>state.user.clone()};
    default:
      return state;
  }
};

