<script lang="ts" setup>
import { ref } from 'vue'
import type { BucketObject } from '~/types';
const props = defineProps<{
  item: BucketObject,
  selectedBucketName: string
}>()

const params = computed(() => ({
  bucket: props.selectedBucketName,
  folder: props.item.key.split('/'),
}))

</script>
<template>
  <NuxtLink :to="{ name: 'files', params }">
    <slot></slot>
  </NuxtLink>
</template>