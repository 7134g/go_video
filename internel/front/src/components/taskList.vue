<script setup>
import {ref} from "vue";

const small = ref(false)
const background = ref(false)
const disabled = ref(false)

import Card from "@/components/card.vue";
import Insert from "@/components/insert.vue";
</script>

<template>
  <div v-if="showDataList" >
    <el-table :data="tableData" style="width: 100%" height="250">
      <el-table-column fixed prop="id" label="任务号" width="150" />
      <el-table-column prop="name" label="任务名称" width="120" />
      <el-table-column prop="video_type" label="视频类型" width="120" />
      <el-table-column prop="type" label="任务类型" width="100" />
      <el-table-column prop="status" label="任务状态" width="100" />
      <el-table-column prop="data" label="任务数据" width="800" />
      <el-table-column fixed="right" label="操作" width="120">
        <template #default="{ row }">
          <el-button link type="primary" size="large" @click="changeCard(row.data)">Detail</el-button>
          <el-button link type="primary" size="large" @click="openForm(row)">Edit</el-button>
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



  <Card v-if="showCard" :message="detail" @closeDataList-flag="changeCard"></Card>

  <Insert v-if="showDataForm" @closeForm-flag="closeForm"></Insert>

</template>

<style scoped>

</style>

<script>

import requestFunc from "@/request/table";
import {useCounterStore} from '@/stores/stores';

export default {

  data() {
    return {
      showDataList: true,

      showCard: false,
      detail: "test",

      showDataForm: false,


      currentPage: 1,
      currentSize: 2,
      total: 10,
      tableData: [{
        "id": 1,
        "name": "test1",
        "video_type": "mp4",
        "type": "url",
        "data": "http://clips.vorwaerts-gmbh.de/big_buck_bunny.mp4",
        "status": 0
      }],

    }
  },

  methods: {
    getTables(dp) {
      dp = useCounterStore().getDataPage()
      console.log("post", dp)
      requestFunc.GetTaskList(dp).then(result => {
        this.tableData = result.list
        this.total = result.total
        // console.log(JSON.stringify(this.tableData))
        this.$message.success('请求成功');
      }).catch(error => {
        console.log(error)
        this.$message.error('请求失败');
      });
    },

    handleCurrentChange(val) {
      let dp = useCounterStore().getDataPage()
      dp.page = val
      console.log(dp)
      this.getTables(dp)
    },

    handleSizeChange(){
      console.log("page-size 改变时触发")
    },


    changeCard(detail) {
      this.showDataList = !this.showDataList // 关闭表结构列表，显示详细数据内容
      this.showCard = !this.showCard
      this.detail = detail
      // useCounterStore().setTableName(id)
    },

    openForm(data){
      this.showDataList = false
      this.showDataForm = true
      // useCounterStore().setDataStruct(id)
    },
    closeForm(){
      this.showDataList = !this.showDataList
      this.showDataForm = !this.showDataForm
    },

  },
  mounted() {
    // 在其他方法或是生命周期中也可以调用方法
    this.getTables(useCounterStore().getDataPage())
  }
}


</script>
