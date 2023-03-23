package main

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// struct yg dibikin di db, akan dibikin tabel sama db
type Product struct {
	gorm.Model
	Nama  string
	Harga uint
}

func main() {
	dsn := "host=localhost user=postgres password=capslock dbname=pertemuan3 port=5432 TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// ini kyk buat table baru, berdasarkan nama product ditambahin s
	// jadi nnt table nya jadi Products.
	// kalo mau ganti nnt di cari di google: gorm custom table name
	err = db.AutoMigrate(&Product{})
	if err != nil {
		panic("failed to migrate")
	}

	// Create
	db.Create(&Product{Nama: "Buku", Harga: 10000})
	db.Create(&Product{Nama: "Pulpen", Harga: 10000})
	db.Create(&Product{Nama: "Tas", Harga: 100000})

	// Read
	var product Product

	db.First((&product), 1)
	fmt.Println("db.First 1 :", product)

	db.First((&product), "nama = ?", "Tas")
	fmt.Println("db.First 2 :", product)

	// Update
	db.Model(&product).Update("Harga", 20000)
	fmt.Println("db.Model.Update:", product)

	db.Model(&product).Updates(Product{Harga: 4000, Nama: "Cilok"})
	fmt.Println("db.Model.Update:", product)

	// Delete
	db.Delete(&product, 1)

	dbSQL, err := db.DB()
	if err != nil {
		panic("failed to get db")
	}

	dbSQL.Close()
}
