import * as React from "react";
import { Graph } from "react-d3-graph";
import { useDispatch } from "react-redux";
import { ACTIONS } from "./store/metamodelReducer";

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
    window.alert(`Clicked node ${nodeId}`);
  };

  const onDoubleClickNode = (nodeId: string) => {
    window.alert(`Double clicked node ${nodeId}`);
  };

  const onRightClickNode = (event: MouseEvent, nodeId: string) => {
    window.alert(`Right clicked node ${nodeId}`);
  };

  const onMouseOverNode = (nodeId: string) => {
    window.alert(`Mouse over node ${nodeId}`);
  };

  const onMouseOutNode = (nodeId: string) => {
    window.alert(`Mouse out node ${nodeId}`);
  };

  const onClickLink = (source: string, target: string) => {
    window.alert(`Clicked link between ${source} and ${target}`);
  };

  const onRightClickLink = (event: MouseEvent, source: string, target: string) => {
    window.alert(`Right clicked link between ${source} and ${target}`);
  };

  const onMouseOverLink = (source: string, target: string) => {
    window.alert(`Mouse over in link between ${source} and ${target}`);
  };

  const onMouseOutLink = (source: string, target: string) => {
    window.alert(`Mouse out link between ${source} and ${target}`);
  };

  const onNodePositionChange = (nodeId: string, x: number, y: number) => {
    window.alert(`Node ${nodeId} is moved to new position. New position is x= ${x} y= ${y}`);
  };

  return (
    <Graph
      id="graph-id" // id is mandatory, if no id is defined rd3g will throw an error
      data={data}
      config={myConfig}
      onClickNode={onClickNode}
      onRightClickNode={onRightClickNode}
      onClickGraph={onClickGraph}
      onClickLink={onClickLink}
      onRightClickLink={onRightClickLink}
      // onMouseOverNode={onMouseOverNode}
      // onMouseOutNode={onMouseOutNode}
      // onMouseOverLink={onMouseOverLink}
      // onMouseOutLink={onMouseOutLink}
      onNodePositionChange={onNodePositionChange}
    />
  );
};

export default DataGraph ;
