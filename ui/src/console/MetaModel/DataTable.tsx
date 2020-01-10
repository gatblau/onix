import * as React from "react";
import { useSelector } from "react-redux";

const DataTable: React.FunctionComponent<{ name: string}> = (props) => {
  const node = useSelector(state => state.MetaModelReducer.node);
  return (
    <>
      ==================
      <h1>{props.name}</h1>
      <h1>******{node}</h1>
      ===================
    </>);
};

export default DataTable;
