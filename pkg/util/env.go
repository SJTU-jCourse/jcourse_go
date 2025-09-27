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

// === Session ===

func GetSessionSecret() string {
	return GetEnvDefault("SESSION_SECRET", "")
}

// === CRSF ===

func GetCSRFSecret() string {
	return GetEnvDefault("CSRF_KEY", "")
}

// === Time ===
func GetTimeLocationStr() string {
	return GetEnvDefault("TIME_ZONE", "Asia/Shanghai")
}
