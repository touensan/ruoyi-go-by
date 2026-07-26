import request from '@/utils/request'
import type {
  AuthActivation, AuthAnomaly, AuthApplication, AuthAuditEvent, AuthDeactivation,
  AuthFeature, AuthInstallation, AuthLease, AuthLicense, AuthPageResult,
  AuthPlan, AuthResult, AuthRuntimeQuery, AuthSigningKey, AuthVersion,
  IssuedAuthLicense
} from '@/types/api/auth'

const base = '/admin/v1/auth'
function newIdempotencyKey() {
  if (typeof crypto.randomUUID === 'function') return crypto.randomUUID()
  const bytes = crypto.getRandomValues(new Uint8Array(16))
  bytes[6] = (bytes[6] & 0x0f) | 0x40
  bytes[8] = (bytes[8] & 0x3f) | 0x80
  const value = Array.from(bytes, byte => byte.toString(16).padStart(2, '0')).join('')
  return `${value.slice(0, 8)}-${value.slice(8, 12)}-${value.slice(12, 16)}-${value.slice(16, 20)}-${value.slice(20)}`
}
const noPersistHeaders = () => ({ repeatSubmit: false, 'Idempotency-Key': newIdempotencyKey() })

export const listAuthApplications = (params: object) =>
  request({ url: `${base}/applications`, method: 'get', params }) as Promise<AuthPageResult<AuthApplication>>
export const getAuthApplication = (id: string) =>
  request({ url: `${base}/applications/${id}`, method: 'get' }) as Promise<AuthResult<AuthApplication>>
export const addAuthApplication = (data: object) =>
  request({ url: `${base}/applications`, method: 'post', data })
export const updateAuthApplication = (id: string, data: object) =>
  request({ url: `${base}/applications/${id}`, method: 'put', data })
export const listAuthVersions = (appId: string) =>
  request({ url: `${base}/applications/${appId}/versions`, method: 'get' }) as Promise<AuthResult<AuthVersion[]>>
export const addAuthVersion = (appId: string, data: object) =>
  request({ url: `${base}/applications/${appId}/versions`, method: 'post', data })
export const updateAuthVersion = (appId: string, id: string, data: object) =>
  request({ url: `${base}/applications/${appId}/versions/${id}`, method: 'put', data })
export const listAuthFeatures = (appId: string) =>
  request({ url: `${base}/applications/${appId}/features`, method: 'get' }) as Promise<AuthResult<AuthFeature[]>>
export const addAuthFeature = (appId: string, data: object) =>
  request({ url: `${base}/applications/${appId}/features`, method: 'post', data })
export const updateAuthFeature = (appId: string, id: string, data: object) =>
  request({ url: `${base}/applications/${appId}/features/${id}`, method: 'put', data })

export const listAuthPlans = (params: object) =>
  request({ url: `${base}/plans`, method: 'get', params }) as Promise<AuthPageResult<AuthPlan>>
export const getAuthPlan = (id: string) =>
  request({ url: `${base}/plans/${id}`, method: 'get' }) as Promise<AuthResult<AuthPlan>>
export const addAuthPlan = (data: object) =>
  request({ url: `${base}/plans`, method: 'post', data })
export const updateAuthPlan = (id: string, data: object) =>
  request({ url: `${base}/plans/${id}`, method: 'put', data })

export const listAuthLicenses = (params: object) =>
  request({ url: `${base}/licenses`, method: 'get', params }) as Promise<AuthPageResult<AuthLicense>>
export const getAuthLicense = (id: string) =>
  request({ url: `${base}/licenses/${id}`, method: 'get' }) as Promise<AuthResult<AuthLicense>>
export const issueAuthLicense = (data: object) =>
  request({ url: `${base}/licenses`, method: 'post', data, headers: noPersistHeaders() }) as Promise<AuthResult<IssuedAuthLicense>>
export const actOnAuthLicense = (id: string, action: string, reason = '') =>
  request({ url: `${base}/licenses/${id}/${action}`, method: 'post', data: { reason }, headers: noPersistHeaders() }) as Promise<AuthResult<AuthLicense | IssuedAuthLicense>>

export const listAuthInstallations = (params: AuthRuntimeQuery) =>
  request({ url: `${base}/installations`, method: 'get', params }) as Promise<AuthPageResult<AuthInstallation>>
export const listAuthActivations = (params: AuthRuntimeQuery) =>
  request({ url: `${base}/activations`, method: 'get', params }) as Promise<AuthPageResult<AuthActivation>>
export const listAuthLeases = (params: AuthRuntimeQuery) =>
  request({ url: `${base}/leases`, method: 'get', params }) as Promise<AuthPageResult<AuthLease>>
export const listAuthDeactivations = (params: AuthRuntimeQuery) =>
  request({ url: `${base}/deactivations`, method: 'get', params }) as Promise<AuthPageResult<AuthDeactivation>>
export const listAuthAnomalies = (params: AuthRuntimeQuery) =>
  request({ url: `${base}/anomalies`, method: 'get', params }) as Promise<AuthPageResult<AuthAnomaly>>
export const deactivateAuthActivation = (id: string, reason: string) =>
  request({ url: `${base}/activations/${id}/deactivate`, method: 'post', data: { reason }, headers: noPersistHeaders() }) as Promise<AuthResult<AuthDeactivation>>

export const listAuthSigningKeys = (params: object) =>
  request({ url: `${base}/signing-keys`, method: 'get', params }) as Promise<AuthPageResult<AuthSigningKey>>
export const registerAuthSigningKey = (data: object) =>
  request({ url: `${base}/signing-keys`, method: 'post', data, headers: { repeatSubmit: false } })
export const actOnAuthSigningKey = (kid: string, action: string, reason = '') =>
  request({ url: `${base}/signing-keys/${encodeURIComponent(kid)}/${action}`, method: 'post', data: action === 'revoke' ? { reason } : {}, headers: { repeatSubmit: false } })

export const listAuthAudit = (params: object) =>
  request({ url: `${base}/audit`, method: 'get', params }) as Promise<AuthPageResult<AuthAuditEvent>>
