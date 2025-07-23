<script setup lang="ts">
import { ref, watch } from 'vue'
import { LogOut, RefreshCw, Home } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue
} from '@/components/ui/select'
import { useBucketStore } from '@/stores'

const nuxtApp = useNuxtApp()
const route = useRoute()
const router = nuxtApp.$router

const bucketStore = useBucketStore()
const { getBuckets, buckets, setSelectedBucketName, setPrefix, refreshBucketObjects } = bucketStore

const currentBucketName = ref(route.params.bucket as string | null)
const currentFolderPath = ref('')

// Set initial selected bucket
setSelectedBucketName(currentBucketName.value)

// Watch for changes in route.params.folder and update prefix in the store
watch(
  () => route.params.folder,
  (newFolder) => {
    const newPrefix = newFolder ? (Array.isArray(newFolder) ? newFolder.join('/') : newFolder) : ''
    setPrefix(newPrefix)
    currentFolderPath.value = newPrefix
  },
  { immediate: true }
)

// Watch for changes in route.params.bucket and update selected bucket
watch(
  () => route.params.bucket,
  (newBucket) => {
    setSelectedBucketName(newBucket as string | null)
    currentBucketName.value = newBucket as string | null
  }
)

const handleBucketSelect = (bucketName: string) => {
  router.push({ path: `/files/${bucketName}` })
}

const handleRefresh = () => {
  if(currentBucketName.value) {
    refreshBucketObjects()
  } else {
    getBuckets()
  }
}

const rootPath = computed(() => {
  return router.resolve({ name: 'files' }).fullPath
})
</script>
<template>
  <!-- Header -->
  <header class="border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
    <Title>Bucket {{ currentBucketName }} files at {{ currentFolderPath }}</Title>
    <div class="flex h-16 items-center justify-between px-6">
      <div class="flex items-center gap-4">
        <h1 class="text-xl font-semibold">explorer451 <span class="text-2xl">🔥</span></h1>
        <Select v-model="currentBucketName" @update:modelValue="handleBucketSelect">
          <SelectTrigger class="w-48">
            <SelectValue placeholder="Select a bucket" />
          </SelectTrigger>
          <SelectContent>
            <SelectGroup>
              <SelectLabel>Buckets</SelectLabel>
              <SelectItem
                v-for="bucket in buckets"
                :key="bucket.name"
                :value="bucket.name"
              >
                {{ bucket.name }}
              </SelectItem>
            </SelectGroup>
          </SelectContent>
        </Select>
      </div>
      <div class="flex items-center gap-2">
        <Button variant="outline" size="sm" @click="nuxtApp.$router.push(rootPath)">
          <Home class="h-4 w-4 mr-2" />
          Home
        </Button>
        <Button variant="outline" size="sm" @click="handleRefresh">
          <RefreshCw class="h-4 w-4 mr-2" />
          Refresh
        </Button>
        <Button variant="outline" size="sm">
          <LogOut class="h-4 w-4 mr-2" />
          Logout
        </Button>
      </div>
    </div>
  </header>
</template>