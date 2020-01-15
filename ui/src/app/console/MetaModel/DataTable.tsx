import * as React from "react";
import { useState } from "react";
import { dispatch, useSelector } from "react-redux";
import { sortable, SortByDirection, Table, TableBody, TableHeader } from "@patternfly/react-table";
import { useEffect } from "react";
import axios from "axios";

const DataTable: React.FunctionComponent<{}> = (props) => {
  const node = useSelector(store => store.MetaModelReducer.node);
  const [state, setState] = useState({
    columns: [
      {title: "Repositories", transforms: [sortable]},
      "Branches",

      {title: "Pull requests", transforms: [sortable]},
      "Workspaces",
      "Last Commit"
    ],
    rows: [],
    sortBy: {}
  });

  const onSort = (_event, index, direction) => {
    const sortedRows = state.rows.sort((a, b) => (a[index] < b[index] ? -1 : a[index] > b[index] ? 1 : 0));
    setState({
      columns: state.columns,
      sortBy: {
        index,
        direction
      },
      rows: direction === SortByDirection.asc ? sortedRows : sortedRows.reverse()
    });
  };

  useEffect(() => {
    axios.get("/" + node + ".json")
      .then(response => {
        console.log("====> Got ", response.data);
        setState({
          columns: state.columns,
          sortBy: {},
          rows: response.data.rows
        });
      })
      .catch(error => console.error(error));

  }, [node]);

  return (
    <>
      <Table aria-label="Sortable Table" sortBy={state.sortBy} onSort={onSort} cells={state.columns} rows={state.rows}>
        <TableHeader/>
        <TableBody/>
      </Table>
    </>
  );
};

export default DataTable;
