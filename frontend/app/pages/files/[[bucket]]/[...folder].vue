<script setup lang="ts">
import { ref, computed } from 'vue'
import { useFileStore } from '@/stores'

definePageMeta({
  name: 'files'
})

const router = useRouter()
const fileStore = useFileStore()
const selectedBucketName = ref(fileStore.selectedBucketName || null)

const currentPath = computed(() => {
  const prefix = fileStore.prefix || ''
  return prefix.split('/').filter(Boolean)
})
const currentFolder = ref(fileStore.prefix || '')
const viewMode = ref(fileStore.viewMode || 'list')
const searchQuery = ref('')
const isUploadDialogOpen = ref(false)
const isCreateFolderDialogOpen = ref(false)

const parentFolderItem = computed(() => {
  if (!currentFolder.value) return null
  const splitPath = currentFolder.value.replace(/\/$/, '').split('/')
  if (splitPath.length <= 1) return { name: '..', key: '', type: 'folder' }
  const parentPath = splitPath.slice(0, -1).join('/') + '/'
  return { name: '..', key: parentPath, type: 'folder' }
})

const filteredItems = computed(() => {
  const filtered = fileStore.getBucketObjectsIncludingText(searchQuery.value)
  if (filtered.length === 0 && searchQuery.value) return []
  return filtered
})

const itemsToRender = computed(() =>
  parentFolderItem.value ? [parentFolderItem.value, ...filteredItems.value] : filteredItems.value
)

const handleViewModeChange = (mode: 'grid' | 'list') => {
  viewMode.value = mode
  fileStore.setViewMode(mode)
}
</script>

<template>
  <div class="min-h-screen flex flex-col bg-background">
    <main v-if="selectedBucketName" class="flex-1 p-6">
      <Toolbar
        :searchQuery="searchQuery"
        :viewMode="viewMode"
        :isUploadDialogOpen="isUploadDialogOpen"
        :isCreateFolderDialogOpen="isCreateFolderDialogOpen"
        @update:viewMode="handleViewModeChange"
        @closeCreateFolderDialog="isCreateFolderDialogOpen = false"
        @closeUploadDialog="isUploadDialogOpen = false"
      />
      <div class="mb-4">
        <FolderCrumbs />
      </div>
      <FileList
        v-if="viewMode === 'list'"
        :items="itemsToRender"
        :selectedBucketName="selectedBucketName"
      />
      <FileGrid
        v-else
        :items="itemsToRender"
        :selectedBucketName="selectedBucketName"
      />
      <EmptyState v-if="filteredItems.length === 0" :searchQuery="searchQuery" />
    </main>
    <BucketPicker v-else />
  </div>
</template>

<style scoped>
/* Add your styles here */
</style>