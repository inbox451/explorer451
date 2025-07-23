<script setup lang="ts">
import { ref, computed } from 'vue'
import {
  ChevronRight, File, Folder, Upload, Download, Trash2, Search, Grid3X3, List, Plus, MoreHorizontal,
  Database
} from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardContent } from '@/components/ui/card'
import { DropdownMenu, DropdownMenuTrigger , DropdownMenuContent, DropdownMenuItem, DropdownMenuSeparator } from '@/components/ui/dropdown-menu'
import { Dialog, DialogTrigger, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Label } from '@/components/ui/label'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { useBucketStore } from '@/stores'
import { Title } from '#components'

definePageMeta({
  name: 'files'
})

const router = useRouter()
const bucketStore = useBucketStore()
const selectedBucketName = ref(bucketStore.selectedBucketName || null)

const currentPath = computed(() => {
  const prefix = bucketStore.prefix || ''
  return prefix.split('/').filter(Boolean)
})
const currentFolder = ref(bucketStore.prefix || '')
const viewMode = ref<'list' | 'grid'>('list')
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
  const filtered = bucketStore.getBucketObjectsIncludingText(searchQuery.value)
  if (filtered.length === 0 && searchQuery.value) return []
  return filtered
})

const itemsToRender = computed(() =>
  parentFolderItem.value ? [parentFolderItem.value, ...filteredItems.value] : filteredItems.value
)
</script>

<template>
  <div class="min-h-screen flex flex-col bg-background">
    <main v-if="selectedBucketName" class="flex-1 p-6">
      <Toolbar
        :searchQuery="searchQuery"
        :viewMode="viewMode"
        :isUploadDialogOpen="isUploadDialogOpen"
        :isCreateFolderDialogOpen="isCreateFolderDialogOpen"
        @update:viewMode="viewMode = $event"
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