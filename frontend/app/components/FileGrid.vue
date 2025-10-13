<template>
  <div class="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 xl:grid-cols-8 gap-4">
    <Card v-for="item in items" :key="item.key" :class="['p-4 hover:shadow-md transition-shadow', item.type === 'folder' ? 'cursor-pointer hover:bg-muted/50' : 'cursor-pointer']">
      <CardContent class="p-0">
        <div class="flex flex-col items-center gap-2">
          <ObjectContextMenu :item="item" :selectedBucketName="selectedBucketName">
            <FileLink :item="item" :selectedBucketName="selectedBucketName">
              <div class="p-3 rounded-lg flex items-center justify-center">
                <FileIcon :item="item" size="h-16 w-16" />
              </div>
              <div class="text-center">
                <div class="relative w-full flex items-center h-6">
                  <p class="text-sm font-medium truncate w-full t-0 l-0 b-0 r-0 truncate absolute" :title="item.name">{{ item.name }}</p>
                </div>
                <p class="text-xs text-muted-foreground" v-if="item.type !== 'folder'">{{ item.formattedSize }}</p>
              </div>
            </FileLink>
          </ObjectContextMenu>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
<script setup lang="ts">
import { Card, CardContent } from '@/components/ui/card'
import FileLink from '@/components/FileLink.vue'
import FileIcon from '@/components/FileIcon.vue'
const props = defineProps<{ items: any[], selectedBucketName: string }>()
</script>
