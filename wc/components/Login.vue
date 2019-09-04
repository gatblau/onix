<template>
  <div class="background">
    <main>
      <form @submit.prevent="handleSubmit">
        <splash/>
        <div v-if="showMsg">
          <p>{{message}}</p>
        </div>
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
          <input type="password" class="form-control" placeholder="Password" aria-label="Password" v-model="user.password">
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
        },
        showMsg: true,
        message: '',
      };
    },
    methods: {
      async handleSubmit() {
        // constructs a basic authentication token (base64 encoded)
        const token = btoa(`${this.user.username}:${this.user.password}`)
        // defines the URL of the resource to try and get
        const url = 'api/item?top=1'
        // sets the authentication header
        this.$axios.setHeader('Authorization', `Basic ${token}`)
        try {
          // hides the error message
          this.showMsg = false
          // issues a call for a resource to check if the authentication token worked
          const res = await this.$axios.$get(url)
          // everything went ok so the token was ok
          // set the user name and token in the user store
          this.$store.commit('user/set', this.user.username, token)
          // send the user to the dashboard page
          this.$router.push('dashboard');
        } catch(e) {
          // oops, an error occurred querying the resource with the token
          // sets the error message
          this.message = e
          // shows the error message
          this.showMsg = true
        }
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
