package main

import (
    "bufio"
    "fmt"
    "os"
    "flag"
    "time"
    "io/ioutil"
    "log"
    "github.com/yeka/zip"
)

func main() {

    start := time.Now()

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
    r, err := zip.OpenReader(*filePtr)
    if err != nil {
        log.Fatal(err)
    }
    defer r.Close()

    file, _ := os.Open(*wordlistPtr)
    fscanner := bufio.NewScanner(file)
    for fscanner.Scan() {
        for _, f := range r.File {
	    f.SetPassword(fscanner.Text())
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
	    log.Printf("\nSize of %v: %v byte(s)\n", f.Name, len(buf))
            log.Printf("============= Found password: %s in %s =============", fscanner.Text(), elapsed)
            os.Exit(0)
        }
    }
} 
