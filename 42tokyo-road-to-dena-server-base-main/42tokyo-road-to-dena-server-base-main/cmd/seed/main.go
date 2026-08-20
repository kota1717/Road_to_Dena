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

	teamNames := []string{
		"横浜DeNAベイスターズ", "阪神タイガース", "読光ジャイアンツ", 
    	"広島東洋カープ", "東京ヤクルトスワローズ", "中日ドラゴンズ",
    	"オリックス・バファローズ", "千葉ロッテマリーンズ", "福岡ソフトバンクホークス",
    	"東北楽天ゴールデンイーグルス", "北海道日本ハムファイターズ", "埼玉西武ライオンズ",
	}

	var teamIDs []int64

	truncateSQL := `TRUNCATE TABLE users, teams, games, 
					reservations, seats, tickets RESTART IDENTITY CASCADE;`
	_, err = db.Exec(truncateSQL)
	if err != nil {
		log.Fatalf("テーブルの初期化に失敗しました: %v", err)
	}

	fmt.Println("テーブルを初期化し、IDの連番をリセットしました。")

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
}