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
    import { transformModel } from '@/services/transform';

    export default {
        name: "models",
        components: {
            D3Network
        },
        beforeCreate() {
            // setToken(this)
            this.$axios.get('/api/model/K8S/data')
                .then(response => {
                    this.chart = transformModel(response.data);
                })
                .catch(error => console.error(error));
        },
        data() {
            return {
                chart: null
            };
        },
        computed: {
            options() {
                return this.$store.state.graph.options
            }
        },
    }
</script>

<style scoped>

</style>