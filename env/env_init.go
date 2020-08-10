package env

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type DefaultValueSupplier = func() string

func init() {
	godotenv.Load(".env.local")
	godotenv.Load() // The Original .env
}

func FindEnv(envname string) (v string, err error) {
	v = os.Getenv(envname)
	if v == "" {
		return "", errors.New(fmt.Sprintf("env:%s not found", envname))
	}
	return
}

func FindEnvOrDefault(envname string, defaultV string) (v string) {
	v = os.Getenv(envname)
	if v == "" {
		return defaultV
	}
	return
}

func FindEnvOrDefaultSupplier(envName string, supplier DefaultValueSupplier) (v string) {
	v = os.Getenv(envName)
	if v == "" {
		return supplier()
	}
	return
}
