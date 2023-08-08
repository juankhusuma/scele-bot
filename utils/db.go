package utils

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

func GetDB() (*sql.DB, error) {
	cs := os.Getenv("DB_CONN")
	db, err := sql.Open("postgres", cs)
	if err != nil {
		return nil, err
	}
	return db, nil
	// defer db.Close()

	// rows, err := db.Query("select version()")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer rows.Close()

	// var version string
	// for rows.Next() {
	// 	err := rows.Scan(&version)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }
	// fmt.Printf("version=%s\n", version)
}
