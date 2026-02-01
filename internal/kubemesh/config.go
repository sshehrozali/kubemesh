package kubemesh

import "os"

func GetEnv(key string, fallback string) string {
	val, found := os.LookupEnv(key)
	if (!found) {
		return fallback
	}

	return val
}