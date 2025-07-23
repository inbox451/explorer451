<template>
  <Card>
    <CardContent class="flex flex-col items-center justify-center py-12">
      <Folder class="h-12 w-12 text-muted-foreground mb-4" />
      <h3 class="text-lg font-semibold mb-2">No items found</h3>
      <p class="text-muted-foreground text-center">
        {{ searchQuery ? 'No files or folders match your search criteria.' : 'This folder is empty. Upload some files or create folders to get started.' }}
      </p>
    </CardContent>
  </Card>
</template>
<script setup lang="ts">
import { Card, CardContent } from '@/components/ui/card'
import { Folder } from 'lucide-vue-next'
const props = defineProps<{ searchQuery: string }>()
</script>
