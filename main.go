package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charleshuang3/subget/lib/subcat"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/log"
)

var (
	searchKeyword = flag.String("s", "", "(Optional) the search keyword")
	subLang       = flag.String(
		"subl",
		envOrDefault("SUB_LANG", "zh-CN"),
		"(Optional) language code used in download, can also pass by env var SUB_LANG. Default is Chinese Simplified")
	vidLang = flag.String(
		"vidl",
		envOrDefault("VIDEO_LANG", "简体中文.chi"),
		"(Optional) language code used in rename, can also pass by env var VIDEO_LANG. Default is Chinese Simplified")
)

func envOrDefault(key string, defaultValue string) string {
	if os.Getenv(key) != "" {
		return os.Getenv(key)
	}
	return defaultValue
}

func main() {
	flag.Parse()

	keyword := getSearchKeyword()
	searchResult := search(keyword)
	sr := selectResult(searchResult)
	subPath := openDownloadPage(sr.Path)
	content := downloadSub(subPath)
	name := renameSub(keyword, subPath)
	writeSubToFile(name, content)
	fmt.Printf("sub file is saved: %s", name)
}

func getSearchKeyword() string {
	if *searchKeyword != "" {
		return *searchKeyword
	}
	var word string
	err := huh.NewInput().
		Title("Enter search keyword").
		Value(&word).
		Run()
	if err != nil {
		log.Fatalf("Error getting input: %v", err)
	}
	if word == "" {
		log.Fatal("Search keyword cannot be empty.")
	}
	return word
}

func search(keyword string) []subcat.SearchResult {
	searchResult := []subcat.SearchResult{}
	searchAction := func() {
		q := strings.ReplaceAll(keyword, " ", "+")
		searchResult = subcat.Search(q)
	}

	err := spinner.New().
		Title(fmt.Sprintf("Searching %q on SubtitleCat", keyword)).
		Action(searchAction).
		Run()

	if err != nil {
		log.Fatalf("Error wait for searching: %v", err)
	}

	return searchResult
}

const (
	// usually subcat can not return more than 5 useful result.
	searchLimit = 10

	searchResultSelectFmt = "%s ↓ %d"
)

func selectResult(results []subcat.SearchResult) subcat.SearchResult {
	opts := []huh.Option[subcat.SearchResult]{}
	for i, sr := range results {
		if i >= searchLimit {
			break
		}
		opts = append(opts, huh.NewOption(
			fmt.Sprintf(searchResultSelectFmt, sr.Title, sr.Downloads),
			sr,
		))
	}

	var selected subcat.SearchResult
	err := huh.NewSelect[subcat.SearchResult]().
		Title("Select a search result").
		Options(
			opts...,
		).
		Value(&selected).
		Run()

	if err != nil {
		log.Fatalf("Error user selecting search result: %v", err)
	}
	return selected
}

func openDownloadPage(path string) string {
	subPath := ""
	action := func() {
		subPath = subcat.OpenDownloadPage(path, *subLang)
	}

	err := spinner.New().
		Title(fmt.Sprintf("Open Download page %s for lang code %s", path, *subLang)).
		Action(action).
		Run()

	if err != nil {
		log.Fatalf("Error wait for open download page: %v", err)
	}

	return subPath
}

func downloadSub(path string) []byte {
	var sub []byte
	action := func() {
		sub = subcat.DownloadSubFile(path)
	}

	err := spinner.New().
		Title(fmt.Sprintf("Download sub file %s", path)).
		Action(action).
		Run()

	if err != nil {
		log.Fatalf("Error wait for download sub file: %v", err)
	}

	return sub
}

func renameSub(keyword string, subPath string) string {
	_, fn := filepath.Split(subPath)

	ext := filepath.Ext(fn)
	n := keyword + "." + *vidLang + ext

	err := huh.NewInput().
		Title("Enter file name").
		Value(&n).
		Placeholder(n).
		Run()
	if err != nil {
		log.Fatalf("Error getting input: %v", err)
	}
	if n == "" {
		log.Fatal("File name cannot be empty.")
	}

	return n
}

func writeSubToFile(filename string, content []byte) {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("open file failed: %v", err)
	}
	defer f.Close()

	_, err = f.Write(content)
	if err != nil {
		log.Fatalf("write file failed: %v", err)
	}
}
