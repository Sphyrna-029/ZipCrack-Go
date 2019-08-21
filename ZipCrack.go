package main

import (
    "bufio"
    "fmt"
    "os"
    //"io"
    "time"
    "log"
    "github.com/alexmullins/zip"
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


    fmt.Println("\nContending secretdata.zip against wordlist.txt")
    r, err := zip.OpenReader("secretdata.zip")
    if err != nil {
        log.Fatal(err)
    }
    defer r.Close()

    file, _ := os.Open("wordlist.txt")
    fscanner := bufio.NewScanner(file)
    for fscanner.Scan() {
        for _, f := range r.File {
	    f.SetPassword(fscanner.Text())
            rc, err := f.Open()
            if err != nil {
                continue
            }
            rc.Close()
            r.Close()
            elapsed := time.Since(start)
            log.Printf("Found password: %s in %s", fscanner.Text(), elapsed)
            os.Exit(0)
        }
    }
} 
