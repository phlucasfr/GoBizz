package util

import "os"

func IsTesting() bool {
	return os.Getenv("TEST_ENV") == "true"
}
