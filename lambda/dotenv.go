package lambda

import (
    "github.com/joho/godotenv"
    "log"
)

func SetupEnv() {
    const path = "/opt/.env"
    if err := godotenv.Load(path); err != nil {
        log.Fatalf("ERROR: env not found: %s", path)
    }
}
