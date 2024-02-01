<script setup>
import {ref} from "vue";

const small = ref(false)
const background = ref(false)
const disabled = ref(false)

import formData from "@/components/formData.vue";
</script>

<template>
  <div v-if="showDataList" >
    <el-table :data="tableData" style="width: 100%" height="100%">
      <el-table-column type="index" label="序号" width="80" />
<!--      <el-table-column fixed prop="id" label="任务号" width="80" />-->
      <el-table-column prop="name" label="任务名称" width="200" />
      <el-table-column prop="video_type" label="视频类型" width="80" />
      <el-table-column prop="type" label="任务类型" width="80" />
      <el-table-column prop="status" label="任务状态" width="80" />
      <el-table-column label="任务进度" width="100">
        <template #default="scope">
          {{ scope.row.score }}%
        </template>
      </el-table-column>
      <el-table-column prop="data" show-overflow-tooltip label="任务数据" width="200" />
      <el-table-column fixed="right" label="操作" width="180">
        <template #default="{ row }">
          <el-button link type="primary" size="large" @click="openCard(row.data)">查看数据</el-button>
          <el-button link type="primary" size="large" @click="editForm(row)">编辑</el-button>
          <el-button link type="primary" size="large" @click="deleteData(row.id)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>


    <div class="demo-pagination-block">
      <el-pagination
          v-model:current-page="currentPage"
          :page-size="currentSize"
          :small="small"
          :disabled="disabled"
          :background="background"
          layout="total, prev, pager, next"
          :total="total"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
      />
    </div>
  </div>



  <Card v-if="showCard" :message="detail" @closeDataList="closeCard"></Card>

  <formData v-if="showDataForm" @close-form="closeForm"></formData>

</template>

<style scoped>

</style>

<script>

import requestFunc from "@/request/task";
import {useCounterStore} from '@/stores/stores';
import Card from "@/components/card.vue";

export default {
  components: { Card },
  emits:['render-task-list'],

  data() {
    return {
      showDataList: true,

      showCard: false,
      detail: "test",

      showDataForm: false,


      currentPage: 0,
      currentSize: 0,
      total: 0,
      // tableData: [{
      //   "id": 1,
      //   "name": "test1",
      //   "video_type": "mp4",
      //   "type": "url",
      //   "data": "http://clips.vorwaerts-gmbh.de/big_buck_bunny.mp4",
      //   "status": 0
      // }],
      tableData:[],

    }
  },

  methods: {
    getTables(dp) {
      dp = useCounterStore().getDataPage()
      // console.log("post", dp)
      requestFunc.GetTaskList(dp).then(result => {
        this.tableData = result.list
        this.total = result.total
        // console.log(this.tableData)
        // console.log(JSON.stringify(this.tableData))
        // this.$message.success('请求成功');
      }).catch(error => {
        console.log(error)
        this.$message.error('请求失败');
      });
    },

    handleCurrentChange(val) {
      let dp = useCounterStore().getDataPage()
      dp.page = val
      // console.log(dp)
      this.getTables(dp)
    },

    handleSizeChange(){
      let dp = useCounterStore().getDataPage()
      dp.size = this.currentSize
      console.log("page-size 改变时触发")
    },

    openCard(detail) {
      this.showDataList = false
      this.showCard = true
      this.detail = detail
    },
    closeCard() {
      this.showDataList = !this.showDataList // 关闭表结构列表，显示详细数据内容
      this.showCard = !this.showCard
    },

    editForm(task){
      console.log("edit task", task)
      useCounterStore().setFormSwitch(2)
      useCounterStore().setTaskData(task)
      this.showDataList = false
      this.showDataForm = true
    },

    closeForm(){
      console.log("close form")
      this.showDataList = !this.showDataList
      this.showDataForm = !this.showDataForm
    },

    deleteData(id){
      requestFunc.DeleteTask(id).then(_ => {
        this.$message.success('请求成功');
      }).catch(error => {
        console.log(error)
        this.$message.error('请求失败');
      });

      this.render()
    },

    render(){
      this.$emit('render-task-list', useCounterStore().getTaskType());
    },


  },
  mounted() {
    let dp = useCounterStore().getDataPage()
    this.currentSize = dp.size
    this.currentPage = dp.page
    // 在其他方法或是生命周期中也可以调用方法
    this.getTables(useCounterStore().getDataPage())
  }
}


</script>
