package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Env represents the application environment.
type Env struct {
	logger              string
	mySQLHost           string
	mySQLName           string
	mySQLPassword       string
	mySQLPort           string
	mySQLUser           string
	port                string
	redisPassword       string
	redisURL            string
	redisCacheKeyPrefix string
}

const (
	configPath = "."
	configType = "env"
)

const (
	envKeyLogger              = "LOGGER"
	envKeyMySQLHost           = "MYSQL_HOST"
	envKeyMySQLName           = "MYSQL_NAME"
	envKeyMySQLPassword       = "MYSQL_PASSWORD"
	envKeyMySQLPort           = "MYSQL_PORT"
	envKeyMySQLUser           = "MYSQL_USER"
	envKeyPort                = "PORT"
	envKeyRedisURL            = "REDIS_URL"
	envKeyRedisPassword       = "REDIS_PASSWORD"
	envKeyRedisCacheKeyPrefix = "REDIS_CACHE_KEY_PREFIX"
)

const (
	defaultLoggerEnvironment = "DEVELOPMENT"
	defaultValuePort         = "3000"
)

// GetEnv returns the application environment.
func GetEnv(fileName string) (*Env, error) {
	v := viper.New()
	v.AutomaticEnv()
	v.AddConfigPath(configPath)
	v.SetConfigType(configType)
	v.SetConfigName(fileName)

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	logger := v.GetString(envKeyLogger)
	if logger == "" {
		logger = defaultLoggerEnvironment
	}

	mySQLHost := v.GetString(envKeyMySQLHost)
	if mySQLHost == "" {
		return nil, missingEnvError(envKeyMySQLHost)
	}

	mySQLPort := v.GetString(envKeyMySQLPort)
	if mySQLPort == "" {
		return nil, missingEnvError(envKeyMySQLPort)
	}

	mySQLName := v.GetString(envKeyMySQLName)
	if mySQLName == "" {
		return nil, missingEnvError(envKeyMySQLName)
	}

	mySQLPassword := v.GetString(envKeyMySQLPassword)

	mySQLUser := v.GetString(envKeyMySQLUser)

	redisURL := v.GetString(envKeyRedisURL)
	if redisURL == "" {
		return nil, missingEnvError(envKeyRedisURL)
	}

	redisPassword := v.GetString(envKeyRedisPassword)

	redisCacheKeyPrefix := v.GetString(envKeyRedisCacheKeyPrefix)

	port := v.GetString(envKeyPort)
	if port == "" {
		port = defaultValuePort
	}

	env := &Env{
		logger:              logger,
		mySQLHost:           mySQLHost,
		mySQLName:           mySQLName,
		mySQLPassword:       mySQLPassword,
		mySQLPort:           mySQLPort,
		mySQLUser:           mySQLUser,
		port:                port,
		redisPassword:       redisPassword,
		redisURL:            redisURL,
		redisCacheKeyPrefix: redisCacheKeyPrefix,
	}

	return env, nil
}

// GetMySQLURL returns the MySQL URL.
func (e *Env) GetMySQLURL() string {
	if e.mySQLUser == "" || e.mySQLPassword == "" {
		return fmt.Sprintf("tcp(%s:%s)/%s?parseTime=true",
			e.mySQLHost,
			e.mySQLPort,
			e.mySQLName)
	}

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		e.mySQLUser,
		e.mySQLPassword,
		e.mySQLHost,
		e.mySQLPort,
		e.mySQLName)
}

// GetRedisURL returns the Redis URL.
func (e *Env) GetRedisURL() string {
	return e.redisURL
}

// GetRedisPassword returs the Redis password.
func (e *Env) GetRedisPassword() string {
	return e.redisPassword
}

// GetLogger returns the logger environment.
func (e *Env) GetLogger() string {
	return e.logger
}

// GetPort return the port.
func (e *Env) GetPort() string {
	return e.port
}

// GetRedisCacheKeyPrefix returns the redis cache key prefix.
func (e *Env) GetRedisCacheKeyPrefix() string {
	return e.redisCacheKeyPrefix
}

func missingEnvError(key string) error {
	return fmt.Errorf("%w: %s", ErrorMissing, key)
}
