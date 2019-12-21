package main

import (
	"crypto/md5"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func getUpdate(url, selector string) (hash, text string) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	text = doc.Find(selector).Text()
	hash = fmt.Sprintf("%x", md5.Sum([]byte(text)))
	return
}

func main() {
	url := os.Getenv("URL")
	selector := os.Getenv("SELECTOR")
	interval := os.Getenv("INTERVAL")
	webhook := os.Getenv("WEBHOOK")
	intervalInt, err := strconv.Atoi(interval)
	if err != nil {
		log.Fatalf("interval can not be convert to int, error: %v", err)
	}
	log.Printf("config loaded successfully.\nurl: %s\ncss selector: %s\ninterval: %d minute\nwebhook: %s",
		url, selector, intervalInt, webhook)
	lastHash, lastText := getUpdate(url, selector)
	log.Printf("target inited, %s\n==========\n||%s||\n==========", lastHash, lastText)
	for {
		time.Sleep(time.Minute * time.Duration(intervalInt))
		thisHash, thisText := getUpdate(url, selector)
		if thisHash != lastHash {
			log.Printf("target changed, %s\n==========\n||%s||\n==========", thisHash, thisText)
			lastHash = thisHash
			lastText = thisText
			http.Get(webhook)
		}
	}
}
