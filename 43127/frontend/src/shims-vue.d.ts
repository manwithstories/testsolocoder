declare module '*.vue' {
  import type { DefineComponent } from 'vue'
  const component: DefineComponent<{}, {}, any>
  export default component
}

declare module 'element-plus/dist/locale/zh-cn'
declare module '@element-plus/icons-vue'
declare module 'echarts'
declare module 'dayjs'
