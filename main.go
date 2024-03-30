package main

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/mole828/onebot-extension/onebot"
	"github.com/samber/lo"
)

var target = "ws://10.0.0.42:7780" // 目标服务器地址

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handler(w http.ResponseWriter, r *http.Request) {
	clientConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}
	defer clientConn.Close()

	url, _ := url.Parse(target)
	targetConn, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}
	defer targetConn.Close()

	go func() {
		for {
			messageType, message, err := clientConn.ReadMessage()
			if err != nil {
				return
			}
			if err := targetConn.WriteMessage(messageType, message); err != nil {
				return
			}
		}
	}()

	for {
		messageType, message, err := targetConn.ReadMessage()
		if err != nil {
			return
		}
		var ret = new(onebot.Ret)
		if err := json.Unmarshal(message, ret); err == nil {
			if ret.Echo == "getFriendList" {
				var getFriendListRet = new(onebot.GetFriendListRet)
				json.Unmarshal(message, getFriendListRet)
				data := lo.UniqBy[onebot.User, int](getFriendListRet.Data, func(user onebot.User) int {
					return user.UserId
				})
				getFriendListRet.Data = data
				message, _ = json.Marshal(getFriendListRet)
			}
		}
		if err := clientConn.WriteMessage(messageType, message); err != nil {
			return
		}
	}
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
