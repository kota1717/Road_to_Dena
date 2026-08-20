package main

import (
	"database/sql"
	"fmt"
	"os"
	"log"
	_ "github.com/lib/pq"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:55432/road_to_dena?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("接続情報の設定に失敗しました: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("DBへの接続に失敗しました: %v", err)
	}
	
	fmt.Println("DB接続成功！これよりマスターデータの投入処理に入ります。")

	truncateSQL := `TRUNCATE TABLE users, teams, games, 
					reservations, seats, tickets RESTART IDENTITY CASCADE;`
	_, err = db.Exec(truncateSQL)
	if err != nil {
		log.Fatalf("テーブルの初期化に失敗しました: %v", err)
	}

	fmt.Println("テーブルを初期化し、IDの連番をリセットしました。")
}