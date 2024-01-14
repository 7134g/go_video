<script setup>

</script>

<template>
  <div class="back-icon" @click="closeDetail">
    <el-button size="large" circle>
      <el-icon >
        <Back />
      </el-icon>
    </el-button>

  </div>

  <el-card class="box-card">
    <el-row>
      <el-text class="w-150px mb-2" truncated>
        详细
      </el-text>
    </el-row>


    <el-row>
      <el-text class="mx-1">
        {{ detail }}
      </el-text>
    </el-row>

  </el-card>
</template>


<script>
export  default {
  props: ['detail'],

  data() {
    return {
    }
  },


  methods: {
    closeDetail(){

    },
  },
}

</script>

<style scoped>
.text {
  font-size: 14px;
}

.item {
  padding: 18px 0;
}

.box-card {
  width: 480px;
}
</style>
