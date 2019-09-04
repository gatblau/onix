<template lang="html">
  <div class="d-flex">

    <div class="bg-light border-right sidebar">
      <img class="logo" src="~assets/logo.png" alt="Onix"/>
      <hr>
      <nav>
        <div @click="title='Dashboard'">
          <nuxt-link to="/dashboard"><i class="material-icons">home</i><span>Dashboard</span></nuxt-link>
        </div>
        <div @click="title='Meta Model'">
          <nuxt-link to="/dashboard/metaModel"><i class="material-icons">device_hub</i><span>Meta-Model</span></nuxt-link>
        </div>
      </nav>
      <hr>
    </div>

    <div class="page-content">
      <b-navbar toggleable="lg" type="dark" variant="dark">
        <b-navbar-brand>{{ title }}</b-navbar-brand>
        <b-navbar-toggle target="nav-collapse"></b-navbar-toggle>
        <b-collapse id="nav-collapse" is-nav>
          <b-navbar-nav class="ml-auto">
            <b-nav-item-dropdown text="Lang">
              <template slot="button-content"><i class="material-icons">translate</i></template>
              <b-dropdown-item href="#"><i class="flag-icon flag-icon-gb"></i> English</b-dropdown-item>
              <b-dropdown-item href="#"><i class="flag-icon flag-icon-es"></i> Español</b-dropdown-item>
              <b-dropdown-item href="#"><i class="flag-icon flag-icon-cn"></i> 中文</b-dropdown-item>
            </b-nav-item-dropdown>

            <b-nav-item-dropdown>
              <template slot="button-content"><i class="material-icons">person</i> {{loggedUser}}</template>
              <b-dropdown-item href="#">Profile</b-dropdown-item>
              <b-dropdown-item href="#">Sign Out</b-dropdown-item>
            </b-nav-item-dropdown>
          </b-navbar-nav>
        </b-collapse>
      </b-navbar>

      <div class="container-fluid">
        <nuxt-child/>
      </div>
    </div>
  </div>
</template>

<script>
  export default {
    data: () => ({
      title: 'Dashboard',
    }),
    computed: {
      loggedUser () {
        return this.$store.state.user.username
      }
    },
  };
</script>

<style lang="scss" scoped>
  .logo {
    width: 80%;
    margin: 1rem 1rem 0 1rem;
  }

  .sidebar {
    min-height: 100vh;
    transition: margin .25s ease-out;
    width: 15rem;

    div {
      padding: .5rem;

      a {
        span {
          vertical-align: top;
          margin-left: 5px;
        }
      }
    }
  }

  .page-content {
    min-width: 100vw;
  }

  @media (min-width: 768px) {
    .sidebar {
      margin-left: 0;
    }

    .page-content {
      min-width: 0;
      width: 100%;
    }
  }
</style>
