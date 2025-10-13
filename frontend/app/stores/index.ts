import { defineStore } from 'pinia'
import type { Bucket, BucketObject } from '../types'
import { filesize } from 'filesize'

interface State {
  buckets: Bucket[]
  selectedBucketName: string | null
  bucketObjects: BucketObject[]
  prefix?: string
  viewMode: 'grid' | 'list'
  isLoadingBuckets: boolean
  isLoadingBucketObjects: boolean
  bucketsError: string | null
  bucketObjectsError: string | null
  invalidBucket: boolean
}

export const getBucketObjectName = (object: BucketObject, prefix?: string): string => {
  let name = object?.name ?? ''
  // create a name from the key, removing the prefix if it exists
  // key is like 'folder1/folder2/file.txt'
  // if prefix is 'folder1/', then name should be 'folder2/file.txt'
  // if prefix is 'folder1/folder2/', then name should be 'file.txt'
  // if prefix is undefined, then name should be the same as key
  // if key is 'folder1/folder2/', then name should be 'folder2
  // if key is 'folder1/folder2/file.txt', then name should be 'file.txt'
  // if key is 'file.txt', then name should be 'file.txt'
  if (prefix && object.key.startsWith(prefix)) {
    name = object.key.slice(prefix.length)
  }
  // if name is empty, then use the last part of the key
  if (!name || name === '') {
    const nameSplit = object.key.split('/')
    name = nameSplit.pop() || nameSplit[0] || object.key
  }
  return name.replace(/\/$/, '') // Remove trailing slash if exists
}

const loadBucketObjects = async ({ bucketName, prefix }: { bucketName: string | null; prefix?: string }) => {
  if (!bucketName) {
    return [] // Return empty array if no bucket is selected
  }
  const api = useNuxtApp().$api
  try {
    const { objects } = await api.getBucketObjects({bucketName, prefix})
    const enrichedObjects = objects.map(object => {
      return {
        ...object,
        name: getBucketObjectName(object, prefix),
        formattedSize: object.size ? filesize(object.size) : undefined, // Use filesize library to format size
      }
    })
    return enrichedObjects
  } catch (error) {
    console.error('Error fetching bucket objects:', error)
    return [] // Return empty array on error
  }
}

export const useFileStore = defineStore('file', () => {
  const buckets = ref([] as Bucket[])
  const selectedBucketName = ref(null as string | null)
  const bucketObjects = ref([] as BucketObject[])
  const prefix = ref(undefined as string | undefined)
  const viewMode = ref('list' as 'grid' | 'list')
  const isLoadingBuckets = ref(false)
  const isLoadingBucketObjects = ref(false)
  const bucketsError = ref(null as string | null)
  const bucketObjectsError = ref(null as string | null)
  const invalidBucket = ref(false)

  // state: (): State => {
  //   const route = useRoute();
  //   const routerFolder = route.params.folder as string[] | null
  //   const currentFolder = ref(routerFolder ? routerFolder.join('/') : '')

  //   return {
  //     buckets: [],
  //     selectedBucketName: null,
  //     bucketObjects: [],
  //     prefix: currentFolder.value,
  //     viewMode: 'list', // Default view mode
  //     isLoadingBuckets: false,
  //     isLoadingBucketObjects: false,
  //     bucketsError: null,
  //     bucketObjectsError: null,
  //     invalidBucket: false,
  //   }
  // },

  async function getBuckets() {
    isLoadingBuckets.value = true
    bucketsError.value = null
    const api = useNuxtApp().$api
    try {
      const response = await api.getBuckets()
      setBuckets(response)
      invalidBucket.value = false
    } catch (error: any) {
      bucketsError.value = error?.message || 'Error fetching buckets'
      invalidBucket.value = true
      console.error('Error fetching buckets:', error)
    } finally {
      isLoadingBuckets.value = false
    }
  }

  async function getBucketObjects() {
    isLoadingBucketObjects.value = true
    bucketObjectsError.value = null
    try {
      const objects = await loadBucketObjects({
        bucketName: selectedBucketName.value,
        prefix: prefix.value
      })
      setBucketObjects(objects)
    } catch (error: any) {
      bucketObjectsError.value = error?.message || 'Error fetching bucket objects'
      console.error('Error fetching bucket objects:', error)
    } finally {
      isLoadingBucketObjects.value = false
    }
  }

  async function refreshBucketObjects() {
    if (selectedBucketName.value) {
      isLoadingBucketObjects.value = true
      bucketObjectsError.value = null
      try {
        const objects = await loadBucketObjects({
          bucketName: selectedBucketName.value,
          prefix: prefix.value
        })
        setBucketObjects(objects)
      } catch (error: any) {
        bucketObjectsError.value = error?.message || 'Error refreshing bucket objects'
        console.error('Error refreshing bucket objects:', error)
      } finally {
        isLoadingBucketObjects.value = false
      }
    }
  }

  function setBuckets(newBuckets: Bucket[]) {
    buckets.value = newBuckets
  }

  function setSelectedBucketName(bucketName: string | null) {
    console.log('Setting selected bucket name:', bucketName, selectedBucketName.value)
    if (bucketName !== selectedBucketName.value) {
      prefix.value = undefined // Reset prefix when changing bucket
      selectedBucketName.value = bucketName
      // Check if bucket is valid
      invalidBucket.value = !buckets.value.some(b => b.name === bucketName)
      console.log('Fetching bucket objects for:', bucketName)
      getBucketObjects()
    }
  }

  function setBucketObjects(newBucketObjects: BucketObject[]) {
    bucketObjects.value = newBucketObjects
  }

  function setPrefix(newPrefix: string | undefined) {
    if (newPrefix !== prefix.value) {
      prefix.value = newPrefix
      if (selectedBucketName.value) {
        getBucketObjects()
      }
    }
  }

  function setViewMode(newViewMode: 'grid' | 'list') {
    viewMode.value = newViewMode
  }

  function getBucketObjectsIncludingText(text: string) {
    return bucketObjects.value.filter(object => object?.name?.includes(text))
  }

  return {
    buckets,
    selectedBucketName,
    bucketObjects,
    prefix,
    viewMode,
    isLoadingBuckets,
    isLoadingBucketObjects,
    bucketsError,
    bucketObjectsError,
    invalidBucket,
    getBuckets,
    getBucketObjects,
    refreshBucketObjects,
    setBuckets,
    setSelectedBucketName,
    setBucketObjects,
    setPrefix,
    setViewMode,
    getBucketObjectsIncludingText
  }
});