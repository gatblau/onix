<template>
    <div class="table-wrapper-scroll-y my-custom-scrollbar">
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
                    >view</button></td>
            </tr>
            </tbody>
        </table>
    </div>
</template>

<script>
    export default {
        name: "itemTypeKey",
        data() {
            return {
                items: [],
            }
        },
        beforeCreate() {
            this.$axios.get('/api/item?type=' + this.$route.params.itemTypeKey)
                .then((items) => {
                    this.items = items.data.values;
            }).catch(error => console.error(error));
        },
        methods: {
            onItemClick(data){
                console.log(data.target.value);
            }
        }
    }
</script>

<style scoped>

</style>