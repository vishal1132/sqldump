package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
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
	formattedWithName := formatted + "_" + db + ".sql"
	fmt.Println("tables ", tables, " formattedName ", formattedWithName)
	args := []string{user, password}
	arr := strings.Split(tables, ",")
	args = append(args, db)
	args = append(args, arr...)
	args = append(args, "--result-file", formattedWithName)
	cmdDump := exec.Command(mysqldumppath, args...)
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
	log.Println(dbs)
	t := time.Now()
	formatted := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	formattedWithName := formatted + "_databases" + ".sql"
	args := []string{user, password, "--databases"}
	arr := strings.Split(dbs, ",")
	args = append(args, arr...)
	args = append(args, "--result-file", formattedWithName)
	cmdDump := exec.Command(mysqldumppath, args...)
	if err := cmdDump.Run(); err != nil {
		log.Println("error occurred ", err)
	}
}

func main() {
	godotenv.Load()

	numDBs, err := strconv.Atoi(getEnv("NumDBs", "1"))
	if err != nil {
		log.Fatal("unknown digit in number destination ", err)
	}
	for i := 0; i < numDBs; i++ {
		dumpSpecificTables(i)
	}
	makeCompleteDBBackup()
}
