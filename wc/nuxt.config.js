/*
 *   Onix Web Console - Copyright (c) 2019 by www.gatblau.org
 *
 *   Licensed under the Apache License, Version 2.0 (the "License");
 *   you may not use this file except in compliance with the License.
 *   You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
 *   Unless required by applicable law or agreed to in writing, software distributed under
 *   the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 *   either express or implied.
 *   See the License for the specific language governing permissions and limitations under the License.
 *
 *   Contributors to this project, hereby assign copyright in this code to the project,
 *   to be licensed under the same terms as the rest of the code.
*/
const pkg = require('./package');

module.exports = {
  // mode: 'spa',

  server: {
    // allow connections outside of the host machine (default is 'localhost')
    host: '0.0.0.0',
  },

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
    {src: '~/plugins/flag-icon-css'},
    {src: '~/plugins/vuetify'},
  ],

  modules: [
    '@nuxtjs/axios', // Doc: https://axios.nuxtjs.org/usage
    'bootstrap-vue/nuxt', // Doc: https://bootstrap-vue.js.org/docs/
    '@nuxtjs/proxy', // Doc: https://github.com/nuxt-community/proxy-module
  ],

  axios: {
    // See https://github.com/nuxt-community/axios-module#options
    proxy: true
  },

  build: {
    /*
    ** You can extend webpack config here
    */
    extend(config, ctx) {
      // NuxtJS debugging support
      // eval-source-map: a SourceMap that matchers exactly to the line number and this help to debug the NuxtJS app in the client
      // inline-source-map: help to debug the NuxtJS app in the server
      config.devtool = ctx.isClient ? 'eval-source-map' : 'inline-source-map'
    }
  },

  env: {
    ox_wapi_uri: process.env.WC_OX_WAPI_URI || 'http://localhost:8080',
    ox_wapi_auth_mode: process.env.WC_OX_WAPI_AUTH_MODE || 'basic',
  },

  proxy: {
    // proxy all calls through to /api to the onix wapi uri
    '/api': {
      target: process.env.WC_OX_WAPI_URI || 'http://localhost:8080',
      pathRewrite: { '^/api/': '' }
    }
  }
};
