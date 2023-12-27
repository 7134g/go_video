
import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './app.vue'
import router from './router/route'
import { ElMessage } from 'element-plus'
import 'element-plus/dist/index.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.config.globalProperties.$message = ElMessage
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
    app.component(key, component)
}
app.mount('#app')
