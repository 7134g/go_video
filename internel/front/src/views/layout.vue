<script lang="ts" setup>

import Aside from "./portion/aside.vue";
import Header from "./portion/header.vue";
import Contain from "./portion/contain.vue";

</script>

<template>
  <el-container class="layout-container-demo" style="height: 800px">
    <el-aside >
      <Aside @render-contain="resetContain"></Aside>
    </el-aside>

    <el-container>
      <el-header style="">
        <Header></Header>
      </el-header>

      <el-main>
        <Contain :key="redirect"></Contain>
      </el-main>
    </el-container>
  </el-container>
</template>

<script lang="ts">
import { useCounterStore } from '../stores/stores.js';


export default {
  data() {
    return {
      redirect: false
    }
  },
  methods: {
    resetContain(type) {
      // 触发子组件B的重新渲染
      const store = useCounterStore();
      // console.log("================", type)
      store.setTaskType(type);
      // console.log("----------------", type)
      // console.log("getTaskType",store.getTaskType())
      // console.log("handleCustomEvent", type)
      this.redirect = !this.redirect;
    }
  }
}
</script>

<style scoped>
.layout-container-demo .el-header {
  height: 100px;
  text-align: center;
  font-size: 15px;
  position: relative;
  background-color: var(--el-color-primary-light-7);
  color: var(--el-text-color-primary);
}
.layout-container-demo .el-aside {
  color: var(--el-text-color-primary);
}
.layout-container-demo .el-menu {
  border-right: none;
}
.layout-container-demo .el-main {
  padding: 0;
}
.layout-container-demo .toolbar {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  right: 20px;
}
</style>
