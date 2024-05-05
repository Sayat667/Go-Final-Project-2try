package utils

import (
	"os"
)

func Setenv() {
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_NAME", "gorm_example")
	os.Setenv("DB_PASSWORD", "admin")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_SSLMODE", "disable")
}
