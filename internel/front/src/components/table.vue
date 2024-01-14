<template>
  <el-scrollbar >
  <div class="back-icon" @click="closeForm">
    <el-icon >
      <Back />
    </el-icon>
  </div>

  <el-form
      ref="ruleFormRef"
      :model="ruleForm"
      :rules="rules"
      label-width="120px"
      class="demo-ruleForm"
      status-icon
  >



    <el-form-item  v-for="(key, value) in ruleForm" :label="value">
      <el-input v-model="formData[value]"/>
    </el-form-item>

    <el-button @click="printLn"></el-button>

    <el-form-item>
      <el-button type="primary" @click="submitForm(formData)">Submit</el-button>
      <el-button @click="resetForm">Reset</el-button>
    </el-form-item>

  </el-form>
  </el-scrollbar>
</template>


<style>
.back-icon {
  font-size: 25px
}
</style>

<script>

import {useCounterStore} from "../../stores/stores";
export default {
  data() {
    return {
      formData:{}
    }
  },

  methods: {
    printLn(){
      console.log(JSON.stringify(this.formData))
    },

    submitForm(formData){

    },
    resetForm(){
      this.formData = {}
    },

    dataStruct(){

      console.log(JSON.stringify(this.formData))

    },

    closeForm(){
      this.$emit('closeForm-flag');
    },
  },


  mounted() {
    // 在其他方法或是生命周期中也可以调用方法
    this.dataStruct()
  }
}

</script>


