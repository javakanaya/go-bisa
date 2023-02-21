package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type User struct {
	ID   uint64 `json:"id"`
	Name string `jsn:"name"`
}

func Connect() (*sql.DB, error) {
	// sql.open hasilnya db dan error
	// open pgx param ke-1 itu driver nya
	// berikutnya itu credential-nya, username root, pass capslock
	// trus konek ke localhost port 5432 slash nama datanya yaitu pgx
	db, err := sql.Open("pgx", "postgres://postgres:capslock@localhost:5432/pgx")
	// kalo error bakal return error
	if err != nil {
		return nil, err
	}

	return db, nil
}

func InsertToDB(db *sql.DB, user User) (*User, error) {
	// rows itu hasil query
	rows, err := db.Query("INSERT INTO users (name) VALUES ($1) RETURNING id, name", user.Name)
	if err != nil {
		return nil, err
	}

	rows.Next()

	result := User{}

	// id dulu baru namanya
	rows.Scan(&result.ID, &result.Name)
	return &result, nil
}

func GetAll(db *sql.DB) ([]User, error) {
	var result []User

	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		user := User{}
		rows.Scan(&user.ID, &user.Name)
		result = append(result, user)
	}

	return result, nil
}

func main() {
	// nyoba open sebuah database
	db, err := Connect()

	// ngecek ada error ndak
	if err != nil {
		fmt.Println(err)
	}

	// ping digunakan untuk konek ke db, kalo error var err bakal ada isinya
	// kalo dak error var err berati kosong/nil
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
	}

	// java := User{
	// 	Name: "Java",
	// }

	res, err := GetAll(db)
	if err != nil {
		fmt.Println(err)
	}

	// print ke webserver
	jsonMap := map[string]interface{}{
		"data": res,
	}

	b, err := json.Marshal(jsonMap)
	if err != nil {
		fmt.Println(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, string(b))
	})

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
