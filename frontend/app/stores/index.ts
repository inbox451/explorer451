import { defineStore } from 'pinia'
import type { Bucket, BucketObject } from '../types'
import { filesize } from 'filesize'

interface State {
  buckets: Bucket[]
  selectedBucketName: string | null
  bucketObjects: BucketObject[]
  prefix?: string
  viewMode: 'grid' | 'list'
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

const getBucketObjects = async ({ bucketName, prefix }: { bucketName: string | null; prefix?: string }) => {
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

export const useBucketStore = defineStore('bucket', {
  state: (): State => {
    const route = useRoute();
    const routerFolder = route.params.folder as string[] | null
    const currentFolder = ref(routerFolder ? routerFolder.join('/') : '')

    return {
      buckets: [],
      selectedBucketName: null,
      bucketObjects: [],
      prefix: currentFolder.value,
      viewMode: 'list', // Default view mode
    }
  },

  actions: {
    async getBuckets() {
      const api = useNuxtApp().$api
      try {
        const response = await api.getBuckets()
        this.setBuckets(response)
      } catch (error) {
        console.error('Error fetching buckets:', error)
      }
    },

    async getBucketObjects() {
      const { selectedBucketName, prefix } = this
      const objects = await getBucketObjects({ 
        bucketName: selectedBucketName, 
        prefix
      })
      this.setBucketObjects(objects)
    },

    async refreshBucketObjects() {
      const { selectedBucketName, prefix } = this
      if (selectedBucketName) {
        const objects = await getBucketObjects({ 
          bucketName: selectedBucketName, 
          prefix
        })
        this.setBucketObjects(objects)
      }
    },

    setBuckets(buckets: Bucket[]) {
      this.buckets = buckets
    },

    setSelectedBucketName(bucketName: string | null) {
      console.log('Setting selected bucket name:', bucketName, this.selectedBucketName)
      if (bucketName !== this.selectedBucketName) {
        this.prefix = undefined // Reset prefix when changing bucket
        this.selectedBucketName = bucketName
        console.log('Fetching bucket objects for:', bucketName)
        this.getBucketObjects()
      }
    },

    setBucketObjects(bucketObjects: BucketObject[]) {
      this.bucketObjects = bucketObjects
    },

    setPrefix(prefix: string | undefined) {
      if (prefix !== this.prefix) {
        this.prefix = prefix
        if (this.selectedBucketName) {
          this.getBucketObjects()
        }
      }
    },

    setViewMode(viewMode: 'grid' | 'list') {
      this.viewMode = viewMode
    },
  },

  getters: {
    bucketByName: (state: State) => (name: string) => {
      return state.buckets.find(bucket => bucket.name === name) || null
    },
    getBucketObjectsIncludingText: (state: State) =>
      (text: string) => state.bucketObjects.filter(object => object?.name?.includes(text)),
  
  }    
});