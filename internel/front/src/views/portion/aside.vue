<script setup>

</script>

<template>
  <!--  侧边栏  -->
<!--  <el-aside class="asideStyle" width="200px" style="height: 90%">-->
    <el-scrollbar>
      <el-menu :default-openeds="['1', '2', '3']">

<!--        <el-sub-menu index="1">-->
<!--          <template #title>-->
<!--            <el-icon></el-icon>-->
<!--            数据库-->
<!--          </template>-->
<!--          <el-menu-item-group>-->
<!--            <el-menu-item @click="setDbType('mysql')" index="1-1">mysql</el-menu-item>-->
<!--            <el-menu-item @click="setDbType('mongo')" index="1-2">mongo</el-menu-item>-->
<!--            <el-menu-item @click="setDbType('sqlite')" index="1-3">sqlite</el-menu-item>-->
<!--          </el-menu-item-group>-->
<!--        </el-sub-menu>-->


        <el-sub-menu index="2">
          <template #title>
            <el-icon></el-icon>任务列表
          </template>

            <el-menu-item @click="setTaskType('all')" index="2-1-1">all</el-menu-item>
            <el-menu-item @click="setTaskType('url')" index="2-1-2">url</el-menu-item>
            <el-menu-item @click="setTaskType('curl')" index="2-1-3">curl</el-menu-item>
            <el-menu-item @click="setTaskType('proxy')" index="2-1-4">proxy</el-menu-item>

        </el-sub-menu>

        <el-sub-menu index="3">
          <template #title>
            <el-icon></el-icon>信息
          </template>

          <el-menu-item index="3-1-1">错误任务</el-menu-item>
          <el-menu-item index="3-1-2">日志</el-menu-item>
<!--          <el-menu-item index="3-1-3">Option 3</el-menu-item>-->

        </el-sub-menu>


      </el-menu>
    </el-scrollbar>
<!--  </el-aside>-->
</template>

<style scoped>
.asideStyle {
  color: var(--el-text-color-primary);
  background: var(--el-color-primary-light-8);
}


el-menu {
  border-right: none;
}


</style>


<script>
import { Menu as IconMenu, Message, Setting } from '@element-plus/icons-vue'
export default {
  methods: {
    setTaskType(type){
      this.$emit('render-contain', type);
    }
  },

}


</script>
