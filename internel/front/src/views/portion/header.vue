<script setup>

</script>

<template>
  <div class="header_layout">
    <!--   头部   -->
    <div :key="taskButtonFlag" class="toolbar">
      <el-button size="large" type="primary" @click="openTask">添加任务</el-button>
      <el-button size="large" v-if="taskButtonFlag" type="success" @click="runTask">启动</el-button>
      <el-button size="large" v-if="!taskButtonFlag" type="danger" @click="runTask">停止</el-button>
    </div>

    <div class="proxy_bar">
<!--      <el-input class="proxy_input" v-model="proxy" placeholder="web被动代理地址" clearable />-->
      <input v-model="proxy" placeholder="web被动代理地址" type="text" />
    </div>

  </div>


</template>

<script>
import {useCounterStore} from '@/stores/stores';
import requestFunc from "@/request/task";

export default {
  data() {
    return {
      insertFlag: false,
      taskButtonFlag: true,
      status: false,
      proxy:"",

      data:{}
    }
  },

  mounted(){
    this.taskStatus()
  },

  methods: {
    taskStatus(){
      requestFunc.GetProgramStatus().then(result => {
        this.status=result.data.status
        this.proxy=result.data.web_proxy
        useCounterStore().setTaskStatus(result.data.status)
        if (this.status){
          this.taskButtonFlag = !this.status
        }

      }).catch(error => {
        console.log(error)
        this.$message.error('请求失败');
      });

    },


    openTask(){
      useCounterStore().setFormSwitch(1)
      console.log("open task......")
      this.$emit('open-form');
    },

    runTask(){
      console.log("run.......")

      requestFunc.RunTask().then(result => {
        // console.log(JSON.stringify(result))
        if (result.data.message !== "") {
          this.$message.success(result.data.message);
        }
        this.taskStatus()
      }).catch(error => {
        console.log(error)
        this.$message.error('请求失败');
      });


    }

  }
}

</script>

<style scoped>

.header_layout {
  height: 100%;
  width: 100%;
  display: flex;
}

.toolbar {
  display: flex; /* 开启 Flexbox 布局 */
  height: 100%;
  width: 40%;
  text-align: left;
  align-items: center; /* 垂直居中 */
  background-color: var(--el-color-primary-light-7);
  color: var(--el-text-color-primary);
}

.proxy_bar {
  display: flex; /* 开启 Flexbox 布局 */
  height: 100%;
  flex-wrap: wrap;
  text-align: left;
  //justify-content: center; /* 水平居中 */
  align-items: center; /* 垂直居中 */
  width: 50%;
}

.proxy_bar input[type="text"] {
  height: 30px;
  border: none;
  background-color: transparent;
  color: #4e8e2f;
  font-size: 30px;
}

</style>
