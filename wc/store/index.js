import Vuex from 'vuex';
import axios from 'axios';

const createStore = () => {
  return new Vuex.Store(
    {
      state: {},
      mutations: {
        setData(state, chartData) {
          state.chartData = chartData;
        }
      },
      actions: {
        getMetaModel(vuexContext, context) {
          return axios
            .get('test.json')
            .then(res => {
              vuexContext.commit('setData', res.data);

              console.log("++++++++++++++++++" + res.data)

            })
            .catch(e => context.error(e));
        }
      },
      getters: {
        chartData(state) {
          return state.chartData;
        }
      }
    }
  );
};

export default createStore;
