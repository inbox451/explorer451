<template>
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
            <FolderPlus class="h-4 w-4 mr-2" />
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
            <Button type="submit" @click="$emit('closeCreateFolderDialog')">Create Folder</Button>
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
            <Button type="submit" @click="$emit('closeUploadDialog')">Upload Files</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
      <div class="flex items-center border rounded-md">
        <Button
          :variant="viewMode === 'list' ? 'default' : 'ghost'"
          size="sm"
          @click="$emit('update:viewMode', 'list')"
          class="rounded-r-none"
        >
          <List class="h-4 w-4" />
        </Button>
        <Button
          :variant="viewMode === 'grid' ? 'default' : 'ghost'"
          size="sm"
          @click="$emit('update:viewMode', 'grid')"
          class="rounded-l-none"
        >
          <LayoutGrid class="h-4 w-4" />
        </Button>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { ref } from 'vue'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Dialog, DialogTrigger, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Label } from '@/components/ui/label'
import { Search, FolderPlus, Upload, List, LayoutGrid } from 'lucide-vue-next'
const props = defineProps<{ searchQuery: string, viewMode: string, isUploadDialogOpen: boolean, isCreateFolderDialogOpen: boolean }>()
const emit = defineEmits(['update:viewMode'])

const searchQuery = ref(props.searchQuery)
const isUploadDialogOpen = ref(props.isUploadDialogOpen)
const isCreateFolderDialogOpen = ref(props.isCreateFolderDialogOpen)
</script>
