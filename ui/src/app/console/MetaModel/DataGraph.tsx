import * as React from "react";
import { useEffect, useState } from "react";
import { Graph } from "react-d3-graph";
import { useDispatch } from "react-redux";
import { ACTIONS } from "./data/metamodelDatalReducer";
import { useParams } from "react-router";
import axios from "axios";

const DataGraph: React.FunctionComponent<{}> = () => {
  const dispatch = useDispatch();

  // graph payload
  const [data, setData] = useState({nodes: [], links: []});

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

  const {id} = useParams();

  useEffect(() => {
    axios.get(`/api/model/${id}/data`
    ).then(response => {
        const itemTypes: any[] = [];
        response.data.itemTypes.forEach((item, idx) => {
          itemTypes.push({id: item.key});
        });

        const linkRules: any[] = [];
        response.data.linkRules.forEach((item, idx) => {
          const linkRule:string[] = item.key.split("->");
          linkRules.push({source: linkRule[0], target: linkRule[1]});
        });

        // @ts-ignore
        setData({nodes: itemTypes, links: linkRules});
      }
    ).catch(error => console.error(error));
  }, []);

  if (data.nodes.length > 0) {
    return (
      <Graph
        id="metaModel" // id is mandatory, if no id is defined rd3g will throw an error
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
  }

  return (
    <h1>Loading</h1>
  );
};

export default DataGraph;
