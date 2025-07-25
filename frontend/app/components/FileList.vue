<template>
  <Table class="w-full">
    <TableHeader>
      <TableRow>
        <TableHead class="min-w-0 flex-1">Name</TableHead>
        <TableHead class="text-right w-[100px]">Size</TableHead>
        <TableHead class="w-[100px]">Modified</TableHead>
        <TableHead class="w-[100px]">Actions</TableHead>
      </TableRow>
    </TableHeader>
    <TableBody>
      <TableRow v-for="item in items" :key="item.key">
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
                  <FileDownload :item="item" :selectedBucketName="selectedBucketName">
                    Download
                  </FileDownload>
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
</template>
<script setup lang="ts">
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import FileLink from '@/components/FileLink.vue'
import FileIcon from '@/components/FileIcon.vue'
import { NuxtTime } from '#components'
import { DropdownMenu, DropdownMenuTrigger, DropdownMenuContent, DropdownMenuItem, DropdownMenuSeparator } from '@/components/ui/dropdown-menu'
import { MoreHorizontal, Download, File, Trash2 } from 'lucide-vue-next'
const props = defineProps<{ items: any[], selectedBucketName: string }>()
</script>
