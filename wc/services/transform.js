/*
 *   Onix Web Console - Copyright (c) 2019 by www.gatblau.org
 *
 *   Licensed under the Apache License, Version 2.0 (the "License");
 *   you may not use this file except in compliance with the License.
 *   You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
 *   Unless required by applicable law or agreed to in writing, software distributed under
 *   the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 *   either express or implied.
 *   See the License for the specific language governing permissions and limitations under the License.
 *
 *   Contributors to this project, hereby assign copyright in this code to the project,
 *   to be licensed under the same terms as the rest of the code.
*/

// transforms model data to the format required by d3-network
export function transformGraph(data) {
    const output = {
        nodes: [
            {
                "id": "MODELS",
                "name": "Models",
                "_color": "orange"
            },
        ],
        links: [],
    };
    output.nodes = output.nodes.concat(transformModels(data.models))
    output.nodes = output.nodes.concat(transformItemTypes(data.itemTypes));
    output.links = output.links.concat(transformLinkRules(data.linkRules));
    output.links = output.links.concat(getModelRootLinks(data));
    return output;
}

function getModelRootLinks(data) {
    const list = [];
    for (let model of data.models) {
        for (let itemType of data.itemTypes) {
            // connects models root with model\
            list.push({
                "sid": "MODELS",
                "tid": model.key
            });
            if (itemType.root && itemType.modelKey == model.key) {
                // connects model with item type root
                list.push({
                    "sid": model.key,
                    "tid": itemType.key
                });
            }
        }
    }
    return list;
}

// transforms models
function transformModels(models) {
    const list = [];
    for (let model of models) {
        list.push({
            "id": model.key,
            "name": model.name,
            "_color": "red"
        });
    }
    return list;
}

// transforms item type data to the format required by d3-network
function transformItemTypes(itemTypes) {
    const list = [];
    for (let item of itemTypes) {
        list.push({
            "id": item.key,
            "name": item.name,
            "_color": "green"
        });
    }
    return list;
}

// transforms link rules data to the format required by d3-network
function transformLinkRules(linkRules) {
    const list = [];
    for (let rule of linkRules) {
        list.push({
            "sid": rule.startItemTypeKey,
            "tid": rule.endItemTypeKey
        });
    }
    return list;
}