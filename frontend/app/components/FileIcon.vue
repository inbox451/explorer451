<script lang="ts" setup>
import type { BucketObject } from '~/types';
import {
  File, Folder, FolderOutput, FileText, ImageIcon, Video, Music, Archive,
} from 'lucide-vue-next'
const props = defineProps<{
  item: BucketObject,
  size?: string,
}>()


const fileIcon = computed(() => {
  const { type, name } = props.item
  if (type === 'folder') return (name == '..' ? FolderOutput : Folder)
  switch (props.item.contentType) {
    case 'image':
    case 'image/gif':
    case 'image/jpeg':
    case 'image/png':
    case 'image/svg+xml':
    case 'image/webp': return ImageIcon
    case 'video': return Video
    case 'audio': return Music
    case 'archive': return Archive
    case 'document': return FileText
    default: return File
  }
})

const size = computed(() => props.size || 'h-4 w-4')
const className = computed(() => props.item.type === 'folder' ? 'text-blue-500' : 'text-muted-foreground')

</script>
<template>
  <component :is="fileIcon" :class="[size, className]" :title="props.item.key" />
</template>