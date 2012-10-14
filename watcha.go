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
	interval        = flag.Int("t", 30, "Interval time to sync (seconds)")
	syncEntriesFile = flag.String("f", "", "Sync entries file")
	interactive     = flag.Bool("i", false, "Specifies if it wants an interactive session")
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
		if !*interactive {
			for {
				time.Sleep(time.Hour)
			}
		}
	}

	if *interactive {
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
	for {
		var option string
		fmt.Println("1. (l)ist, 2. (a)dd, 3. (r)emove, 4. (q)uit")
		fmt.Print("What to do? > ")
		fmt.Scanf("%s", &option)
		switch option {
		case "1", "l":
			list(files)
		case "2", "a":
			add(files)
		case "3", "r":
			fmt.Println("What!? Removing a file that was already syncing")
		case "4", "q":
			fmt.Println("Bye.")
			os.Exit(0) // Probably clean up before exit to not leave corrupted sync files
		default:
			fmt.Println("Incorrect option")
		}

	}
}

func list(files []syncFile) {
	if len(files) == 0 {
		fmt.Println("No files are actually being synced. Please add one")
	} else {
		for i, file := range files {
			fmt.Println(string(i) + ". " + file.String())
		}
	}
}

func add(files []syncFile) {
	var src, dst string

	fmt.Println("Please src and dst path separated by a space:")
	fmt.Scanf("%s %s", &src, &dst)
	files = append(files, syncFile{src, dst})
}

func performSync(files []syncFile) {
	for _, sf := range files {
		sf.sync()
	}
}
