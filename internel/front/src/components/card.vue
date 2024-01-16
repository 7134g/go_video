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

  <el-card style="width: 90%" class="box-card">
    <el-row>
      <el-text class="w-150px mb-2" truncated>
        详细
      </el-text>
    </el-row>


    <el-row>
      <el-text class="mx-1">
        <p v-for="(segment, index) in messageSegments" :key="index">{{ segment }}</p>
      </el-text>
    </el-row>

  </el-card>
</template>


<script>
export  default {
  props: ['message'],
  emits: ['closeDataList'],

  data() {
    return {
      detail:"",
    }
  },

  computed: {
    messageSegments() {
      return this.message.split("\\")
    }
  },

  methods: {
    closeDetail(){
      this.$emit('closeDataList');
    },
  },
}

</script>

<style scoped>
.box-card {
  width: 480px;
}
</style>
