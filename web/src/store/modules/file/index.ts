export interface File {
  path: string
}

export const useFileStore = defineStore('file', {
  state: (): File => {
    return {
      path: '/'
    }
  },
  actions: {
    set(info: File) {
      this.path = info.path
    }
  },
  persist: true
})
