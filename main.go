package main

import (
	// "fmt"
	"github.com/joho/godotenv"
)

func main() {

	// 環境変数読み込み
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}

	// チェックマーク済み（読了済み）の漫画を、Notionから削除

	titles := getComicTitles()

	comics := scrape(titles)

	insertComics(comics)
}
