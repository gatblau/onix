<template>
    <div>
        <d3-network
                v-if="chart != null"
                :net-nodes="chart.nodes"
                :net-links="chart.links"
                :options="options"/>
    </div>
</template>

<script>
    import D3Network from 'vue-d3-network';
    import 'vue-d3-network/dist/vue-d3-network.css';
    import axios from 'axios';

    const options = {
        canvas: false,
        force: 3000,
        nodeSize: 50,
        nodeLabels: true,
        linkWidth: 2,
        linkLabels: true,
    };

    export default {
        components: {
            D3Network
        },
        beforeCreate() {
            axios.get('/test.json')
                .then(response => {
                    this.chart = response.data;
                })
                .catch(error => console.error(error));
        },
        data() {
            return {
                options: options,
                chart: null
            };
        }
    };
</script>

<style scoped>
</style>
