package conn

import (
	"chatroom/message"
	"chatroom/libs"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 解决跨域问题
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var connMap sync.Map

const JOIN = "上线啦"
const LEAVE = "下线了"
const TEXT_TYPE = "text_type"
const STATUS_TYPE = "status_type"

type connInfo struct {
	Uid      string
	Gravatar string
	Conn     *websocket.Conn
	mutex    sync.Mutex
}

var data []byte

func Connection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	uid  := r.FormValue("uid")
	gravatar := libs.UrlSize(uid, 32)

	msg := &message.Message{
		Uid:uid,
		Content: uid + JOIN,
		Gravatar:gravatar,
		Time:time.Now().Format("2006-01-02 15:04:05"),
		Type:STATUS_TYPE,
	}
	data, _ = proto.Marshal(msg)

	Broadcast(data)

	connInfo := &connInfo{Uid: uid, Gravatar: gravatar, Conn: conn}
	connMap.Store(uid, connInfo)

	Receive(connInfo)
}

func Receive(connInfo *connInfo) {
	for {
		_, p, err := connInfo.Conn.ReadMessage()
		if err != nil {
			connMap.Delete(connInfo.Uid)
			
			msg := &message.Message{
				Uid:connInfo.Uid,
				Content: connInfo.Uid + LEAVE,
				Gravatar:connInfo.Gravatar,
				Time:time.Now().Format("2006-01-02 15:04:05"),
				Type:STATUS_TYPE,
			}
			data, _ = proto.Marshal(msg)
			Broadcast(data)

			return
		}

		msg := &message.Message{
			Uid:connInfo.Uid,
			Content: string(p),
			Gravatar:connInfo.Gravatar,
			Time:time.Now().Format("2006-01-02 15:04:05"),
			Type:TEXT_TYPE,
		}
		data, _ = proto.Marshal(msg)

		go Broadcast(data)
	}
}

func Broadcast(data []byte) {
	connMap.Range(func(key, value interface{}) bool {
		tmpval, _ := value.(*connInfo)

		tmpval.mutex.Lock()
		defer tmpval.mutex.Unlock()

		if err := tmpval.Conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
			return true
		}
		return true
	})

}