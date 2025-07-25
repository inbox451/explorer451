import type { ApiFetch } from './main'
import type { Bucket, ObjectsAPIResponse } from '~/types'

interface BucketObjectFileUrl {
  url: string
}
export interface BucketsAPI {
  getBucketDetail: (id: string) => Promise<Bucket>
  getBuckets: () => Promise<Bucket[]>
  getBucketObjects: ({ bucketName, prefix }: { bucketName: string; prefix?: string }) => Promise<ObjectsAPIResponse>
  getBucketObjectFileUrl: ({ bucketName, key }: { bucketName: string; key: string }) => Promise<BucketObjectFileUrl>
}

export default (apiFetch: ApiFetch): BucketsAPI => ({
  getBucketDetail(id: string) {
    return apiFetch(`/buckets/${id}`, {
      method: 'GET'
    })
  },

  getBuckets() {
    return apiFetch(`/buckets`, {
      method: 'GET'
    })
  },

  getBucketObjects({ bucketName, prefix }) {
    return apiFetch(`/buckets/${bucketName}/objects`, {
      method: 'GET',
      query: { prefix }
    })
  },

  getBucketObjectFileUrl({ bucketName, key }) {
    return apiFetch(`/buckets/${bucketName}/objects/${key}`, {
      method: 'GET',
      query: { presigned: true }
    })
  }
})
