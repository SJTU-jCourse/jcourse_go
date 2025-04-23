package password_hash

type Algorithm string

const (
	AlgorithmPBKDF2SHA256 Algorithm = "pbkdf2_sha256"
)

type Hasher interface {
	HashPassword(password string) string
	GetAlgorithm() Algorithm
	GetIteration() int
	GetSalt() string
}

func NewHasher(algorithm Algorithm, iteration int, salt string) Hasher {
	switch algorithm {
	case AlgorithmPBKDF2SHA256:
		return NewPBK2DFSHA256Hasher(salt, iteration)
	}
	return nil
}
