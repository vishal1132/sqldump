package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
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
	ls := exec.Command("ls")
	grepPattern2 := "\\.sql"
	grep := exec.Command("grep", grepPattern2)
	grep.Stdin, err = ls.StdoutPipe()
	if err != nil {
		log.Fatalln(err)
	}
	var b bytes.Buffer
	grep.Stdout = &b
	err = grep.Start()
	err = ls.Run()
	err = grep.Wait()
	if err != nil {
		log.Println("unable to grep .sql files in the directory ", err)
	}
	s := strings.Split(b.String(), "\n")
	files := s[:len(s)-1]
	maketar("sql", files)
	remove("sql", files)
}

func remove(grep string, files []string) error {
	for _, file := range files {
		os.Remove(file)
	}
	return nil
}

func maketar(grep string, files []string) error {
	t := time.Now()
	formatted := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	formattedWithName := formatted + ".sql.gz"
	file, err := os.Create(formattedWithName)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not create tarball file '%s', got error '%s'", formattedWithName, err.Error()))
	}
	defer file.Close()

	gzipWriter := gzip.NewWriter(file)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()
	for _, filePath := range files {
		err := addFileToTarWriter(filePath, tarWriter)
		if err != nil {
			return errors.New(fmt.Sprintf("Could not add file '%s', to tarball, got error '%s'", filePath, err.Error()))
		}
	}
	return nil
}

func addFileToTarWriter(filePath string, tarWriter *tar.Writer) error {
	file, err := os.Open(filePath)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not open file '%s', got error '%s'", filePath, err.Error()))
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return errors.New(fmt.Sprintf("Could not get stat for file '%s', got error '%s'", filePath, err.Error()))
	}

	header := &tar.Header{
		Name:    filePath,
		Size:    stat.Size(),
		Mode:    int64(stat.Mode()),
		ModTime: stat.ModTime(),
	}

	err = tarWriter.WriteHeader(header)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not write header for file '%s', got error '%s'", filePath, err.Error()))
	}

	_, err = io.Copy(tarWriter, file)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not copy the file '%s' data to the tarball, got error '%s'", filePath, err.Error()))
	}

	return nil
}
