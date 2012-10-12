package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

var (
	interval        = flag.Int("i", 30, "Interval time to sync (seconds)")
	syncEntriesFile = flag.String("f", "", "Sync entries file")
)

func main() {
	flag.Parse()

	// Initial creation of sync files slice
	files := make([]syncFile, 0, 5)

	ticker := time.NewTicker(time.Second * 2)
	go func() {
		for _ = range ticker.C {
			performSync(files)
		}
	}()

	if *syncEntriesFile != "" {
		files, _ = readFromFile(*syncEntriesFile)
		for {
			time.Sleep(time.Hour)
		}
	} else {
		enterRepl(files)
	}
}

// Reads entries from file
func readFromFile(syncEntriesFileName string) ([]syncFile, error) {
	f, err := os.Open(syncEntriesFileName)
	if err != nil {
		return nil, err
	}

	files := make([]syncFile, 0, 5)

	reader := bufio.NewReader(f)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}

		pairs := strings.SplitN(string(line), " ", 2)
		files = append(files, syncFile{pairs[0], pairs[1]})
	}

	return files, nil
}

func enterRepl(files []syncFile) {
	var src, dst string
	for {
		fmt.Println("Do you want to add a path?")
		if n, err := fmt.Scanf("%s %s", &src, &dst); err == nil && n == 2 {
			files = append(files, syncFile{src, dst})
		}
	}
}

func performSync(files []syncFile) {
	for _, sf := range files {
		sf.sync()
	}
}
