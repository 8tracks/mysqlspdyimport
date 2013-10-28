package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"

	"database/sql"
	_ "github.com/8tracks.com/mysqlspdyimport/vendor/go-sql-driver/mysql"
)

const NotInsertStatement = "Not an insert statement"

var numGoFuncs int
var maxErrors int
var wg = sync.WaitGroup{}

type MysqlConf struct {
	username string
	password string
	host     string
	port     string
	database string
}

var mysqlConf = MysqlConf{}

func (m *MysqlConf) ConnectionString() string {
	var str string

	if m.username != "" {
		str += m.username + ":"
	}

	if m.password != "" {
		str += m.password
	}

	if m.username != "" {
		str += "@"
	}

	str += "tcp("
	str += m.host
	str += ":" + m.port
	str += ")/"
	str += m.database

	return str
}

func execQueries(db *sql.DB, query <-chan string) {
	for q := range query {
		_, err := db.Exec(q)
		if err != nil {
			log.Println("Error:", err)
			log.Println("=> Query skipped:", q)
		}
	}

	wg.Done()
}

func startUp() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	if mysqlConf.database == "" {
		flag.Usage()
		fmt.Println("Please specify a database name to connect to.")
		os.Exit(1)
	}
}

func cleanUpQuery(q string) (string, error) {
	q = strings.Trim(q, "\t\n ")

	insertPart := q[0:6]
	if len(insertPart) != 6 || strings.ToLower(insertPart) != "insert" {
		return "", errors.New(NotInsertStatement)
	}

	return q, nil
}

func init() {
	flag.StringVar(&mysqlConf.username, "u", "", "Database username")
	flag.StringVar(&mysqlConf.password, "p", "", "Database password")
	flag.StringVar(&mysqlConf.host, "h", "localhost", "Database host")
	flag.StringVar(&mysqlConf.port, "P", "3306", "Database port")
	flag.StringVar(&mysqlConf.database, "d", "", "Database name")
	flag.IntVar(&numGoFuncs, "c", 3, "Number of goroutines to execute queries.")
}

func main() {
	startUp()
	log.Print("Start 'er up!")

	input := bufio.NewReader(os.Stdin)
	pipe := make(chan string)

	db, err := sql.Open("mysql", mysqlConf.ConnectionString())
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	// Without this db connections are recycled pretty fast and will cause
	// errors relating to too many file descriptors open.
	//
	// With go1.2 we should also call SetMaxOpenConns to set total number of connections
	// to leave open.
	db.SetMaxIdleConns(numGoFuncs)

	// Ensure we can connect to this sucker
	err = db.Ping()
	if err != nil {
		log.Fatalln("Could not connect to database!")
	}

	log.Printf("Spinning up %d processes", numGoFuncs)
	for i := 0; i < numGoFuncs; i++ {
		wg.Add(1)
		go execQueries(db, pipe)
	}

	for {
		switch line, err := input.ReadString('\n'); err {
		case nil:
			q, err := cleanUpQuery(line)
			if err != nil {
				log.Printf("Could not process line: %s", line)
				continue
			}
			pipe <- q

		case io.EOF:
			close(pipe)
			wg.Wait()
			log.Print("")
			log.Print("All done!")
			os.Exit(0)

		default:
			log.Fatalln("No input!", err)
		}
	}

}
