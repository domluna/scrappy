package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gocolly/colly/v2"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: scrappy <url> <element>")
		os.Exit(1)
	}

	url := os.Args[1]
	element := os.Args[2]

	dbPath := filepath.Join(os.Getenv("HOME"), ".scrappy", "scrappy_notes.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create table if not exists
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS notes (
                url TEXT PRIMARY KEY,
                content TEXT
        )`)
	if err != nil {
		log.Fatal(err)
	}

	// Check if entry already exists
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM notes WHERE url = ?", url).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}

	if count > 0 {
		fmt.Print("An entry for this URL already exists. Do you want to overwrite it? (y/n): ")
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" {
			fmt.Println("Operation cancelled.")
			return
		}
	}

	// Initialize Colly collector
	c := colly.NewCollector()

	var content string
	c.OnHTML(element, func(e *colly.HTMLElement) {
		content += e.Text + "\n"
	})

	// Start scraping
	err = c.Visit(url)
	if err != nil {
		log.Fatal(err)
	}

	// Save or update content in database
	_, err = db.Exec("INSERT OR REPLACE INTO notes (url, content) VALUES (?, ?)", url, content)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Content saved successfully.")
}
