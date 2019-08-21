package main

import (
    "bufio"
    "fmt"
	"os"
	"io"
	"log"
	"github.com/alexmullins/zip"
)

func main() {
    // Open a zip archive for reading.
    r, err := zip.OpenReader("encryptedFile.zip")
    if err != nil {
        log.Fatal(err)
    }
    defer r.Close()

    file, _ := os.Open("wordList.txt")
    fscanner := bufio.NewScanner(file)
    for fscanner.Scan() {
		for _, f := range r.File {
			f.SetPassword(fscanner.Text())
			rc, err := f.Open()
			if err != nil {
				//log.Printf(err.Error())
				continue
			}
			_, err = io.CopyN(os.Stdout, rc, 68)
			if err != nil {
				log.Printf(err.Error())
			}
			rc.Close()
			r.Close()
			fmt.Println("\n\nPassword found: ", fscanner.Text())
			os.Exit(0)
		}
    }
} 
