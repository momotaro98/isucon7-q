package main

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"strings"
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

type Image struct {
	Name  string `db:"name"`
	Image []byte `db:"name"`
}

func main() {
	//dest := []Image{}
	// err := db.Select(&dest, "SELECT name, data FROM image")
	// if err != nil {
	// 	panic(err)
	// }

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

		img, _, err := image.Decode(bytes.NewReader(data))
		if err != nil {
			fmt.Println("image.Decode error:", err)
			continue
		}
		switch true {
		case strings.HasSuffix(name, ".jpg"), strings.HasSuffix(name, ".jpeg"):
			out, _ := os.Create(fmt.Sprintf("/home/isucon/isubata/webapp/public/icons/%s", name))
			err := jpeg.Encode(out, img, nil)
			if err != nil {
				fmt.Println("jpg error:", err)
			}
			out.Close()
		case strings.HasSuffix(name, ".gif"):
			out, _ := os.Create(fmt.Sprintf("/home/isucon/isubata/webapp/public/icons/%s", name))
			err := gif.Encode(out, img, nil)
			if err != nil {
				fmt.Println("gif error:", err)
			}
			out.Close()
		case strings.HasSuffix(name, ".png"):
			out, _ := os.Create(fmt.Sprintf("/home/isucon/isubata/webapp/public/icons/%s", name))
			err := png.Encode(out, img)
			if err != nil {
				fmt.Println("png error:", err)
			}
			out.Close()
		default:
			log.Fatalln("unknown format")
		}
	}

}
