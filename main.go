package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	s "strings"
	"sync"
)

type Container struct {
	mu sync.Mutex
	m  map[string][]os.DirEntry
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func hash(filePath string) []byte {
	dat, err := os.Open(filePath)
	check(err)
	defer dat.Close()
	h := sha256.New()
	_, err = io.Copy(h, dat)
	check(err)
	hash := h.Sum(nil)
	return hash
}

func worker(files <-chan os.DirEntry, filePath string, c *Container) {
	m := c.m
	for file := range files {
		parts := s.Split(file.Name(), ".")
		last := parts[len(parts)-1]
		if !file.IsDir() {
			hash := hash(filePath + "/" + file.Name())
			c.mu.Lock()
			m[string(hash)] = append(m[string(hash)], file)
			l := len(m[string(hash)])
			c.mu.Unlock()
			if l > 1 {
				//create a directory with "duplicate entities"
				dirpath := filePath + "/" + string(hash)
				err := os.MkdirAll(dirpath, 0755)
				check(err)
				err = os.Rename(filePath+"/"+file.Name(), filePath+"/"+string(hash)+"/"+file.Name())
				check(err)
			} else {
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
}

func main() {
	fmt.Println("DropSort starting...")
	filePath := "/Users/sjulakanti/Desktop/test"
	c, err := os.ReadDir(filePath)
	check(err)
	files := make(chan os.DirEntry)
	container := Container{m: make(map[string][]os.DirEntry)}
	var wg sync.WaitGroup
	for w := 1; w <= 3; w++ {
		wg.Go(func() {
			worker(files, filePath, &container)
		})
	}
	for _, entry := range c {
		files <- entry
	}
	close(files)
	wg.Wait()
}
