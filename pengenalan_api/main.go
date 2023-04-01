package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"todo/functionality"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type List struct {
	ID     uint64 `json:"id"`
	Status bool   `json:"status"`
	Item   string `json:"item"`
}

// CREATE TABLE tdlist ( ID serial PRIMARY KEY, Status BOOLEAN NOT NULL, Item VARCHAR(100) NOT NULL);

func InsertToDB(db *sql.DB, list List) (*List, error) {
	rows, err := db.Query("INSERT INTO tdlist (status, item) VALUES ($1, $2) RETURNING id, status, item", list.Status, list.Item)
	if err != nil {
		return nil, err
	}

	rows.Next()
	result := List{}
	rows.Scan(&result.ID, &result.Status, &result.Item)
	return &result, nil
}

func GetAllFromDB(db *sql.DB) ([]List, error) {
	var result []List

	rows, err := db.Query("SELECT * FROM tdlist ORDER BY id")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		list := List{}
		rows.Scan(&list.ID, &list.Status, &list.Item)
		result = append(result, list)
	}

	return result, nil
}

func GetByIDFromDB(db *sql.DB, id uint64) (*List, error) {
	rows, err := db.Query("SELECT * FROM tdlist WHERE id=$1", id)
	if err != nil {
		return nil, err
	}

	fmt.Println(rows)
	rows.Next()

	result := List{}
	rows.Scan(&result.ID, &result.Status, &result.Item)
	if result.ID == 0 {
		return nil, errors.New("Data tidak ditemukan")
	}
	fmt.Println(rows)
	fmt.Println(result)
	return &result, nil
}

func DeleteByIdFromDB(db *sql.DB, id uint64) (*List, error) {
	rows, err := db.Query("DELETE FROM tdlist WHERE id=$1 RETURNING id, status, item", id)
	if err != nil {
		return nil, err
	}

	rows.Next()
	result := List{}
	rows.Scan(&result.ID, &result.Status, &result.Item)
	if result.ID == 0 {
		return nil, errors.New("Data tidak ditemukan")
	}
	return &result, nil
}

func ToggleByIdFromDB(db *sql.DB, id uint64) (*List, error) {
	rows, err := db.Query("UPDATE tdlist SET status = NOT status WHERE ID = $1 RETURNING id, status, item", id)
	if err != nil {
		return nil, err
	}

	rows.Next()
	result := List{}
	rows.Scan(&result.ID, &result.Status, &result.Item)
	if result.ID == 0 {
		return nil, errors.New("Data tidak ditemukan")
	}
	return &result, nil
}

func main() {
	db, err := functionality.ConnectToDB()
	if err != nil {
		fmt.Println(err)
	}

	err = db.Ping()
	if err != nil {
		fmt.Println(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Halo RPL")
	})

	// INSERT
	// /insert
	http.HandleFunc("/insert", func(w http.ResponseWriter, r *http.Request) {
		newList := List{}

		err := json.NewDecoder(r.Body).Decode(&newList)
		if err != nil {
			fmt.Println(err)
		}

		jsonMap := map[string]interface{}{}
		fmt.Println(newList.Item)
		if newList.Item != "" {
			res, err := InsertToDB(db, newList)
			if err != nil {
				fmt.Println(err)
			}
			jsonMap["data"] = res
			jsonMap["message"] = "berhasil ditambahkan"
		} else {
			jsonMap["data"] = nil
			jsonMap["message"] = "tidak ada data yang ditambahkan"
		}

		b, err := json.Marshal(jsonMap)
		if err != nil {
			fmt.Println(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, string(b))
	})

	// GET ALL
	// /getall
	http.HandleFunc("/getall", func(w http.ResponseWriter, r *http.Request) {
		res, err := GetAllFromDB(db)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(res)

		jsonMap := map[string]interface{}{
			"data": res,
		}

		b, err := json.Marshal(jsonMap)
		if err != nil {
			fmt.Println(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, string(b))
	})

	// GET BY ID
	// get/:id
	http.HandleFunc("/get/", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(r.URL.Path[len("/get/"):], 10, 64)
		if err != nil {
			fmt.Println(err)
		}

		res, err := GetByIDFromDB(db, id)
		if err != nil {
			fmt.Println(err)
		}

		jsonMap := map[string]interface{}{}
		if res != nil {
			jsonMap["data"] = res
			jsonMap["message"] = "data ditemukan"
		} else {
			jsonMap["data"] = nil
			jsonMap["message"] = "data tidak ditemukan"
		}

		b, err := json.Marshal(jsonMap)
		if err != nil {
			fmt.Println(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, string(b))
	})

	// DELETE BY ID
	// /delete/:id
	http.HandleFunc("/delete/", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(r.URL.Path[len("/delete/"):], 10, 64)
		if err != nil {
			fmt.Println(err)
		}
		res, err := DeleteByIdFromDB(db, id)
		if err != nil {
			fmt.Println(err)
		}

		jsonMap := map[string]interface{}{}
		if res != nil {
			jsonMap["data"] = res
			jsonMap["message"] = "data berhasil dihapus"
		} else {
			jsonMap["data"] = nil
			jsonMap["message"] = "data gagal dihapus, data tidak ditemukan dalam database"
		}

		b, err := json.Marshal(jsonMap)
		if err != nil {
			fmt.Println(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, string(b))
	})

	// TOGlE BY ID
	// /toggle/:id
	http.HandleFunc("/toggle/", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(r.URL.Path[len("/toggle/"):], 10, 64)
		if err != nil {
			fmt.Println(err)
		}
		res, err := ToggleByIdFromDB(db, id)
		if err != nil {
			fmt.Println(err)
		}

		jsonMap := map[string]interface{}{}
		if res != nil {
			jsonMap["data"] = res
			jsonMap["message"] = "berhasil mengubah status"
		} else {
			jsonMap["data"] = nil
			jsonMap["message"] = "gagal mengubah status, data tidak ditemukan"
		}

		b, err := json.Marshal(jsonMap)
		if err != nil {
			fmt.Println(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, string(b))
	})

	fmt.Println("starting web server at http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}
