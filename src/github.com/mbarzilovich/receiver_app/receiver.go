package main

import (
 //   "fmt"
    "log"
    "net/http"
    "github.com/go-stomp/stomp"
    "strings"
)


func send_message(text string) {

    log.Println("Start connecting to brocker")
    conn, err := stomp.Dial("tcp", "brocker:61613")
    if err != nil {
        println("cannot connect to server", err.Error())
        return
    }
    log.Println("Connected to brocker")
    err = conn.Send("/queue/SampleQueue", "text/plain", []byte(text))
    if err != nil {
        println("failed to send to server", err)
        return
    }
    log.Println("Message sent to brocker")
    conn.Disconnect()
  
}

func handler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        r.ParseForm()
        log.Println("text:", r.Form["text"])
        go send_message(strings.Join(r.Form["text"], ", "))
    } else {
        http.Error(w, "Method not allowed", 405)
        log.Println("We can handle only POST requests")
        return
    }  
  

}



func main() {
    http.HandleFunc("/", handler)
    log.Fatal(http.ListenAndServe(":80", nil))
}