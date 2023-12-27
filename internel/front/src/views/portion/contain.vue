<script setup>

</script>

<template>
  <!--   内容   -->

</template>

<style scoped>


</style>


<script>

import requestFunc from "../../request/table";
import {useCounterStore} from '../../stores/stores';

export default {
  data() {
    return {
      showStruct: true,
      showDataList: false,
      showDataForm: false,
      tableStruct: [],

    }
  },

  methods: {


  },
  mounted() {
    // 在其他方法或是生命周期中也可以调用方法
    this.getTables()
  }
}


</script>
