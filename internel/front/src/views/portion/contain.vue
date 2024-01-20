<script setup>

import TaskList from "@/components/taskList.vue";
</script>

<template>
  <!--   内容   -->
<!--  <taskList v-if="showTable"></taskList>-->
  <taskList @render-task-list="resetTaskList" :key="redirect" v-if="showTable"></taskList>
</template>

<style scoped>


</style>


<script>


import {useCounterStore} from "@/stores/stores";

export default {
  data() {
    return {
      redirect: false,

      showTable:true,
    }
  },

  methods: {
    resetTaskList(type){
      console.log("redirect task list", type)
      const store = useCounterStore();
      store.setTaskType(type);
      this.redirect = !this.redirect
    },

  },

}


</script>
