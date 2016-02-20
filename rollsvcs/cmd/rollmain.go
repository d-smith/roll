package main

import (
	"flag"
	"fmt"
	"github.com/xtraclabs/roll/dbutil"
	"github.com/xtraclabs/roll/roll"
	"github.com/xtraclabs/roll/rollsvcs"
	"log"
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
	log.Println(unsecureBanner)
	return rollsvcs.DefaultUnsecureConfig()
}

func createDynamoDBConfig() *roll.CoreConfig {
	return rollsvcs.DefaultConfig()
}

func createUnsecureMariaDBConfig() *roll.CoreConfig {
	log.Println(unsecureBanner)
	return rollsvcs.MariaDBConfig()
}

func creatMariaDBConfig() *roll.CoreConfig {
	return nil
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
		log.Println(unsecureBanner)
		if dbutil.UseMariaDB() {
			log.Println("Using maria db")
			coreConfig = rollsvcs.MariaDBConfig()
		} else {
			log.Println("Using dynamo db")
			coreConfig = rollsvcs.DefaultUnsecureConfig()
		}
	} else {
		coreConfig = rollsvcs.DefaultConfig()
	}

	rollsvcs.RunRoll(*port, coreConfig)
}
