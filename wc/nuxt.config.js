const pkg = require('./package');

module.exports = {
  mode: 'spa',

  head: {
    title: pkg.name,
    meta: [
      {charset: 'utf-8'},
      {name: 'viewport', content: 'width=device-width, initial-scale=1'},
      {hid: 'description', name: 'description', content: pkg.description}
    ],
    link: [
      {rel: 'icon', type: 'image/png', href: '/favicon.png'}
    ]
  },

  /*
  ** Customize the progress-bar color
  */
  loading: {color: 'red'},

  /*
  ** Global CSS
  */
  css: [],

  plugins: [
    {src: '~/plugins/material-icons'},
    {src: '~/plugins/flag-icon-css'}
  ],

  modules: [
    '@nuxtjs/axios', // Doc: https://axios.nuxtjs.org/usage
    'bootstrap-vue/nuxt' // Doc: https://bootstrap-vue.js.org/docs/
  ],

  axios: {
    // See https://github.com/nuxt-community/axios-module#options
  },

  build: {
    /*
    ** You can extend webpack config here
    */
    extend(config, ctx) {
    }
  }
};
