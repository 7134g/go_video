<template>
  <div class="modal">
    <div class="modal-content">
      <h2>表单</h2>

      <!-- 表单 -->
      <el-form :model="formData" label-width="auto">

        <el-form-item label="任务名字">
          <el-input v-model="formData.name" />
        </el-form-item>

        <el-form-item label="视频类型">
          <el-select v-model="formData.video_type" placeholder="填写视频类型">
            <el-option label="mp4" value="mp4" />
            <el-option label="m3u8" value="m3u8" />
          </el-select>
        </el-form-item>

        <el-form-item label="任务类型">
          <el-select v-model="formData.type" placeholder="填写任务类型">
            <el-option label="url" value="url" />
            <el-option label="curl" value="curl" />
          </el-select>
        </el-form-item>

        <el-form-item label="任务数据">
<!--          <el-input size="large" label-width="300px" v-model="formData.data" type="textarea" />-->

          <el-input
              v-model="formData.data"
              :autosize="{ minRows: 5, maxRows: 10 }"
              type="textarea"
              placeholder="Please input"
          />

        </el-form-item>

        <el-form-item>
          <el-button v-if="insertFlag" type="primary" @click="addTask">提交</el-button>
          <el-button v-if="!insertFlag" type="primary" @click="updateTask">更新</el-button>
          <el-button @click="cancelButton">取消</el-button>
        </el-form-item>

      </el-form>


      <!-- 关闭弹出框按钮 -->
<!--      <span class=close-btn @click=closeModal>X</span>-->
    </div>
  </div>


  <!-- 遮罩层 -->
<!--  <div class=overlay></div>-->


</template>

<script>
import requestFunc from "@/request/task";
import {useCounterStore} from "@/stores/stores";

export default {
  data(){
    return {
      insertFlag: true,

      formData :{
        name: '',
        video_type: '',
        type: '',
        data: '',
      }
    }
  },


  mounted() {
    // 渲染组件时候调用
    this.loadForm()
  },



  methods:{
    loadForm(){
      let store = useCounterStore()
      switch (store.formSwitch) {
        case 1:
          this.insertFlag=true
          this.formData={
              name: '',
              video_type: '',
              type: '',
              data: '',
            }
          break
        case 2:
          this.insertFlag=false
          this.formData={
            id: store.taskData.id,
            name: store.taskData.name,
            video_type: store.taskData.video_type,
            type: store.taskData.type,
            data: store.taskData.data,
          }
          break
      }


    },

    cancelButton(){
      console.log("cancel button")
      let fs = useCounterStore().formSwitch
      switch (fs) {
        case 1:
          break
        case 2:
          break
      }
      this.$emit('close-form');
    },

    addTask(){
      requestFunc.InsertTask(this.formData).then(result => {
        // console.log(result)
        this.$message.success('请求成功');
        this.$emit('close-form');
      }).catch(error => {
        console.log(error)
        this.$message.error('请求失败');
      });

    },

    updateTask(){
      requestFunc.UpdateTask(this.formData).then(result => {
        // console.log(result)
        this.$message.success('请求成功');
      }).catch(error => {
        console.log(error)
        this.$message.error('请求失败');
      });

      this.$emit('close-form');
    },


  }
}


</script>


<style>
.modal {
  position: fixed;
  top: calc(50% - (400px /2));
  left: calc(50% - (1000px /2));
  width: 1000px;
}

.modal-content {
  background-color: #f9f9f9;
  padding:20px;
  border-radius:10px;
  box-shadow:0px 0px .8em rgba(0,0,0,.3);
}


</style>
