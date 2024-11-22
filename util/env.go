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

// === redis ===

func GetRedisHost() string {
	return GetEnvDefault("REDIS_HOST", "localhost")
}

func GetRedisPort() string {
	return GetEnvDefault("REDIS_PORT", "6379")
}
func GetRedisPassword() string {
	return GetEnvDefault("REDIS_PASSWORD", "")
}

// === SMTP ===

func GetSMTPHost() string {
	return GetEnvDefault("SMTP_HOST", "localhost")
}

func GetSMTPPort() string {
	return GetEnvDefault("SMTP_PORT", "465")
}

func GetSMTPUser() string {
	return GetEnvDefault("SMTP_USERNAME", "user")
}

func GetSMTPPassword() string {
	return GetEnvDefault("SMTP_PASSWORD", "pass")
}

func GetSMTPSender() string {
	return GetEnvDefault("SMTP_SENDER", "user")
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
