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
    <!-- Main Content -->
    <main v-if="selectedBucketName" class="flex-1 p-6">
      <!-- Toolbar -->
      <div class="flex items-center justify-between mb-6">
        <div class="flex items-center gap-2">
          <div class="relative">
            <Search class="absolute left-2 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
            <Input
              placeholder="Search files and folders..."
              v-model="searchQuery"
              class="pl-8 w-80"
            />
          </div>
        </div>
        <div class="flex items-center gap-2">
          <Dialog v-model="isCreateFolderDialogOpen">
            <DialogTrigger>
              <Button variant="outline" size="sm">
                <Plus class="h-4 w-4 mr-2" />
                New Folder
              </Button>
            </DialogTrigger>
            <DialogContent>
              <DialogHeader>
                <DialogTitle>Create New Folder</DialogTitle>
                <DialogDescription>Enter a name for the new folder.</DialogDescription>
              </DialogHeader>
              <div class="grid gap-4 py-4">
                <div class="grid grid-cols-4 items-center gap-4">
                  <Label for="folder-name" class="text-right">Name</Label>
                  <Input id="folder-name" placeholder="Folder name" class="col-span-3" />
                </div>
              </div>
              <DialogFooter>
                <Button type="submit" @click="isCreateFolderDialogOpen = false">Create Folder</Button>
              </DialogFooter>
            </DialogContent>
          </Dialog>
          <Dialog v-model="isUploadDialogOpen">
            <template #trigger>
              <Button size="sm">
                <Upload class="h-4 w-4 mr-2" />
                Upload
              </Button>
            </template>
            <DialogContent>
              <DialogHeader>
                <DialogTitle>Upload Files</DialogTitle>
                <DialogDescription>Select files to upload to the current folder.</DialogDescription>
              </DialogHeader>
              <div class="grid gap-4 py-4">
                <div class="grid grid-cols-4 items-center gap-4">
                  <Label for="file-upload" class="text-right">Files</Label>
                  <Input id="file-upload" type="file" multiple class="col-span-3" />
                </div>
              </div>
              <DialogFooter>
                <Button type="submit" @click="isUploadDialogOpen = false">Upload Files</Button>
              </DialogFooter>
            </DialogContent>
          </Dialog>
          <div class="flex items-center border rounded-md">
            <Button
              :variant="viewMode === 'list' ? 'default' : 'ghost'"
              size="sm"
              @click="viewMode = 'list'"
              class="rounded-r-none"
            >
              <List class="h-4 w-4" />
            </Button>
            <Button
              :variant="viewMode === 'grid' ? 'default' : 'ghost'"
              size="sm"
              @click="viewMode = 'grid'"
              class="rounded-l-none"
            >
              <Grid3X3 class="h-4 w-4" />
            </Button>
          </div>
        </div>
      </div>
      <!-- Folder Crumbs -->
      <div class="mb-4">
        <FolderCrumbs />
      </div>
      <!-- File List -->
      <Table class="w-full" v-if="viewMode === 'list'">
        <TableHeader>
          <TableRow>
            <TableHead class="min-w-0 flex-1">Name</TableHead>
            <TableHead class="text-right w-[100px]">Size</TableHead>
            <TableHead class="w-[100px]">Modified</TableHead>
            <TableHead class="w-[100px]">Actions</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow
            v-for="item in itemsToRender"
            :key="item.key"
          >
            <TableCell class="gap-2 min-w-0">
              <FileLink :item="item" :selectedBucketName="selectedBucketName" :Title="item.name" class="relative w-full flex items-center">
                <div class="flex items-center w-full t-0 l-0 b-0 r-0 truncate absolute">
                  <FileIcon :item="item" class="inline" />
                  <span class="ml-2 inline-block overflow-hidden text-ellipsis">{{ item.name }}</span>
                </div>
              </FileLink>
            </TableCell>
            <TableCell class="text-right">
              {{ item.type === 'folder' ? '' : item.formattedSize }}
            </TableCell>
            <TableCell>
              <NuxtTime v-if="item.type !== 'folder' && item.lastModified" 
                :datetime="item.lastModified" 
                dateStyle="short"
                timeStyle="medium" />
            </TableCell>
            <TableCell>
              <DropdownMenu>
                <DropdownMenuTrigger>
                  <MoreHorizontal class="h-4 w-4 cursor-pointer" />
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end">
                  <template v-if="item.type !== 'folder'">
                    <DropdownMenuItem>
                      <Download class="h-4 w-4 mr-2" />
                      Download
                    </DropdownMenuItem>
                    <DropdownMenuItem>
                      <File class="h-4 w-4 mr-2" />
                      Copy URL
                    </DropdownMenuItem>
                    <DropdownMenuSeparator />
                  </template>
                  <DropdownMenuItem class="text-destructive">
                    <Trash2 class="h-4 w-4 mr-2" />
                    Delete
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>
      <div v-else class="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 xl:grid-cols-8 gap-4">
        <Card
          v-for="item in itemsToRender"
          :key="item.key"
          :class="['p-4 hover:shadow-md transition-shadow', item.type === 'folder' ? 'cursor-pointer hover:bg-muted/50' : 'cursor-pointer']"
        >
          <CardContent class="p-0">
            <div class="flex flex-col items-center gap-2">
              <FileLink :item="item" :selectedBucketName="selectedBucketName">
                <div class="p-3 rounded-lg bg-muted">
                  <FileIcon :item="item" size="h-8 w-8" />
                </div>
                <div class="text-center">
                  <p class="text-sm font-medium truncate w-full overflow-hidden text-ellipsis" :title="item.name">{{ item.name }}</p>
                  <p class="text-xs text-muted-foreground" v-if="item.type !== 'folder'">{{ item.formattedSize }}</p>
                </div>
              </FileLink>
            </div>
          </CardContent>
        </Card>
      </div>
      <Card v-if="filteredItems.length === 0">
        <CardContent class="flex flex-col items-center justify-center py-12">
          <Folder class="h-12 w-12 text-muted-foreground mb-4" />
          <h3 class="text-lg font-semibold mb-2">No items found</h3>
          <p class="text-muted-foreground text-center">
            {{ searchQuery ? 'No files or folders match your search criteria.' : 'This folder is empty. Upload some files or create folders to get started.' }}
          </p>
        </CardContent>
      </Card>
    </main>
    <main v-else class="flex-1 p-6">
      <h2 class="text-center">Pick a bucket</h2>
      <p class="text-muted-foreground text-center">
        Please select a bucket to view its contents.
      </p>
      <div class="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 xl:grid-cols-8 gap-4">
        <Title>Available Buckets</Title>
        <NuxtLink :to="`/files/${bucket.name}`" v-for="bucket in bucketStore.buckets" :key="bucket.name">
          <Card
            class="p-4 hover:shadow-md transition-shadow cursor-pointer"
          >
            <CardContent class="flex flex-col items-center gap-2">
              <div class="p-3 rounded-lg bg-muted">
                <Database class="h-8 w-8 text-muted-foreground" />
              </div>
              <div class="text-center">
                <p class="text-sm font-medium truncate w-full overflow-hidden text-ellipsis" :title="bucket.name">{{ bucket.name }}</p>
                <p class="text-xs text-muted-foreground">Bucket</p>
              </div>
            </CardContent>
          </Card>
        </NuxtLink>
      </div>
    </main>
  </div>
</template>

<style scoped>
/* Add your styles here */
</style>