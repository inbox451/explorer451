export default defineNuxtPlugin(() => {
  if (import.meta.client) {
    const fileStore = useFileStore()
    fileStore.getBuckets()
    // fileStore.getBucketObjects()
  }
})
