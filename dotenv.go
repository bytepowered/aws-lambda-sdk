package lambda

import (
	"github.com/joho/godotenv"
	"io/fs"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
)

func SetupEnv() {
	const dirpath = "/opt"
	err := fs.WalkDir(os.DirFS(dirpath), ".", func(fpath string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Printf("WARN: on env file, path: %s, error: %s", fpath, err)
			return nil
		}
		if d.IsDir() {
			return nil
		}
		if !(strings.HasPrefix(fpath, ".env") || strings.HasSuffix(fpath, ".env")) {
			return nil
		}
		// path: .env, .env-nft
		fullpath := path.Join(dirpath, fpath)
		if err := godotenv.Load(fullpath); err != nil {
			log.Printf("ERROR: load env file: %s, error: %s", fullpath, err)
			return err
		}
		log.Printf("INFO: load env file: %s", fullpath)
		return nil
	})
	if err != nil {
		log.Fatalf("ERROR: walk env dir: %s, error: %s", dirpath, err)
	}
}

func RequiredEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic("Load env NOT-FOUND, key: " + key)
	}
	return v
}

func RequiredStrEnv(key string, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func RequiredIntEnv(key string, def int) int {
	str := RequiredEnv(key)
	if str == "" {
		return def
	}
	if v, err := strconv.Atoi(str); err != nil {
		return def
	} else {
		return v
	}
}
