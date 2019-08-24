package main

import (
    "bufio"
    "fmt"
    "os"
    "flag"
    "io/ioutil"
    "log"
    "bytes"
    "encoding/json"
    "errors"
    "net/http"
    "time"
    "github.com/alexmullins/zip"
)

var start = time.Now()

type SlackRequestBody struct {
    Text string `json:"text"`
}

func slack(webhookUrl string, msg string) error {
    //From https://golangcode.com/send-slack-messages-without-a-library/
    slackBody, _ := json.Marshal(SlackRequestBody{Text: msg})
    req, err := http.NewRequest(http.MethodPost, webhookUrl, bytes.NewBuffer(slackBody))
    if err != nil {
        return err
    }

    req.Header.Add("Content-Type", "application/json")

    client := &http.Client{Timeout: 10 * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }

    buf := new(bytes.Buffer)
    buf.ReadFrom(resp.Body)

    if buf.String() != "ok" {
        return errors.New("Non-ok response returned from Slack")
    }
    return nil
}

func decrypt(r *zip.ReadCloser, password string, slackHook *string) {

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
        message := fmt.Sprintf(":heavy_check_mark: Found password *%s* for *%s* in *%s* :heavy_check_mark:", password, f.Name, elapsed)
        if *slackHook != "" {
            log.Printf("Sending solution to slack.")
            slack(*slackHook, message)
            if err != nil {
                log.Fatal(err)
            }
        }
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
    slackHook := flag.String("slack", "", "Slack web hook url. (Optional)")
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

    r, err := zip.OpenReader(*filePtr)
    if err != nil {
        log.Fatal(err)
    }

    for fscanner.Scan() {
        go decrypt(r, fscanner.Text(), slackHook)
    }
    log.Printf("Password not found.")
    message := ":no_entry_sign: Password not found :no_entry_sign:"
    if *slackHook != "" {
        slack(*slackHook, message)
        if err != nil {
            log.Fatal(err)
        }
    }
}