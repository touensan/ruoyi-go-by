import request from '@/utils/request'

export const getAccountOverview = () => request({ url: '/account/v1/overview', method: 'get' })
export const listAccountLedger = () => request({ url: '/account/v1/ledger', method: 'get' })
export const listAccountRecharges = () => request({ url: '/account/v1/recharges', method: 'get' })
export const createAccountRecharge = (points: number, payType: string) =>
  request({
    url: '/account/v1/recharges', method: 'post', data: { points, payType },
    headers: { 'Idempotency-Key': window.crypto.randomUUID(), repeatSubmit: false }
  })
export const getExchangeCenter = () => request({ url: '/account/v1/exchange', method: 'get' })
export const submitRewardClaim = (taskPublicId: string, providerSubject: string) =>
  request({ url: '/account/v1/exchange/claims', method: 'post', data: { taskPublicId, providerSubject } })
export const redeemExchangeCode = (code: string) =>
  request({ url: '/account/v1/exchange/redeem', method: 'post', data: { code } })

export const listPointAccounts = (params: object) => request({ url: '/admin/v1/platform/accounts', method: 'get', params })
export const adjustPoints = (data: object) => request({ url: '/admin/v1/platform/accounts/adjust', method: 'post', data })
export const getPlatformSettings = () => request({ url: '/admin/v1/platform/exchange/settings', method: 'get' })
export const savePlatformSettings = (data: object) => request({ url: '/admin/v1/platform/exchange/settings', method: 'put', data })
export const testRainyun = () => request({ url: '/admin/v1/platform/exchange/rainyun/test', method: 'post' })
export const listTasks = () => request({ url: '/admin/v1/platform/exchange/tasks', method: 'get' })
export const createTask = (data: object) => request({ url: '/admin/v1/platform/exchange/tasks', method: 'post', data })
export const updateTask = (publicId: string, data: object) => request({ url: `/admin/v1/platform/exchange/tasks/${publicId}`, method: 'put', data })
export const listClaims = (params: object) => request({ url: '/admin/v1/platform/exchange/claims', method: 'get', params })
export const reviewClaim = (publicId: string, data: object) => request({ url: `/admin/v1/platform/exchange/claims/${publicId}/review`, method: 'post', data })
export const listCodes = (params: object) => request({ url: '/admin/v1/platform/exchange/codes', method: 'get', params })
export const generateCodes = (data: object) => request({ url: '/admin/v1/platform/exchange/codes', method: 'post', data })
