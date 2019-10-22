package main

import (
	"os"
	"log"
	"fmt"
	"net/http"
	"golang.org/x/net/websocket"
	"strings"
	//"encoding/json"
	"bytes"
	"runtime"
	"strconv"
	//"reflect"
	"WSMatcher/room"
)

func getRoomId(uri string) string {
	uri_head := "/wsmatcher/"
	b := strings.Index(uri, uri_head) + len(uri_head)
	e := strings.Index(uri, "?")

	if b < 0 || e < 0 || e < b {
		return ""
	}

	return string([]byte(uri)[b:e])
}


func goID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

func verifyConn(ws *websocket.Conn) {
	//RequestURI /wsmatcher/abcdefg?type=1
	remoteAddr := ws.Request().RemoteAddr
	ws.Request().ParseForm();
	roomType := ws.Request().Form.Get("type")
	roomId := getRoomId(ws.Request().RequestURI)

	log.Println("open: ", remoteAddr, goID(), ws.Request().RequestURI)
	client := room.Matcher(roomId, roomType, nil)

	closeMsg := ""
	defer func() {
		log.Println("close: ", remoteAddr, closeMsg, ws.Request().RequestURI)
		if client != nil { client.Close() }
		ws.Close()
	}()

	if client == nil {
		closeMsg = "client is nil"
		return
	}
	
	client.SetFun_WriteString(func(str string){	websocket.Message.Send(ws, str)	})
	client.SetFun_Disconnect(func(){ws.Close()})

	for {
		var str string
		if err := websocket.Message.Receive(ws, &str); err != nil {
			closeMsg = "client is disconnect"
			break;
		}
		client.ReadString(str)
	}
}

func main() {
	addr := ":3333"

	if len(os.Args) == 2 {
		addr = os.Args[1]
		log.Println("Listening:", addr)
	} else{
		log.Println("USEAGE: WSMatcher {listenAddr}")
		log.Println("Listening Default:", addr)
	}

	http.Handle("/wsmatcher/", websocket.Handler(verifyConn))
	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}