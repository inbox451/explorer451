<template>
  <div class="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 xl:grid-cols-8 gap-4">
    <Card v-for="item in items" :key="item.key" :class="['p-4 hover:shadow-md transition-shadow', item.type === 'folder' ? 'cursor-pointer hover:bg-muted/50' : 'cursor-pointer']">
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
</template>
<script setup lang="ts">
import { Card, CardContent } from '@/components/ui/card'
import FileLink from '@/components/FileLink.vue'
import FileIcon from '@/components/FileIcon.vue'
const props = defineProps<{ items: any[], selectedBucketName: string }>()
</script>
