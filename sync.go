package main

import (
	"fmt"
	"io"
	"os"
)

type syncFile struct {
	src string
	dst string
}

func (sf *syncFile) sync() (int64, error) {
	// Read from the src
	// (Over)Write to the dst
	//fmt.Printf("Syncing from '%s' to '%s'\n", sf.src, sf.dst)

	src, err := os.Open(sf.src)
	if err != nil {
		fmt.Println("There was an error opening file " + sf.src)
	}
	defer src.Close()

	dst, err := os.Create(sf.dst)
	if err != nil {
		fmt.Println("There was an error creating file " + sf.dst + ". " + err.Error())
	}
	defer dst.Close()

	return io.Copy(dst, src)
}

func (sf *syncFile) String() string {
	return "'" + sf.src + "' syncing to '" + sf.dst + "'"
}
