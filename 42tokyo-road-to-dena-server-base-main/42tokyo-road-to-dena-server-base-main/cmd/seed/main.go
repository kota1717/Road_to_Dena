package main

import (
	"database/sql"
	"fmt"
	"os"
	"log"
	_ "github.com/lib/pq"
)

func main() {
	db := connectDB()
	defer db.Close()

	truncateTables(db)

	teamIDs := seedTeams(db)
}

func connectDB() *sql.DB {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:55432/road_to_dena?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("接続設定エラー: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("DB接続エラー: %v", err)
	}
	return db
}

func truncateTables(db *sql.DB) {
	truncateSQL := `TRUNCATE TABLE users, teams, games, 
					reservations, seats, tickets RESTART IDENTITY CASCADE;`
	if _, err := db.Exec(truncateSQL); err != nil {
		log.Fatalf("テーブル初期化失敗: %v", err)
	}
	fmt.Println("テーブル初期化")
}

func seedTeams(db *sql.DB) {
	teamNames := []string{
			"横浜DeNAベイスターズ", "阪神タイガース", "読光ジャイアンツ", 
    		"広島東洋カープ", "東京ヤクルトスワローズ", "中日ドラゴンズ",
    		"オリックス・バファローズ", "千葉ロッテマリーンズ", "福岡ソフトバンクホークス",
    		"東北楽天ゴールデンイーグルス", "北海道日本ハムファイターズ", "埼玉西武ライオンズ",
	}

	var teamIDs []int64

	for _, name := range teamNames{
		var id int64
		query := `INSERT INTO teams (name) VALUES ($1) RETURNING id;`

		err := db.QueryRow(query, name).Scan(&id)
		if err != nil {
			log.Fatalf("チームデータの挿入に失敗しました (%s): %v", name, err)
		}
		
		teamIDs = append(teamIDs, id)
	}

	fmt.Printf("%d 件のチームデータを投入しました。 \n", len(teamIDs))
	return teamIDs
}