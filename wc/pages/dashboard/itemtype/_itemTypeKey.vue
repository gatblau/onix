<template>
    <div class="table-wrapper-scroll-y my-custom-scrollbar">
        <!--<button type="button" class="btn btn-outline-primary" onclick="window.history.back()"><- Back</button>-->
        <table class="table table-bordered table-striped table-hover mb-0">
            <thead class="thead-light">
            <tr>
                <th scope="col">Key</th>
                <th scope="col">Name</th>
                <th scope="col">Description</th>
                <th scope="col">Action</th>
            </tr>
            </thead>
            <tbody>
            <tr v-for="item in items">
                <td>{{ item.key }}</td>
                <td>{{ item.name }}</td>
                <td>{{ item.description }}</td>
                <td>
                    <button
                        type="button"
                        class="btn btn-primary"
                        v-on:click="onItemClick"
                        :value="item.key"
                    >view</button>
                </td>
            </tr>
            </tbody>
        </table>
    </div>
</template>

<script>
    export default {
        name: "itemTypeKey",
        methods: {
            onItemClick(data){
                this.$router.push('../item/' + data.target.value);
                this.$forceUpdate();
            }
        },
        computed: {
            items () {
                return this.$store.state.graph.items;
            }
        },
        // watchQuery: true,
        async asyncData ({ params, $axios, app, store }) {
            $axios.get('/api/item?type=' + params.itemTypeKey)
                .then((result) => {
                    store.commit('graph/setItems', { items: result.data.values, app: app });
                }).catch(error => console.error(error));
        }
    }
</script>

<style scoped>

</style>