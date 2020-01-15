import * as React from "react";
import { Graph } from "react-d3-graph";
import { useDispatch } from "react-redux";
import { ACTIONS } from "./data/metamodelReducer";

const DataGraph: React.FunctionComponent<{}> = () => {
  const dispatch = useDispatch();

// graph payload (with minimalist structure)
  const data = {
    nodes: [{id: "Bob", symbolType: "diamond"}, {id: "Carol"}, {id: "Alice"}],
    links: [{source: "Bob", target: "Carol"}, {source: "Bob", target: "Alice"}]
  };

// the graph configuration, you only need to pass down properties
// that you want to override, otherwise default ones will be used
  const myConfig = {

    nodeHighlightBehavior: true,
    node: {
      color: "red",
      size: 500,
      highlightStrokeColor: "blue"
    },
    link: {
      highlightColor: "lightblue"
    }
  };

// graph event callbacks
  const onClickGraph = function () {
    window.alert(`Clicked the graph background`);
  };

  const onClickNode = (nodeId: string) => {
    dispatch({type: ACTIONS.SET_NODE, node: nodeId});
  };

  const onDoubleClickNode = (nodeId: string) => {
  };

  const onRightClickNode = (event: MouseEvent, nodeId: string) => {
  };

  const onMouseOverNode = (nodeId: string) => {
  };

  const onMouseOutNode = (nodeId: string) => {
  };

  const onClickLink = (source: string, target: string) => {
  };

  const onRightClickLink = (event: MouseEvent, source: string, target: string) => {
  };

  const onMouseOverLink = (source: string, target: string) => {
  };

  const onMouseOutLink = (source: string, target: string) => {
  };

  const onNodePositionChange = (nodeId: string, x: number, y: number) => {
  };

  return (
    <Graph
      id="graph-id" // id is mandatory, if no id is defined rd3g will throw an error
      data={data}
      config={myConfig}
      onClickNode={onClickNode}
      onDoubleClickNode={onDoubleClickNode}
      onRightClickNode={onRightClickNode}
      onClickGraph={onClickGraph}
      onClickLink={onClickLink}
      onRightClickLink={onRightClickLink}
      onMouseOverNode={onMouseOverNode}
      onMouseOutNode={onMouseOutNode}
      onMouseOverLink={onMouseOverLink}
      onMouseOutLink={onMouseOutLink}
      onNodePositionChange={onNodePositionChange}
    />
  );
};

export default DataGraph ;
