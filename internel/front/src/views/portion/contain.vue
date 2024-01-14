<script setup>

import TaskList from "@/components/taskList.vue";
</script>

<template>
  <!--   内容   -->
  <taskList v-if="showTable"></taskList>
</template>

<style scoped>


</style>


<script>

import requestFunc from "@/request/table";
import {useCounterStore} from '@/stores/stores';

export default {
  data() {
    return {
      showTable:true,
    }
  },

  methods: {


  },

}


</script>
