package main

import (
    "bufio"
    "fmt"
    "os"
    "flag"
	"io/ioutil"
	"log"
	"time"
    "github.com/yeka/zip"
)

var start = time.Now()

func decrypt(zipFile string, password string) {
	r, err := zip.OpenReader(zipFile)
    if err != nil {
        log.Fatal(err)
    }
	defer r.Close()
	
    for _, f := range r.File {
		f.SetPassword(password)
		rc, err := f.Open()
		if err != nil {
			continue
		}
		buf, err := ioutil.ReadAll(rc)
		if err != nil {
			continue
		}
		rc.Close()
		r.Close()
		elapsed := time.Since(start)
		log.Printf("Size of %v: %v byte(s)\n", f.Name, len(buf))
		log.Printf("!============= Found password: %s in %s =============!", password, elapsed)
		os.Exit(0)
	}
}

func main() {

    banner := 
    `
    ▒███████▒ ██▓ ██▓███   ▄████▄   ██▀███   ▄▄▄       ▄████▄   ██ ▄█▀
    ▒ ▒ ▒ ▄▀░▓██▒▓██░  ██▒▒██▀ ▀█  ▓██ ▒ ██▒▒████▄    ▒██▀ ▀█   ██▄█▒ 
    ░ ▒ ▄▀▒░ ▒██▒▓██░ ██▓▒▒▓█    ▄ ▓██ ░▄█ ▒▒██  ▀█▄  ▒▓█    ▄ ▓███▄░ 
      ▄▀▒   ░░██░▒██▄█▓▒ ▒▒▓▓▄ ▄██▒▒██▀▀█▄  ░██▄▄▄▄██ ▒▓▓▄ ▄██▒▓██ █▄ 
    ▒███████▒░██░▒██▒ ░  ░▒ ▓███▀ ░░██▓ ▒██▒ ▓█   ▓██▒▒ ▓███▀ ░▒██▒ █▄
    ░▒▒ ▓░▒░▒░▓  ▒▓▒░ ░  ░░ ░▒ ▒  ░░ ▒▓ ░▒▓░ ▒▒   ▓▒█░░ ░▒ ▒  ░▒ ▒▒ ▓▒
    ░░▒ ▒ ░ ▒ ▒ ░░▒ ░       ░  ▒     ░▒ ░ ▒░  ▒   ▒▒ ░  ░  ▒   ░ ░▒ ▒░
    ░ ░ ░ ░ ░ ▒ ░░░       ░          ░░   ░   ░   ▒   ░        ░ ░░ ░ 
      ░ ░     ░           ░ ░         ░           ░  ░░ ░      ░  ░   
    ░                     ░                           ░      
    `
    fmt.Println(banner)

    filePtr := flag.String("zip", "", "Path to zip file.")
    wordlistPtr := flag.String("wordlist", "", "Path to wordlist.")
    flag.Parse()

    if *filePtr == "" {
        log.Fatal("Zip file not found.")
    }
    if *wordlistPtr == "" {
        log.Fatal("Wordlist not found.")
    }

    log.Printf("Contending %s against %s", *filePtr, *wordlistPtr)

    file, _ := os.Open(*wordlistPtr)
    fscanner := bufio.NewScanner(file)
    for fscanner.Scan() {
        go decrypt(*filePtr, fscanner.Text())
	}
} 
