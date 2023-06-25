package lambda

import (
    "github.com/joho/godotenv"
    "io/fs"
    "log"
    "os"
    "path"
    "strings"
)

func SetupEnv() {
    const dirpath = "/opt"
    err := fs.WalkDir(os.DirFS(dirpath), ".", func(fpath string, d fs.DirEntry, err error) error {
        if err != nil {
            log.Printf("ERROR: walk env file error, path: %s, error: %s", fpath, err)
            return nil
        }
        if d.IsDir() {
            return nil
        }
        if !strings.HasPrefix(fpath, ".env") {
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
        log.Fatalf("ERROR: walk env file dir: %s, error: %s", dirpath, err)
    }
}
