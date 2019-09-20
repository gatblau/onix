<template>
    <div class="table-wrapper-scroll-y my-custom-scrollbar">
        <table>
            <td><button type="button" class="btn btn-outline-primary" onclick="window.history.back()"><i class="material-icons">keyboard_arrow_up</i></button></td>
            <td><h3>Children of <i>{{title}}</i></h3></td>
        </table>
        <table class="table table-bordered table-striped table-hover mb-0">
            <thead class="thead-light">
            <tr>
                <th scope="col">Key</th>
                <th scope="col">Name</th>
                <th scope="col">Description</th>
                <th scope="col">Type</th>
                <th scope="col">Action</th>
            </tr>
            </thead>
            <tbody>
            <tr v-for="item in items">
                <td>{{ item.key }}</td>
                <td>{{ item.name }}</td>
                <td>{{ item.description }}</td>
                <td>{{ item.typeName }}</td>
                <td>
                    <button
                            type="button"
                            class="btn btn-outline-primary"
                            v-on:click="onItemIn"
                            :value="item.key"
                    >in</button>
                    <button
                            type="button"
                            class="btn btn-outline-primary"
                            v-on:click="onItemView"
                            :value="item.key"
                    >view</button>
                </td>
            </tr>
            </tbody>
        </table>
        <div v-if="showMeta">
            <h4>Metadata</h4>
            <div class="overflow-auto" style="max-height: 250px;">
                <code>{{meta}}</code>
            </div>
        </div>
    </div>
</template>

<script>
    export default {
        name: "itemKey",
        methods: {
            onItemIn(data) {
                this.$router.push('../item/' + data.target.value);
            },
            onItemView(data) {
                this.$store.commit('graph/setMeta', { itemKey: data.target.value, app: this });
            }
        },
        computed: {
            items() {
                return this.$store.state.graph.items;
            },
            meta() {
                return JSON.stringify(this.$store.state.graph.meta);
            },
            showMeta() {
                return this.$store.state.graph.meta != null;
            },
            title(){
                return this.$route.params.itemKey;
            }
        },
        async asyncData ({ params, $axios, app, store }) {
            store.commit('graph/setMeta', { itemKey: "", app: this });
            store.commit('graph/setTitle', { title: params.itemKey, app:this });
            $axios.get('/api/item/' + params.itemKey + "/children")
                .then((result) => {
                    store.commit('graph/setItems', { items: result.data.values, app: app });
                }).catch(error => console.error(error));
        }
    }
</script>

<style scoped>

</style>