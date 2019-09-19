<template>
    <div>
        <nuxt-link :to="itemTypeLink"><i class="material-icons"></i><span>{{this.itemTypeLinkLabel}}</span></nuxt-link>
        <d3-network
                v-if="chart != null"
                :net-nodes="chart.nodes"
                :net-links="chart.links"
                v-on:node-click="onNodeClick"
                :options="options"/>
    </div>
</template>

<script>
    import D3Network from 'vue-d3-network';
    import 'vue-d3-network/dist/vue-d3-network.css';
    import {transformGraph} from '@/services/transform';

    export default {
        name: "models",
        components: {
            D3Network
        },
        beforeCreate() {
            // get a list of all models
            this.$axios.get('/api/model')
                .then(modelList => {
                    // defines the info object & adds the models to the info object
                    let info = {
                        models: modelList.data.values,
                        itemTypes: [],
                        linkTypes: [],
                        linkRules: [],
                    };
                    for (let model of modelList.data.values) {
                        // query each model for their content
                        this.$axios.get('/api/model/' + model.key + '/data')
                            .then(modelItem => {
                                // add the model data to the info object
                                info.itemTypes = info.itemTypes.concat(modelItem.data.itemTypes);
                                info.linkTypes = info.linkTypes.concat(modelItem.data.linkTypes);
                                info.linkRules = info.linkRules.concat(modelItem.data.linkRules);
                                // bind to the chart
                                this.chart = transformGraph(info);
                            })
                            .catch(error => console.error(error));
                    }
                })
                .catch(error => console.error(error));
        },
        data() {
            return {
                chart: null,
                itemTypeKey: "",
                itemTypeName: "",
            };
        },
        computed: {
            options() {
                return this.$store.state.graph.options
            },
            itemTypeLink() {
                return "itemType/" + this.itemTypeKey;
            },
            itemTypeLinkLabel() {
                if (this.itemTypeKey == "") {
                    return "";
                } else if (this.itemTypeName.endsWith("y")){
                    return "View " + this.itemTypeName.substr(0, this.itemTypeName.length - 1) + "ies";
                } else {
                    return "View " + this.itemTypeName + "s";
                }
            }
        },
        methods: {
            onNodeClick(event, node) {
                if (node._color == "green") {
                    this.itemTypeKey = node.id;
                    this.itemTypeName = node.name;
                } else {
                    this.itemTypeKey = "";
                    this.itemTypeName = "";
                }
            }
        },
    }
</script>

<style scoped>

</style>