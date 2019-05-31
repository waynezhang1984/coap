package mysqlclient

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)
}

func CreateMysqlClient() (*sql.DB, error) {

	db, err := sql.Open("mysql", "root:123456@tcp(192.168.101.189:3306)/orcadt?charset=utf8")
	if err != nil {
		panic(err.Error())
		log.WithFields(log.Fields{
			"filename": "mysql.go",
		}).Fatal("CreateClient")

		return nil, err
	}
	defer db.Close()
	return db, err
}

func QueryMysql() error {
	db, err := sql.Open("mysql", "root:123456@tcp(192.168.101.189:3306)/orcadt?charset=utf8")
	if err != nil {
		panic(err.Error())
		log.WithFields(log.Fields{
			"filename": "mysql.go",
		}).Fatal("CreateClient: ", err)

		return err
	}

	select_sql := "select * from orcadt_dev_login"
	rows, err1 := db.Query(select_sql)
	defer rows.Close()
	if err1 != nil {
		log.WithFields(log.Fields{
			"filename": "mysql.go",
		}).Fatal("QueryMysql: ", err1)
		return err1
	}

	for rows.Next() {
		var id, dev_id, ctime string
		var cid int
		var mtime interface{}

		err = rows.Scan(&id, &dev_id, &cid, &ctime, &mtime)
		if err != nil {
			log.WithFields(log.Fields{
				"filename": "mysql.go",
			}).Fatal("QueryMysql: ", err)
		}
		log.WithFields(log.Fields{
			"filename": "mysql.go",
			"id":       id,
			"dev_id":   dev_id,
		}).Info("QueryMysql")
	}

	defer db.Close()
	return nil
}

func QueryRawbyte() {
	// Open database connection
	db, err := sql.Open("mysql", "root:123456@tcp(192.168.101.189:3306)/orcadt?charset=utf8")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	// Execute the query
	rows, err := db.Query("select * from orcadt_dev_login")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	// Make a slice for the values
	values := make([]sql.RawBytes, len(columns))

	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	// See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	// Fetch rows
	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}

		// Now do something with the data.
		// Here we just print each column as a string.
		var value string
		for i, col := range values {
			// Here we can check if the value is nil (NULL value)
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			fmt.Println(columns[i], ": ", value)
		}
		fmt.Println("-----------------------------------")
	}
	if err = rows.Err(); err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
}
