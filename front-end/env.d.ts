/// <reference types="vite/client" />

interface ImportMetaEnv {
    readonly VITE_AUTH_SERVICE_API: string
  }
  
  interface ImportMeta {
    readonly env: ImportMetaEnv
  }