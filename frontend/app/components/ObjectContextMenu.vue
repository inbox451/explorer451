<script setup lang="ts">
import {
  ContextMenu,
  ContextMenuContent,
  ContextMenuItem,
  ContextMenuSeparator,
  ContextMenuTrigger,
} from '@/components/ui/context-menu'
import { MoreHorizontal, Download, File, Trash2 } from 'lucide-vue-next'
import type { BucketObject } from '~/types';

const props = defineProps<{
  item: BucketObject,
  selectedBucketName: string
}>()
</script>

<template>
  <ContextMenu>
    <ContextMenuTrigger><slot/></ContextMenuTrigger>
    <ContextMenuContent>
      <template v-if="item.type !== 'folder'">
        <ContextMenuItem v-if="item.type !== 'folder'">
          <FileDownload :item="item" :selectedBucketName="selectedBucketName">
            <Download class="h-4 w-4 mr-2" />
            Download
          </FileDownload>
        </ContextMenuItem>
        <ContextMenuItem v-if="item.type !== 'folder'">
          <File class="h-4 w-4 mr-2" />
          Copy URL
        </ContextMenuItem>
      </template>
      <template v-else>
        <ContextMenuItem>
          <FileLink :item="item" :selectedBucketName="selectedBucketName" class="flex items-center">
            <FileIcon :item="item" class="h-4 w-4 mr-2" />
            Open Folder
          </FileLink>
        </ContextMenuItem>
      </template>
      <ContextMenuSeparator />
      <ContextMenuItem class="text-destructive" v-if="item.key">
        <Trash2 class="h-4 w-4 mr-2" />
        Delete
      </ContextMenuItem>
    </ContextMenuContent>
  </ContextMenu>  
</template>