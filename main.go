package main

import (
	"fmt"
	"os"
	s "strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	fmt.Println("DropSort starting...")
	filePath := "/Users/sjulakanti/Desktop/test"
	c, err := os.ReadDir(filePath)
	check(err)
	fmt.Println("Listing Test dir")
	var files []string
	for _, entry := range c {
		parts := s.Split(entry.Name(), ".")
		last := parts[len(parts)-1]
		if !entry.IsDir() {
			// It is not a directory, so move the file based on the "last"
			var flag = false
			for file := range files {
				if last == files[file] {
					flag = true
				}
			}
			if !flag {
				//create a directory with "last" first
				files = append(files, last)
				dirpath := filePath + "/" + last
				err := os.Mkdir(dirpath, 0755)
				check(err)
			}
			//move the file
			err = os.Rename(filePath+"/"+entry.Name(), filePath+"/"+last+"/"+entry.Name())
			check(err)
		}
	}
}
