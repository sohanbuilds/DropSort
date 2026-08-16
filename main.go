package main

import (
	"fmt"
	"os"
	s "strings"
	"sync"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func worker(files <-chan os.DirEntry, filePath string) {
	for file := range files {
		parts := s.Split(file.Name(), ".")
		last := parts[len(parts)-1]
		if !file.IsDir() {
			// It is not a directory, so move the file based on the "last"
			//create a directory with "last" first
			dirpath := filePath + "/" + last
			err := os.MkdirAll(dirpath, 0755)
			check(err)
			//move the file
			err = os.Rename(filePath+"/"+file.Name(), filePath+"/"+last+"/"+file.Name())
			check(err)
		}
	}
}

func main() {
	fmt.Println("DropSort starting...")
	filePath := "/Users/sjulakanti/Desktop/test"
	c, err := os.ReadDir(filePath)
	check(err)
	files := make(chan os.DirEntry)
	var wg sync.WaitGroup
	for w := 1; w <= 3; w++ {
		wg.Go(func() {
			worker(files, filePath)
		})
	}
	for _, entry := range c {
		files <- entry
	}
	close(files)
	wg.Wait()
}
