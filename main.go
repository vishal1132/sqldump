package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

var (
	mysqldumppath string = getEnv("MysqlDumpPath", "/ivr/software/thirdparty/system/lampp/lampp/bin/mysqldump")
	user          string = "-u" + getEnv("Username", "root")
	password      string = "-p" + getEnv("Password", "r00t")
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

func dumpSpecificTables(i int) {
	dbenvKey := "DB" + strconv.Itoa(i+1)
	t := time.Now()
	formatted := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	db := getEnv(dbenvKey, "cdrlog")
	tableEnvKey := "Tables" + strconv.Itoa(i+1)
	tables := getEnv(tableEnvKey, "")
	formattedWithName := formatted + "_" + db
	cmdDump := &exec.Cmd{
		Path:   mysqldumppath,
		Stdout: os.Stdout,
		Stdin:  os.Stdin,
		Stderr: os.Stderr,
		Args:   []string{mysqldumppath, user, password, db, tables, "--result-file", formattedWithName + ".sql"},
	}
	log.Println(cmdDump.String())
	if err := cmdDump.Run(); err != nil {
		log.Println("error occurred ", err)
	}
}

func makeCompleteDBBackup() {
	dbs := getEnv("DBNames", "")
	if dbs == "" {
		return
	}
	t := time.Now()
	formatted := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	formattedWithName := formatted + "_completedbs"
	cmdDump := &exec.Cmd{
		Path:   mysqldumppath,
		Stdout: os.Stdout,
		Stdin:  os.Stdin,
		Stderr: os.Stderr,
		Args:   []string{mysqldumppath, user, password, "--databases", dbs, "--result-file", formattedWithName + ".sql"},
	}
	log.Println(cmdDump.String())
	if err := cmdDump.Run(); err != nil {
		log.Println("error occurred ", err)
	}
}

func main() {
	godotenv.Load()
	makeCompleteDBBackup()
	numDBs, err := strconv.Atoi(getEnv("NumDestination", "1"))
	if err != nil {
		log.Fatal("unknown digit in number destination ", err)
	}
	for i := 0; i < numDBs; i++ {
		dumpSpecificTables(i)
	}
}
