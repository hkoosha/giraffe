package gtesting

import (
	"os"
	"testing"

	. "github.com/hkoosha/giraffe/internal/dot0"
)

const (
	envDisallowTest       = "GIRAFFE_DISALLOW_TEST"
	envSkipOnceValidation = "GIRAFFE_SKIP_SETUP_VALIDATION"
)

func IsTest() bool {
	_, defined := os.LookupEnv(envDisallowTest)

	return !defined
}

func IsBypassOnceValidation() bool {
	return IsTest() && os.Getenv(envSkipOnceValidation) == "true"
}

func EnsureTesting() {
	if !IsTest() {
		panic(EF("test code called outside tests"))
	}
}

func SetNoTesting() {
	if err := os.Setenv(envDisallowTest, "true"); err != nil {
		panic(EF("failed to set env var %s: %v", envDisallowTest, err))
	}
}

func SetSkipOnceValidation(t *testing.T) {
	t.Helper()
	EnsureTesting()

	t.Setenv(envSkipOnceValidation, "true")
}
