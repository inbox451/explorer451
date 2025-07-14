import type { ApiFetch } from './main'
import type { Bucket, ObjectsAPIResponse } from '~/types'

export interface BucketsAPI {
  getBucketDetail: (id: string) => Promise<Bucket>
  getBuckets: () => Promise<Bucket[]>
  getBucketObjects: ({ bucketName, prefix }: { bucketName: string; prefix?: string }) => Promise<ObjectsAPIResponse>
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
  }
})
