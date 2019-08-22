package main

import (
    "bufio"
    "fmt"
    "os"
    "flag"
    "time"
    "log"
    "sync"
    "github.com/alexmullins/zip"
)

var wg sync.WaitGroup 

func decrypt(zipFile string, password string, passChan chan string) {
    r, err := zip.OpenReader(zipFile)
    if err != nil {
        log.Fatal(err)
    }
	
    f := r.File 
    f.SetPassword(password)
    rc, err := f.Open()
    if err != nil {
        close(passChan)
    }
	//If solved clean up our files and send back password
    passChan <- password
    rc.Close()
    r.Close()
    close(passChan)
}

func main() {
    start := time.Now()

    //Parse cmd line arguments
    filePtr := flag.String("zip", "", "Path to zip file.")
    wordlistPtr := flag.String("wordlist", "", "Path to wordlist.")
    flag.Parse()

    if *filePtr == "" {
        log.Fatal("Zip file not found.")
    }
    if *wordlistPtr == "" {
        log.Fatal("Wordlist not found.")
    }
	
    //Display art, info
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
    log.Printf("Contending %s against %s", *filePtr, *wordlistPtr)
	
    //Initialize channel for sending solved password back to the main thread
    passChan := make(chan string)

    //Begin the main loop
    file, _ := os.Open(*wordlistPtr)
    fscanner := bufio.NewScanner(file)
    for fscanner.Scan() {
	go decrypt(*filePtr, fscanner.Text(), passChan)

	passwd := <-passChan

	if passwd != "" {
	    close(passChan)
	    elapsed := time.Since(start)
            log.Printf("Found password: %s in %s", passwd, elapsed)
            os.Exit(0)
	}
    }
} 
