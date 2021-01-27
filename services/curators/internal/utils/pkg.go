package utils

import "log"

// LogFatalIfErr calls log.Fatal(err) if err is not nil.
func LogFatalIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
