// Copyright © 2019-2020 catenocrypt.  See LICENSE file for license information.

package workcache

import (
	"fmt"
	"log"
	"math"
	"strconv"

	"github.com/spf13/viper"
)

var configRead bool = false
var configFileName string = "config"

func SetConfigFile(configFile string) {
	configFileName = configFile
}

func readConfigIfNeeded() {
	if configRead {
		return
	} // already read

	// set defaults
	// no default for "Main.NodeRpc", must be set
	viper.SetDefault("Main.ListenIpPort", ":7176")
	viper.SetDefault("Main.CachePeristFileName", "")
	viper.SetDefault("Main.RestMaxActiveRequests", 200)
	viper.SetDefault("Main.BackgroundWorkerCount", 4)
	viper.SetDefault("Main.MaxOutRequests", 8)
	viper.SetDefault("Main.EnablePregeneration", 1)
	viper.SetDefault("Main.PregenerationQueueSize", 10000)
	viper.SetDefault("Main.MaxCacheAgeDays", 30)

	// read config file
	viper.SetConfigName(configFileName) // name of config file (without extension)
	viper.AddConfigPath(".")            // optionally look for config in the working directory
	viper.AddConfigPath("/")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	configRead = true
	log.Println("Config file has been read")
}

func ConfigGetString(keyName string) string {
	readConfigIfNeeded()
	return viper.GetString(keyName)
}

func ConfigGetIntWithDefault(keyName string, defaultVal int) int {
	str := ConfigGetString(keyName)
	val, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		log.Println("Invalid int config value", str)
		return defaultVal
	}
	return int(val)
}

func EnablePregeneration() int {
	return ConfigGetIntWithDefault("Main.EnablePregeneration", 1)
}

func PregenerationQueueSize() int {
	val := ConfigGetIntWithDefault("Main.PregenerationQueueSize", 10000)
	val = int(math.Max(float64(val), float64(0)))
	val = int(math.Min(float64(val), float64(100000)))
	return val
}
