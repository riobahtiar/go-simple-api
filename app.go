package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

const (
	dbhost = "127.0.0.1"
	dbport = "5432"
	dbuser = "postgres"
	dbpass = "liberate123"
	dbname = "go_simple_api"
)

func main() {
	http.HandleFunc("/api/staff/", staffHandler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

type staff struct {
	ID      int
	Name    string
	Address string
	Status  int
}

type team struct {
	Team []staff
}

func staffHandler(w http.ResponseWriter, r *http.Request) {
	staff := team{}

	err := getStaff(&staff)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	out, err := json.Marshal(staff)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Fprintf(w, string(out))
}

func getStaff(staffs *team) error {
	rows, err := db.Query(`
		SELECT
			id,
			name,
			address,
			status
		FROM staff
		ORDER BY id DESC`)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		man := staff{}
		err = rows.Scan(
			&man.ID,
			&man.Name,
			&man.Address,
			&man.Status,
		)
		if err != nil {
			return err
		}
		staffs.Team = append(staffs.Team, man)
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	return nil
}

func initDb() {
	config := dbConfig()
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config[dbhost], config[dbport],
		config[dbuser], config[dbpass], config[dbname])

	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected!")
}

func dbConfig() map[string]string {
	conf := make(map[string]string)
	host, ok := os.LookupEnv(dbhost)
	if !ok {
		panic("DBHOST environment variable required but not set")
	}
	port, ok := os.LookupEnv(dbport)
	if !ok {
		panic("DBPORT environment variable required but not set")
	}
	user, ok := os.LookupEnv(dbuser)
	if !ok {
		panic("DBUSER environment variable required but not set")
	}
	password, ok := os.LookupEnv(dbpass)
	if !ok {
		panic("DBPASS environment variable required but not set")
	}
	name, ok := os.LookupEnv(dbname)
	if !ok {
		panic("DBNAME environment variable required but not set")
	}
	conf[dbhost] = host
	conf[dbport] = port
	conf[dbuser] = user
	conf[dbpass] = password
	conf[dbname] = name
	return conf
}
