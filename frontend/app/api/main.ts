import type { BucketsAPI } from './buckets'
import bucketsApi from './buckets'

export type ApiFetch = (
  url: string,
  options?: {
    method?: 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE'
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    body?: any
    query?: Record<string, string | number | boolean | undefined>
  },
// eslint-disable-next-line @typescript-eslint/no-explicit-any
) => Promise<any>

export interface Api extends BucketsAPI {}

export default (apiFetch: ApiFetch): Api => ({
  ...bucketsApi(apiFetch),
})
