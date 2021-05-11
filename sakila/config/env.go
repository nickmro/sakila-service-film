package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Env represents the application environment.
type Env struct {
	logger         string
	mySQLHost      string
	mySQLName      string
	mySQLPassword  string
	mySQLPort      string
	mySQLUser      string
	port           string
	redisHost      string
	redisPassword  string
	redisPort      int
	redisKeyPrefix string
}

const (
	configPath = "."
	configType = "env"
)

const (
	envKeyLogger         = "LOGGER"
	envKeyMySQLHost      = "MYSQL_HOST"
	envKeyMySQLName      = "MYSQL_NAME"
	envKeyMySQLPassword  = "MYSQL_PASSWORD"
	envKeyMySQLPort      = "MYSQL_PORT"
	envKeyMySQLUser      = "MYSQL_USER"
	envKeyPort           = "PORT"
	envKeyRedisHost      = "REDIS_HOST"
	envKeyRedisPassword  = "REDIS_PASSWORD"
	envKeyRedisPort      = "REDIS_PORT"
	envKeyRedisKeyPrefix = "REDIS_KEY_PREFIX"
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

	redisHost := v.GetString(envKeyRedisHost)
	if redisHost == "" {
		return nil, missingEnvError(envKeyRedisHost)
	}

	redisPort := v.GetInt(envKeyRedisPort)
	if redisPort == 0 {
		return nil, missingEnvError(envKeyRedisPort)
	}

	redisPassword := v.GetString(envKeyRedisPassword)

	redisKeyPrefix := v.GetString(envKeyRedisKeyPrefix)

	port := v.GetString(envKeyPort)
	if port == "" {
		port = defaultValuePort
	}

	env := &Env{
		logger:         logger,
		mySQLHost:      mySQLHost,
		mySQLName:      mySQLName,
		mySQLPassword:  mySQLPassword,
		mySQLPort:      mySQLPort,
		mySQLUser:      mySQLUser,
		port:           port,
		redisHost:      redisHost,
		redisPassword:  redisPassword,
		redisPort:      redisPort,
		redisKeyPrefix: redisKeyPrefix,
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

// GetRedisHost returns the Redis host.
func (e *Env) GetRedisHost() string {
	return e.redisHost
}

// GetRedisPort returns the Redis port.
func (e *Env) GetRedisPort() int {
	return e.redisPort
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

// GetRedisKeyPrefix returns the redis cache key prefix.
func (e *Env) GetRedisKeyPrefix() string {
	return e.redisKeyPrefix
}

func missingEnvError(key string) error {
	return fmt.Errorf("%w: %s", ErrorMissing, key)
}
