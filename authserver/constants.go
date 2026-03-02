package authserver

const TokenFormatJwt = "jwt"
const TokenFormatOpaque = "opaque"

const TokenLifetimeLong = 1800
const TokenLifetimeShort = 300
const TokenLifetimeDefault = 0

const ClientTypePublic = "public"
const ClientTypeConfidential = "confidential"

const AuthMethodClientSecretBasic = "client_secret_basic"
const AuthMethodClientSecretJwt = "client_secret_jwt"
const AuthMethodPrivateKeyJwt = "private_key_jwt"
const AuthMethodNone = "none"
