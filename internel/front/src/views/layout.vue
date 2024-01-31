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
        <Header @open-form="openForm"></Header>
      </el-header>

      <el-main>
        <Contain v-if="containHidden" :key="redirect"></Contain>
        <formData v-if="formFlag" @close-form="closeForm"></formData>
      </el-main>
    </el-container>
  </el-container>
</template>

<script lang="ts">
import { useCounterStore } from '../stores/stores.js';
import formData from "../components/formData.vue";


export default {
  components: { formData },

  data() {
    return {
      redirect: false, // 重新加载

      containHidden: true,// 隐藏列表
      formFlag: false, // 显示表单
    }
  },
  methods: {
    resetContain(type) {
      console.log("render-contain", type)
      // 触发子组件B的重新渲染
      const store = useCounterStore();
      store.setTaskType(type);
      this.redirect = !this.redirect;

      this.formFlag = false
      this.containHidden = true
    },

    openForm(){
      console.log("layout open form")
      this.formFlag = !this.formFlag
      this.containHidden = !this.containHidden;
    },

    closeForm(){
      console.log("layout close form")
      this.formFlag = false
      this.containHidden = true
      this.resetContain(useCounterStore().getTaskType())
    },
  }
}
</script>

<style scoped>
.layout-container-demo .el-header {
  height: 100px;
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
