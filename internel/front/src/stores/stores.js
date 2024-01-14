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
      dataPage: {
        where: {
          type: "all",
          video_type: null,
        },
        page: 1,
        size: 10,
      },
    }
  },
  // 也可以定义为
  // state: () => ({ count: 0 })
  actions: {
    getDataPage() {
      return this.dataPage
    },

    setTaskType(type) {
      // console.log("aaaa", type)
      this.dataPage.where.type = type
    },
    getTaskType() {
      return this.dataPage.where.type
    },

    setVideoType(video_type) {
      this.dataPage.where.video_type = video_type
    },
    getVideoType() {
      return this.dataPage.where.video_type
    },



  },
})
