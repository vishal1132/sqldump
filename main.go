package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/joho/godotenv"
)

//to get environment values
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

type options struct {
	ExecutionStartDate time.Time
}

func main() {
	godotenv.Load()
	mysqldumppath := getEnv("MysqlDumpPath", "/ivr/software/thirdparty/system/lampp/lampp/bin/mysqldump")
	user := "-u" + getEnv("Username", "root")
	password := "-p" + getEnv("Password", "r00t")
	dbs := getEnv("DBNames", "subs_engine0")
	t := time.Now()
	formatted := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	cmdDump := &exec.Cmd{
		Path:   mysqldumppath,
		Stdout: os.Stdout,
		Stdin:  os.Stdin,
		Stderr: os.Stderr,
		Args:   []string{mysqldumppath, user, password, "--databases", dbs, "--result-file", formatted + ".sql"},
	}
	log.Println(cmdDump.String())
	if err := cmdDump.Run(); err != nil {
		log.Println("error occurred ", err)
	}
}
