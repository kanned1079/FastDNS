// import { StrictMode,  } from 'react'
import { createRoot } from 'react-dom/client'
import {ConfigProvider, theme} from "antd"
import './index.css'
import App from './App.tsx'

createRoot(document.getElementById('root')!).render(
  <>

      <ConfigProvider
          theme={{
              // 1. 单独使用暗色算法
              algorithm: theme.defaultAlgorithm,

              // 2. 组合使用暗色算法与紧凑算法
              // algorithm: [theme.darkAlgorithm, theme.compactAlgorithm],
          }}
      >
          <App />

      </ConfigProvider>

  </>
)
