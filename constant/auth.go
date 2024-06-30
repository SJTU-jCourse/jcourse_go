package constant

const (
	AuthVerifyCodeLen = 6
	VerifyCodeDigits  = "0123456789"
)

type UserRole = string

const (
	UserRoleNormal = "normal"
	UserRoleAdmin  = "admin"
)

type HashAlgorithmType = string

const (
	HashAlgorithmPBK2DF HashAlgorithmType = "pbkdf2_sha256"
)

const PasswordHashIteration = 720000
