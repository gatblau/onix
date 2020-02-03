const ACTIONS = {
  SET_NODE: 'SET_NODE'
};

const initialState = {
  node: 'None Selected'
};

const REDUCER = (state = initialState, action) => {
  let result = state;

  if (action.type === ACTIONS.SET_NODE) {
    result = {
      ...state,
      node: action.node
    };
  }

  return result;
};

export {ACTIONS, REDUCER};
