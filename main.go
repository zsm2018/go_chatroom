package main

import (
	"chatroom/conn"
	"html/template"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./template/index.tpl")
	_ = t.Execute(w, map[string]string{})
}

func Join(w http.ResponseWriter, r *http.Request) {
	uid  := r.FormValue("uid")
	if uid == "" {
		t, _ := template.ParseFiles("./template/index.tpl")
		_ = t.Execute(w, map[string]string{})
		return
	}

	t, _ := template.ParseFiles("./template/room.tpl")
	_ = t.Execute(w, map[string]string{"uid": uid})
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))  // 启动静态文件服务

	http.HandleFunc("/chat", conn.Connection)
	http.HandleFunc("/join", Join)
	http.HandleFunc("/", Index)

	_ = http.ListenAndServe(":2222", nil)
}