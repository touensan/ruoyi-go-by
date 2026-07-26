export interface AuthPageQuery {
  pageNum: number
  pageSize: number
}

export interface AuthPageResult<T> {
  code: number
  msg: string
  rows: T[]
  total: number
}

export interface AuthResult<T> {
  code: number
  msg: string
  data: T
}

export interface AuthApplication {
  publicId: string
  appCode: string
  name: string
  status: string
  lockVersion: number
  createdAt: string
  updatedAt: string
}

export interface AuthVersion {
  publicId: string
  version: string
  channel: string
  status: string
  minimumProtocolVersion: number
  releasedAt?: string
  createdAt: string
  updatedAt: string
}

export interface AuthFeature {
  publicId: string
  featureCode: string
  name: string
  valueType: string
  status: string
  createdAt: string
  updatedAt: string
}

export interface AuthEntitlement {
  featureCode: string
  valueType?: string
  value: unknown
}

export interface AuthPlan {
  publicId: string
  applicationPublicId: string
  applicationCode: string
  planCode: string
  name: string
  licenseKind: string
  validityDays?: number
  maxActivations: number
  maxConcurrentSeats: number
  allowedPlatforms: string[]
  minimumClientVersion?: string
  refreshIntervalSeconds: number
  leaseOfflineSeconds: number
  failureGraceSeconds: number
  offlineLicenseSeconds?: number
  status: string
  lockVersion: number
  entitlements?: AuthEntitlement[]
  createdAt: string
  updatedAt: string
}

export interface AuthLicense {
  publicId: string
  customerPublicId: string
  planPublicId: string
  applicationCode: string
  status: string
  keyMask: string
  startsAt?: string
  expiresAt?: string
  revocationSerial: number
  replacementPublicId?: string
  lockVersion: number
  createdAt: string
  updatedAt: string
}

export interface IssuedAuthLicense {
  license: AuthLicense
  licenseKey: string
  replay: boolean
}

export interface AuthRuntimeQuery extends AuthPageQuery {
  applicationPublicId?: string
  licensePublicId?: string
  status?: string
  search?: string
}

export interface AuthInstallation {
  publicId: string
  customerPublicId: string
  applicationPublicId: string
  applicationCode: string
  bindingType: string
  bindingDisplay: string
  platform: string
  clientVersion: string
  status: string
  firstSeenAt: string
  lastSeenAt: string
}

export interface AuthActivation {
  publicId: string
  licensePublicId: string
  licenseKeyMask: string
  applicationPublicId: string
  applicationCode: string
  installationPublicId: string
  bindingDisplay: string
  status: string
  activatedAt: string
  deactivatedAt?: string
}

export interface AuthLease {
  publicId: string
  licensePublicId: string
  licenseKeyMask: string
  activationPublicId: string
  signingKeyKid: string
  serial: number
  issuedAt: string
  notBefore: string
  refreshAfter: string
  offlineUntil: string
  status: string
}

export interface AuthDeactivation {
  publicId: string
  activationPublicId: string
  licensePublicId: string
  licenseKeyMask: string
  installationPublicId: string
  reason?: string
  actorType: string
  actorRef?: string
  deactivatedAt: string
  replay: boolean
}

export interface AuthAnomaly {
  anomalyCode: string
  severity: 'HIGH' | 'MEDIUM'
  subjectType: string
  subjectPublicId: string
  summary: string
  detectedAt: string
}

export interface AuthSigningKey {
  kid: string
  purpose: string
  algorithm: string
  publicKey: string
  status: string
  notBefore: string
  signUntil?: string
  verifyUntil: string
  createdAt: string
  updatedAt: string
}

export interface AuthAuditEvent {
  publicId: string
  occurredAt: string
  actorType: string
  actorRef?: string
  action: string
  targetType: string
  targetPublicId?: string
  requestId?: string
  result: string
  metadata?: Record<string, unknown>
}
