package main

import (
	"fmt"
	"log"
	"os"
	"time"
	"database/sql"

	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql"

)

var (
	db *sqlx.DB
)

func init() {
	db_host := os.Getenv("ISUBATA_DB_HOST")
	if db_host == "" {
		db_host = "127.0.0.1"
	}
	db_port := os.Getenv("ISUBATA_DB_PORT")
	if db_port == "" {
		db_port = "3306"
	}
	db_user := os.Getenv("ISUBATA_DB_USER")
	if db_user == "" {
		db_user = "isucon"
	}
	db_password := os.Getenv("ISUBATA_DB_PASSWORD")
	if db_password == "" {
		db_password = ":" + "isucon"
	}

	dsn := fmt.Sprintf("%s%s@tcp(%s:%s)/isubata?parseTime=true&loc=Local&charset=utf8mb4",
		db_user, db_password, db_host, db_port)

	log.Printf("Connecting to db: %q", dsn)
	db, _ = sqlx.Connect("mysql", dsn)
	for {
		err := db.Ping()
		if err == nil {
			break
		}
		log.Println(err)
		time.Sleep(time.Second * 3)
	}

	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(5 * time.Minute)

	log.Printf("Succeeded to connect db.")
}

func main() {
	for i:= 0; i < 1500; i++ {
		var name string
		var data []byte
		err := db.QueryRow("SELECT name, data FROM image WHERE id = ?", i).Scan(&name, &data)
		if err == sql.ErrNoRows {
			continue
		}
		if err != nil {
			log.Fatalln(err)
		}

		out, _ := os.Create(fmt.Sprintf("/home/isucon/isubata/webapp/public/icons/%s", name))
		if err != nil {
			fmt.Println("write error:", err)
		}
		_, err = out.Write(data)
		if err != nil {
			fmt.Println("write error:", err)
		}

		out.Close()
	}

}
