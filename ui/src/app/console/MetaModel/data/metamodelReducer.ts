const ACTIONS = {
  SET_MODELS: "SET_MODELS"
};

const initialState = {
  models: []
};

const REDUCER = (state = initialState, action) => {
  let result = state;

  if (action.type === ACTIONS.SET_MODELS) {
    result = {
      ...state,
      models: action.models
    };
  }

  return result;
};

export {ACTIONS, REDUCER};
