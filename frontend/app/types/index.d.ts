export interface Bucket {
  name: string
  creationDate: string
  region: string
}

export interface BucketObject {
  key: string
  name?: string
  size: number
  formattedSize?: string
  isFolder: boolean
  type: 'file' | 'folder'
  contentType?: string
  lastModified: string
  storageClass: string
  etag: string
}

export interface ApiResponse<T> {
  data: T[]
  pagination: {
    total: number
    limit: number
    offset: number
  }
}

export interface ObjectsAPIResponse {
  objects: BucketObject[]
  isTruncated: boolean
  itemsInPage: number
  pageSize: number
}