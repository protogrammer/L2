package env

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"
)

func envString(key string) string {
	s, ok := os.LookupEnv(key)
	if !ok {
		log.Panicf("[env] Cannot find var `%s`", key)
	}
	return s
}

var (
	port                string
	reactBuildDirectory string
	domain              string
	databaseDirectory   string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Print("[env] Cannot find .env file")
	}

	port = strings.TrimLeft(envString(`PORT`), ":")
	reactBuildDirectory = envString(`REACT_BUILD_DIRECTORY`)
	domain = envString("DOMAIN")
	databaseDirectory = envString("DATABASE_DIRECTORY")
}

func Port() string {
	return port
}

func ReactBuildDirectory() string {
	return reactBuildDirectory
}

func Domain() string {
	return domain
}

func DatabaseDirectory() string {
	return databaseDirectory
}
