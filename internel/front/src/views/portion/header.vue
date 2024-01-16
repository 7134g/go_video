<script setup>

</script>

<template>
  <!--   头部   -->
  <div class="toolbar">
    <el-button size="large" type="success" @click="runTask">Run</el-button>
    <el-button size="large" type="danger" @click="runTask">Stop</el-button>

  </div>
</template>

<script>

import requestFunc from "@/request/task";

export default {
  data() {
    return {

    }
  },

  methods: {
    runTask(){
      console.log("run.......")

      requestFunc.RunTask().then(result => {
        // console.log(JSON.stringify(result))
        if (result.data.message !== "") {
          this.$message.success(result.data.message);
        }

      }).catch(error => {
        console.log(error)
        this.$message.error('请求失败');
      });


    }

  }
}

</script>

<style scoped>
.toolbar {
  height: 100%;
  position: relative;
  background-color: var(--el-color-primary-light-7);
  color: var(--el-text-color-primary);
}


</style>
