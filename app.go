package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

const (
	dbhost = "DBHOST"
	dbport = "DBPORT"
	dbuser = "DBUSER"
	dbpass = "DBPASS"
	dbname = "DBNAME"
)

type individual struct {
	ID               int    `json:"id"`
	ShopID           int    `json:"shop_id"`
	ShopName         string `json:"shop_name"`
	Region           string `json:"region"`
	TotalPoints      int    `json:"total_points"`
	PreviousPoints   int    `json:"previous_points"`
	CurrentPosition  int    `json:"current_position"`
	PreviousPosition int    `json:"previous_position"`
	ProcessedDTTM    string `json:"last_change"`
}

type region struct {
	ID               int    `json:"id"`
	Region           string `json:"region"`
	TotalPoints      int    `json:"total_points"`
	PreviousPoints   int    `json:"previous_points"`
	CurrentPosition  int    `json:"current_position"`
	PreviousPosition int    `json:"previous_position"`
	ProcessedDTTM    string `json:"last_change"`
}

type individuals struct {
	Individuals []individual `json:"individuals"`
}

type regions struct {
	Regions []region `json:"individuals"`
}
type totalData struct {
	TotalIndividual int `json:"individual"`
	TotalRegion     int `json:"region"`
}
type header struct {
	TotalData   totalData `json:"total_data"`
	ProcessTime float64   `json:"process_time"`
}

type data struct {
	Individual []individual `json:"individual"`
	Region     []region     `json:"region"`
}

type response struct {
	Head header `json:"header"`
	Data data   `json:"data"`
}

func main() {
	initDb()
	defer db.Close()
	http.HandleFunc("/microsite/v1/finalist", FinalistHandler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

// func inResponses(w http.ResponseWriter) {
// 	users := individuals{}
// 	regions := regions{}
// 	header := header{}
// 	totalData := totalData{}

// 	errIndividual := queryIndividual(&users)
// 	if errIndividual != nil {
// 		http.Error(w, errIndividual.Error(), 500)
// 		return
// 	}

// 	errRegion := queryRegion(&regions)
// 	if errRegion != nil {
// 		http.Error(w, errRegion.Error(), 500)
// 		return
// 	}

// 	datas := data{Individual: users.Individuals, Region: regions.Regions}
// }

// FinalistHandler
func FinalistHandler(w http.ResponseWriter, req *http.Request) {
	start := time.Now()
	users := individuals{}
	regions := regions{}

	errIndividual := queryIndividual(&users)
	if errIndividual != nil {
		http.Error(w, errIndividual.Error(), 500)
		return
	}

	errRegion := queryRegion(&regions)
	if errRegion != nil {
		http.Error(w, errRegion.Error(), 500)
		return
	}

	elapsed := time.Since(start)
	totalData := totalData{TotalIndividual: len(users.Individuals), TotalRegion: len(regions.Regions)}
	datas := data{Individual: users.Individuals, Region: regions.Regions}
	header := header{TotalData: totalData, ProcessTime: float64(math.Floor(elapsed.Seconds()*100) / 100)}
	response := response{Head: header, Data: datas}

	out, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Fprintf(w, string(out))
}

// queryIndividual
func queryIndividual(users *individuals) error {
	rows, err := db.Query(`SELECT id,shop_id,shop_name,region,total_points,previous_points,current_position,previous_position,processed_dttm FROM personal ORDER BY id ASC`)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		user := individual{}
		err = rows.Scan(
			&user.ID,
			&user.ShopID,
			&user.ShopName,
			&user.Region,
			&user.TotalPoints,
			&user.PreviousPoints,
			&user.CurrentPosition,
			&user.PreviousPosition,
			&user.ProcessedDTTM)

		if err != nil {
			return err
		}
		users.Individuals = append(users.Individuals, user)
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	return nil
}

// queryIndividual
func queryRegion(locations *regions) error {
	rows, err := db.Query(`SELECT id,region,total_points,previous_points,current_position,previous_position,processed_dttm FROM region ORDER BY id ASC`)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		rgn := region{}
		err = rows.Scan(
			&rgn.ID,
			&rgn.Region,
			&rgn.TotalPoints,
			&rgn.PreviousPoints,
			&rgn.CurrentPosition,
			&rgn.PreviousPosition,
			&rgn.ProcessedDTTM)

		if err != nil {
			return err
		}
		locations.Regions = append(locations.Regions, rgn)
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	return nil
}

// initDB
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

// dbConfig
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
