import { combineReducers } from "redux";
import { REDUCER as MetaModelReducer } from "app/console/MetaModel/data/metamodelReducer";
import { REDUCER as AuthReducer } from "app/auth/authReducer";

export default combineReducers({
  AuthReducer,
  MetaModelReducer
});
