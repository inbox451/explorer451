<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import {
  ChevronRight,
} from 'lucide-vue-next'
import { Breadcrumb, BreadcrumbItem, BreadcrumbLink, BreadcrumbList, BreadcrumbSeparator } from './ui/breadcrumb'

const route = useRoute()
const router = useRouter()
const folderPath = computed(() => route.params.folder as string[] || [])

interface Crumb {
  name: string
  path?: string
}

const rootPath = computed(() => {
  const bucket = route.params.bucket as string
  return router.resolve({ name: 'files', params: { bucket } }).fullPath
})

const crumbs = computed(() => folderPath.value.reduce((acc, crumb, index) => {
  if (crumb === '') return acc // Skip empty crumbs
  const path = index < folderPath.value.length 
    ? router.resolve({
        name: 'files',
        params: {
          bucket: route.params.bucket,
          folder: folderPath.value.slice(0, index + 1).filter(Boolean),
        }
      }).fullPath
    : undefined
  acc.push({
    name: crumb,
    path,
  })
  return acc
}, [{ name: 'Root', path: rootPath.value }] as Crumb[]))

</script>
<template>
  <Breadcrumb>
    <BreadcrumbList>
      <BreadcrumbItem v-for="(crumb, index) in crumbs" :key="index">
        <template v-if="index < crumbs.length - 1">
          <BreadcrumbLink :href="crumb.path">
            {{ crumb.name }}
          </BreadcrumbLink>
          <BreadcrumbSeparator />
        </template>
        <template v-else>
          <span class="text-muted-foreground">{{ crumb.name }}</span>
        </template>
      </BreadcrumbItem>
    </BreadcrumbList>
  </Breadcrumb>
</template>