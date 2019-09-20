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
export const state = () => ({
    options: {
        canvas: false,
        force: 3000,
        nodeSize: 30,
        size: { w:1400, h:700 },
        offset: { x:-250, y:0 },
        nodeLabels: true,
        linkWidth: 2,
        linkLabels: true,
        fontSize: 16,
    },
    items : [],
    itemTypes: [],
    meta : {},
    title: "",
})

export const mutations = {
    setItems(state, data) {
        state.items = [];
        state.items = data.items;
    },
    setItemTypes(state, data) {
        state.itemTypes = [];
        state.itemTypes = data.itemTypes;
    },
    setMeta(state, data) {
        let key = data.itemKey;
        if (key == "") {
            state.meta = null;
        } else {
            for (let i of state.items) {
                if (i.key == key) {
                    state.meta = i.meta;
                    break;
                }
            }
        }
    },
    setTitle(state, data) {
        state.title = data.title;
    }
}

export const getters = {
}