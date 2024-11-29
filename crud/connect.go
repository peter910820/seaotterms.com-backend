package crud

import (
	"fmt"
	"os"
)

var dsn = fmt.Sprintf("user=%s password=%s dbname=%s port=%s sslmode=disable",
	os.Getenv("DB_OWNER"),
	os.Getenv("DB_PASSWORD"),
	os.Getenv("DB_NAME"),
	os.Getenv("DB_PORT"))
