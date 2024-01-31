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
      formSwitch:1,
      taskStatus: false,
      taskData:{
        id:0,
        name: '',
        video_type: '',
        type: '',
        data: '',
      },

      dataPage: {
        where: {
          type: "all",
          video_type: null,
        },
        page: 1,
        size: 20,
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
    getTaskType(){
      return this.dataPage.where.type
    },



    setTaskStatus(status){
      this.taskStatus = status
    },

    setTaskData(task) {
      this.taskData = task
    },

    setFormSwitch(value) {
      this.formSwitch = value
    },

  },
})
