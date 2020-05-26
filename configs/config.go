package configs

import (
	"os"
)

type Config struct {
	Key   string
	Value string
}

func retrieveEnvOrDefault(key string, defaultValue string) string {
	result := os.Getenv(key)
	if len(result) == 0 {
		result = defaultValue
	}
	return result
}

//redis基本配置
var RedisHost string = retrieveEnvOrDefault("REDIS_HOST", "127.0.0.1")
var RedisPort string = retrieveEnvOrDefault("REDIS_PORT", "6379")
var RedisChannel string = retrieveEnvOrDefault("REDIS_CHANNEL", "0")
var EncryptKey string = retrieveEnvOrDefault("ENCRYPT_KEY", "xxxxxxxxxxxxxxxx")
var BackendGoHost string = retrieveEnvOrDefault("BACKENND_GO_HOST", "127.0.0.1")
var BackendGoPort string = retrieveEnvOrDefault("BACKENND_GO_PORT", "8080")
