package password

const storeFormat = "%s$%d$%s$%s"

type Algorithm string

const (
	AlgorithmPBKDF2SHA256 Algorithm = "pbkdf2_sha256"
)

type Hasher interface {
	HashPassword(password string) (string, error)
	GetAlgorithm() Algorithm
}

type Validator interface {
	ValidatePassword(password string, hash string) bool
}
