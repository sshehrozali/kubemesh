package kubemesh

import (
	"log"
	"regexp"
)

func IsValidPort(port string) bool {
	match, err := regexp.Match("^[0-9]+$", []byte(port))
	if err != nil {
		log.Fatal("Error compiling regex for port validation")
	}

	return match
}

func IsValidNodeNic(nic string) bool {
	return true
}