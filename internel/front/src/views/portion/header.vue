<script setup>

import {Setting} from "@element-plus/icons-vue";
</script>

<template>
  <!--   头部   -->
  <div class="toolbar">
    <el-dropdown>
      <el-icon style="margin-right: 8px; margin-top: 1px"
      ><setting
      /></el-icon>
      <template #dropdown>
        <el-dropdown-menu>
          <el-dropdown-item>View</el-dropdown-item>
          <el-dropdown-item>Add</el-dropdown-item>
          <el-dropdown-item>Delete</el-dropdown-item>
        </el-dropdown-menu>
      </template>
    </el-dropdown>
    <span>Tom</span>
  </div>
</template>

<!--<script src="../../assets/css/header.css"></script>-->
<style scoped>
.toolbar {
  text-align: right;
  height: 100%;
  position: relative;
  background-color: var(--el-color-primary-light-7);
  color: var(--el-text-color-primary);
}

.toolbar span.name {
font-size: 30px;
  margin-right: 2%;
}

</style>
