<template>
  <main class="flex-1 p-6">
    <h2 class="text-center">Pick a bucket</h2>
    <p class="text-muted-foreground text-center">
      Please select a bucket to view its contents.
    </p>
    <div class="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 xl:grid-cols-8 gap-4">
      <Title>Available Buckets</Title>
      <NuxtLink :to="`/files/${bucket.name}`" v-for="bucket in fileStore.buckets" :key="bucket.name">
        <Card class="p-4 hover:shadow-md transition-shadow cursor-pointer">
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
</template>
<script setup lang="ts">
import { Card, CardContent } from '@/components/ui/card'
import { Database } from 'lucide-vue-next'
import { Title } from '#components'
import { useFileStore } from '@/stores'
const fileStore = useFileStore()
</script>
