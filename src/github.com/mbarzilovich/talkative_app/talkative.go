package main

import (
 //  "fmt"
    "log"
    "net/http"
    "github.com/go-stomp/stomp"
    "github.com/gorilla/websocket"
    "text/template"
//    "sort"
)

var currentMessage = "Message will be here"
var homeTempl = template.Must(template.New("").Parse(homeHTML))
var messageChan = make(chan []byte)
var upgrader  = websocket.Upgrader{}

func receiveMessage(subscribed chan bool) {
    conn, err := stomp.Dial("tcp", "brocker:61613")
    if err != nil {
        println("cannot connect to server", err.Error())
        return
    }
    log.Println("Connected to brocker")
    sub, err := conn.Subscribe("/queue/SampleQueue", stomp.AckAuto)
    if err != nil {
        println("failed to subscribe to queue", err)
        return
    }
    log.Println("Subscribed to queue")
    close(subscribed)
    
    for  {
        msg := <- sub.C
        log.Println("Message received", string(msg.Body))
        currentMessage = string(msg.Body)
        messageChan <- msg.Body
    }
}

func serveHome(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
    
    /*trueHost := r.Header().Get("X-Forwarded-For")
    var keys []string
    for k := range r.Header {
        keys = append(keys, k)
    }
    sort.Strings(keys)
    for _, k := range keys {
           log.Println(k, ":", r.Header[k])
    } */
    var v = struct {
		Host    string
		Data    string
	}{
		r.Host,
		currentMessage,
	}
	homeTempl.Execute(w, &v)
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil { return }

    go writer(ws)    
    
}

func writer(ws *websocket.Conn) {
    defer func() {
		ws.Close()
	}()
    for {
        log.Println("Waiting for chanel..")
        msg := <- messageChan
        if err := ws.WriteMessage(websocket.TextMessage, msg); err != nil {
            return
        }
    }
    
}

func main() {
    subscribed := make(chan bool)
    go receiveMessage(subscribed)
	// wait until we know the receiver has subscribed
	<-subscribed
    
    http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", serveWs)
    log.Fatal(http.ListenAndServe(":80", nil))
}

const homeHTML = `<!DOCTYPE html>
<html lang="en">
    <head>
        <title>WebSocket Messager</title>
    </head>
    <body>
        <pre id="Message">{{.Data}}</pre>
        <script type="text/javascript">
            (function() {
                var data = document.getElementById("Message");
                var conn = new WebSocket("ws://{{.Host}}/ws");
                conn.onclose = function(evt) {
                    data.textContent = 'Connection closed';
                }
                conn.onmessage = function(evt) {
                    console.log('New message received');
                    data.textContent = evt.data;
                }
            })();
        </script>
    </body>
</html>
`