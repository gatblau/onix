import { combineReducers } from "redux";
import { REDUCER as MetaModelReducer } from "console/MetaModel/store/metamodelReducer";
import { REDUCER as AuthReducer } from "../auth/authReducer";

export default combineReducers({
  AuthReducer,
  MetaModelReducer
});
