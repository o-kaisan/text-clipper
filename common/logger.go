package common

import (
	"io/ioutil"
	"log"
)

func SetupGlobalLogger(discard bool) {
	// You can also specify DEBUG mode by environment variable DEBUG. It's handy to debug in runtime.
	if Env("DEBUG", "") == "" && discard {
		log.SetOutput(ioutil.Discard)
		return
	}
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
}
