package main

import (
	"flag"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/xtraclabs/roll/dbutil"
	"github.com/xtraclabs/roll/roll"
	"github.com/xtraclabs/roll/rollsvcs"
	"os"
	"strings"
)

var unsecureBanner = `
_   _                                     ___  ___          _
| | | |                                    |  \/  |         | |
| | | |_ __  ___  ___  ___ _   _ _ __ ___  | .  . | ___   __| | ___
| | | | '_ \/ __|/ _ \/ __| | | | '__/ _ \ | |\/| |/ _ \ / _  |/ _ \
| |_| | | | \__ \  __/ (__| |_| | | |  __/ | |  | | (_) | (_| |  __/
 \___/|_| |_|___/\___|\___|\__,_|_|  \___| \_|  |_/\___/ \__,_|\___|
`

func createUnsecureDynamoDBConfig() *roll.CoreConfig {
	log.Info(unsecureBanner)
	return rollsvcs.DefaultUnsecureConfig()
}

func createDynamoDBConfig() *roll.CoreConfig {
	return rollsvcs.DefaultConfig()
}

func createUnsecureMariaDBConfig() *roll.CoreConfig {
	log.Info(unsecureBanner)
	return rollsvcs.MariaDBUnsecureConfig()
}

func createMariaDBConfig() *roll.CoreConfig {
	return rollsvcs.MariaDBSecureConfig()
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	setLoggingLevel()
}

func setLoggingLevel() {

	logLevel := strings.ToLower(os.Getenv("ROLL_LOGGING_LEVEL"))
	switch logLevel {
	default:
		log.SetLevel(log.InfoLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
		//Note - makes no sense to set the default log levels to fatal or to panic
	}

	log.Info("log level set: ", log.GetLevel())
}

func main() {

	var port = flag.Int("port", -1, "Port to listen on")
	var unsecureMode = flag.Bool("unsecure", false, "Boot in unsecure mode")
	flag.Parse()
	if *port == -1 {
		fmt.Println("Must specify a -port argument")
		return
	}

	var coreConfig *roll.CoreConfig

	if *unsecureMode == true {
		fmt.Println(unsecureBanner)
		if dbutil.UseMariaDB() {
			log.Info("Using maria db")
			coreConfig = rollsvcs.MariaDBUnsecureConfig()
		} else {
			log.Info("Using dynamo db")
			coreConfig = rollsvcs.DefaultUnsecureConfig()
		}
	} else {
		if dbutil.UseMariaDB() {
			log.Info("Using maria db")
			coreConfig = rollsvcs.MariaDBSecureConfig()
		} else {
			log.Info("Using dynamo db")
			coreConfig = rollsvcs.DefaultConfig()
		}
	}

	rollsvcs.RunRoll(*port, coreConfig)
}
