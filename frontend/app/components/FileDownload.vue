<script lang="ts" setup>
import { ref } from 'vue'
import type { BucketObject } from '~/types';

const api = useNuxtApp().$api;
const props = defineProps<{
  item: BucketObject,
  selectedBucketName: string
}>()

const onDownloadHandle = async () => {
  const data = await api.getBucketObjectFileUrl({
    bucketName: props.selectedBucketName,
    key: props.item.key
  })
  const link = document.createElement('a')
  link.href = data.url
  link.download = props.item.key.split('/').pop() || 'download'
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}
const isDownloading = ref(false)
const onDownloadClick = async () => {
  isDownloading.value = true
  try {
    await onDownloadHandle()
  } catch (error) {
    console.error('Download failed:', error)
  } finally {
    isDownloading.value = false
  }
}

</script>
<template>
  <div class="flex items-center gap-2" :class="{ 'cursor-not-allowed opacity-50': isDownloading }" @click="onDownloadClick">
    <slot></slot>
  </div>
</template>