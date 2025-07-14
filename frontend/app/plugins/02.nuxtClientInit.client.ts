export default defineNuxtPlugin(() => {
  if (import.meta.client) {
    const bucketStore = useBucketStore()
    bucketStore.getBuckets()
    // bucketStore.getBucketObjects()
  }
})
