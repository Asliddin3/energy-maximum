package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

type Config struct {
	Host string
	Port int

	DB_HOST                string
	DB_PORT                int
	DB_NAME                string
	DB_PASS                string
	DB_USER                string
	StaticFilePath         string
	AccessTokenPublicKey   string
	AccessTokenPrivateKey  string
	RefreshTokenPublicKey  string
	RefreshTokenPrivateKey string
	AccessTokenExpiresIn   time.Duration
	RefreshTokenExpiresIn  time.Duration
	AccessTokenMaxAge      int
	RefreshTokenMaxAge     int
}

func Load() Config {
	c := Config{}
	c.Host = cast.ToString(getOrReturnDefault("HOST", "localhost"))
	c.Port = cast.ToInt(getOrReturnDefault("PORT", 8000))
	c.DB_HOST = cast.ToString(getOrReturnDefault("DB_HOST", "localhost"))
	c.DB_NAME = cast.ToString(getOrReturnDefault("DB_NAME", "market"))
	c.DB_PASS = cast.ToString(getOrReturnDefault("DB_PASS", "admin"))
	c.DB_USER = cast.ToString(getOrReturnDefault("DB_USER", "postgres"))
	c.DB_PORT = cast.ToInt(getOrReturnDefault("DB_PORT", 5432))
	c.StaticFilePath = cast.ToString(getOrReturnDefault("STATIC_FILE_PATH", "/Users/devop/go/src/myProjects/EnergyMaximum/media/"))
	c.AccessTokenPrivateKey = cast.ToString(getOrReturnDefault("ACCESS_TOKEN_PRIVATE_KEY", "LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlCUEFJQkFBSkJBTzVIKytVM0xrWC91SlRvRHhWN01CUURXSTdGU0l0VXNjbGFFKzlaUUg5Q2VpOGIxcUVmCnJxR0hSVDVWUis4c3UxVWtCUVpZTER3MnN3RTVWbjg5c0ZVQ0F3RUFBUUpCQUw4ZjRBMUlDSWEvQ2ZmdWR3TGMKNzRCdCtwOXg0TEZaZXMwdHdtV3Vha3hub3NaV0w4eVpSTUJpRmI4a25VL0hwb3piTnNxMmN1ZU9wKzVWdGRXNApiTlVDSVFENm9JdWxqcHdrZTFGY1VPaldnaXRQSjNnbFBma3NHVFBhdFYwYnJJVVI5d0loQVBOanJ1enB4ckhsCkUxRmJxeGtUNFZ5bWhCOU1HazU0Wk1jWnVjSmZOcjBUQWlFQWhML3UxOVZPdlVBWVd6Wjc3Y3JxMTdWSFBTcXoKUlhsZjd2TnJpdEg1ZGdjQ0lRRHR5QmFPdUxuNDlIOFIvZ2ZEZ1V1cjg3YWl5UHZ1YStxeEpXMzQrb0tFNXdJZwpQbG1KYXZsbW9jUG4rTkVRdGhLcTZuZFVYRGpXTTlTbktQQTVlUDZSUEs0PQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQ=="))
	c.AccessTokenPublicKey = cast.ToString(getOrReturnDefault("ACCESS_TOKEN_PUBLIC_KEY", "LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUZ3d0RRWUpLb1pJaHZjTkFRRUJCUUFEU3dBd1NBSkJBTzVIKytVM0xrWC91SlRvRHhWN01CUURXSTdGU0l0VQpzY2xhRSs5WlFIOUNlaThiMXFFZnJxR0hSVDVWUis4c3UxVWtCUVpZTER3MnN3RTVWbjg5c0ZVQ0F3RUFBUT09Ci0tLS0tRU5EIFBVQkxJQyBLRVktLS0tLQ=="))
	c.RefreshTokenPrivateKey = cast.ToString(getOrReturnDefault("REFRESH_TOKEN_PRIVATE_KEY", "LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlCT1FJQkFBSkJBSWFJcXZXeldCSndnYjR1SEhFQ01RdHFZMTI5b2F5RzVZMGlGcG51a0J1VHpRZVlQWkE4Cmx4OC9lTUh3Rys1MlJGR3VxMmE2N084d2s3TDR5dnY5dVY4Q0F3RUFBUUpBRUZ6aEJqOUk3LzAxR285N01CZUgKSlk5TUJLUEMzVHdQQVdwcSswL3p3UmE2ZkZtbXQ5NXNrN21qT3czRzNEZ3M5T2RTeWdsbTlVdndNWXh6SXFERAplUUloQVA5UStrMTBQbGxNd2ZJbDZtdjdTMFRYOGJDUlRaZVI1ZFZZb3FTeW40YmpBaUVBaHVUa2JtZ1NobFlZCnRyclNWZjN0QWZJcWNVUjZ3aDdMOXR5MVlvalZVRlVDSUhzOENlVHkwOWxrbkVTV0dvV09ZUEZVemhyc3Q2Z08KU3dKa2F2VFdKdndEQWlBdWhnVU8yeEFBaXZNdEdwUHVtb3hDam8zNjBMNXg4d012bWdGcEFYNW9uUUlnQzEvSwpNWG1heWtsaFRDeWtXRnpHMHBMWVdkNGRGdTI5M1M2ZUxJUlNIS009Ci0tLS0tRU5EIFJTQSBQUklWQVRFIEtFWS0tLS0t"))
	c.RefreshTokenPublicKey = cast.ToString(getOrReturnDefault("REFRESH_TOKEN_PUBLIC_KEY", "LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUZ3d0RRWUpLb1pJaHZjTkFRRUJCUUFEU3dBd1NBSkJBSWFJcXZXeldCSndnYjR1SEhFQ01RdHFZMTI5b2F5Rwo1WTBpRnBudWtCdVR6UWVZUFpBOGx4OC9lTUh3Rys1MlJGR3VxMmE2N084d2s3TDR5dnY5dVY4Q0F3RUFBUT09Ci0tLS0tRU5EIFBVQkxJQyBLRVktLS0tLQ=="))
	c.AccessTokenExpiresIn = cast.ToDuration(getOrReturnDefault("ACCESS_TOKEN_EXPIRED_IN", time.Duration(time.Minute*60)))
	c.RefreshTokenExpiresIn = cast.ToDuration(getOrReturnDefault("REFRESH_TOKEN_EXPIRED_IN", time.Duration(time.Minute*300)))
	c.AccessTokenMaxAge = cast.ToInt(getOrReturnDefault("ACCESS_TOKEN_MAXAGE", 60))
	c.RefreshTokenMaxAge = cast.ToInt(getOrReturnDefault("REFRESH_TOKEN_MAXAGE", 300))

	return c
}
func getOrReturnDefault(key string, defaultValue interface{}) interface{} {
	err := godotenv.Load(".env")
	if err != nil {
		// log.Fatalf("Error loading .env file")
		return defaultValue
	}
	_, exists := os.LookupEnv(key)
	if exists {
		return os.Getenv(key)
	}
	return defaultValue
}
