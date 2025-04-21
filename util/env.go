package util

import "os"

func IsDebug() bool {
	return os.Getenv("DEBUG") != ""
}

func GetEnvDefault(key, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		val = defaultValue
	}
	return val
}

// === auth ===

func IsNoLoginMode() bool {
	return GetEnvDefault("NO_LOGIN_MODE", "") != ""
}

//  === db ===

func GetPostgresHost() string {
	return GetEnvDefault("POSTGRES_HOST", "localhost")
}

func GetPostgresPort() string {
	return GetEnvDefault("POSTGRES_PORT", "5432")
}

func GetPostgresUser() string {
	return GetEnvDefault("POSTGRES_USER", "postgres")
}

func GetPostgresPassword() string {
	return GetEnvDefault("POSTGRES_PASSWORD", "postgres")
}

func GetPostgresDBName() string {
	return GetEnvDefault("POSTGRES_DBNAME", "postgres")
}

func IsEnableSJTUFeature() bool {
	return GetEnvDefault("ENABLE_SJTU_FEATURE", "") != ""
}

func GetHashSalt() string {
	return GetEnvDefault("HASH_SALT", "")
}
