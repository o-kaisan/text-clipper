package common

import (
	"io/ioutil"
	"log"
	"os"
)

func Env(name string, defaultValue string) string {
	value := os.Getenv(name)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func SetupGlobalLogger(discard bool) {
	// You can also specify DEBUG mode by environment variable DEBUG. It's handy to debug in runtime.
	if Env("DEBUG", "") == "" && discard {
		log.SetOutput(ioutil.Discard)
		return
	}
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
}
