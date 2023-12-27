// import { createPinia } from 'pinia';
//
// const pinia = createPinia();
//
// export const useStore = pinia.createStore({
//   state: () => ({
//     _type: '',
//   }),
//   actions: {
//     setDbType(dbType) {
//       this.db_type = dbType;
//     },
//   },
// });



import { defineStore } from 'pinia'

export const useCounterStore = defineStore('counter', {
  state: () => {
    return {

    }
  },
  // 也可以定义为
  // state: () => ({ count: 0 })
  actions: {


  },
})
