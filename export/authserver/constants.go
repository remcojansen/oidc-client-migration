package authserver

const AccessTokenFormatJwt = "jwt"
const AccessTokenFormatOpaque = "opaque"

const AccessTokenLifetimeLong = 1800
const AccessTokenLifetimeShort = 300
const AccessTokenLifetimeDefault = 0 // avoid inclusion if not set

const OfflineSessionMaxLifetimeDefault = 0 // avoid inclusion if not set
const OfflineSessionIdleTimeoutDefault = 0 // avoid inclusion if not set
const SessionMaxLifetimeDefault = 0        // avoid inclusion if not set
const SessionIdleTimeoutDefault = 0        // avoid inclusion if not set

const ClientTypePublic = "public"
const ClientTypeConfidential = "confidential"

const AuthMethodClientSecretBasic = "client_secret_basic"
const AuthMethodClientSecretJwt = "client_secret_jwt"
const AuthMethodPrivateKeyJwt = "private_key_jwt"
const AuthMethodNone = "none"

const SubjectTypePublic = "public"
const SubjectTypePairwise = "pairwise"

const ApplicationTypeWeb = "web"
const ApplicationTypeNative = "native"

const GrantTypeAuthorizationCode = "authorization_code"
const GrantTypeImplicit = "implicit"
const GrantTypeClientCredentials = "client_credentials"
const GrantTypeResourceOwnerPassword = "password"
const GrantTypeRefreshToken = "refresh_token"
const GrantTypeDeviceCode = "device_code"
const GrantTypeTokenExchange = "urn:ietf:params:oauth:grant-type:token-exchange"
const GrantTypeAccessTokenValidation = "introspect"
