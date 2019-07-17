<template>
  <div class="background">
    <main>
      <form @submit.prevent="handleSubmit">
        <splash/>
        <div class="input-group">
          <div class="input-group-prepend">
            <span class="input-group-text"><i class="material-icons">person</i></span>
          </div>
          <input type="text" class="form-control" placeholder="Username" aria-label="Username" v-model="user.username">
        </div>

        <div class="input-group">
          <div class="input-group-prepend">
            <span class="input-group-text"><i class="material-icons">vpn_key</i></span>
          </div>
          <input type="text" class="form-control" placeholder="Password" aria-label="Password">
        </div>

        <div class="form-group">
          <button class="btn btn-primary" type="submit">LOGIN</button>
        </div>
      </form>
    </main>
  </div>
</template>

<script>
  import splash from '@/components/Splash';

  export default {
    name: 'login',
    components: {
      splash
    },
    data() {
      return {
        user: {
          username: '',
          password: ''
        }
      };
    },
    methods: {
      handleSubmit() {

        // this.$emit('clicked', this.user.username);
        this.$router.push('dashboard');
      }
    }
  };
</script>

<style lang="scss" scoped>
  $image-path: "~assets";

  div.background {
    position: fixed;
    height: 100%;
    width: 100%;
    overflow-y: auto;
    background: url("#{$image-path}/bg.png");
    background-size: cover;

    main {
      padding: 0 10px 0 10px;
      position: relative;
      max-width: 500px;
      min-width: 400px;
      overflow-y: auto;
      z-index: 100;
      top: 50%;
      margin: auto;
      transform: translateY(-50%);
      border-radius: 5px;
      border: 1px solid gray;
      background: rgba(255, 255, 255, 1);
      box-shadow: 7px 7px 5px 0 rgba(30, 30, 50, 0.8);
    }
  }

  form {
    button {
      width: 100%;
      margin-top: 5px;
    }

    div.input-group {
      margin-top: 5px;
    }

    div.error {
      margin-top: 2px;
      margin-bottom: 3px;
    }
  }
</style>
