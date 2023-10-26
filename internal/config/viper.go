package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

func initViper(configPath string) {
	fmt.Printf("[QiYin Backend] Load config from %s ... ", configPath)
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		log.Panicf("fatal error config file: %s", err.Error())
	}

	fmt.Println("Load ok")
}

// GetInterface get interface config.
func GetInterface(key string) interface{} {
	return viper.Get(key)
}

// GetString get string config.
func GetString(key string) string {
	return viper.GetString(key)
}

// GetInt get int config.
func GetInt(key string) int {
	return viper.GetInt(key)
}

// GetBool get bool config.
func GetBool(key string) bool {
	return viper.GetBool(key)
}

// GetStringOrDefault get string config or use default.
func GetStringOrDefault(key string, defaultValue string) string {
	if viper.IsSet(key) {
		return viper.GetString(key)
	}

	return defaultValue
}

// GetIntOrDefault get int config or use default.
func GetIntOrDefault(key string, defaultValue int) int {
	if viper.IsSet(key) {
		return viper.GetInt(key)
	}

	return defaultValue
}

// GetBoolOrDefault get bool config or use default.
func GetBoolOrDefault(key string, defaultValue bool) bool {
	if viper.IsSet(key) {
		return viper.GetBool(key)
	}

	return defaultValue
}
