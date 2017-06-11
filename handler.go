package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

const SHOWMAXROW = 40

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// ping is used to test whether the http service starts normally
// This is design commonly for the http service
func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("PONG"))
}

// getMainInterface get main interface
func getMainInterface(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "Go Server")
	w.Header().Set("Content-type", "text/html")
	html, err := Asset("static/index.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Write(html)
}

// getAsset get static resource including css and js
func getAsset(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "Go Server")
	path := strings.TrimLeft(r.URL.RequestURI(), "/")
	if strings.HasSuffix(path, ".css") {
		w.Header().Set("Content-type", "text/css")
	}
	if strings.HasSuffix(path, ".js") {
		w.Header().Set("Content-type", "application/x - javascript")
	}
	cssJs, err := Asset(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Write(cssJs)
}

// getServerList from .toml file
func getServerList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Server", "Go Server")
	w.Header().Set("Content-type", "application/json")
	isIncludeStatus := r.URL.Query().Get(`status`)
	switch isIncludeStatus {
	// case `include`:
	// 	io.WriteString(w, getServers(isIncludeStatus))
	// case `exclude`:
	// 	io.WriteString(w, getServers(isIncludeStatus))
	case `serverinfo`:
		io.WriteString(w, getServerInfos())
	case `remove`:
		io.WriteString(w, removeServerInfo(r.URL.Query().Get(`name`)))
	}
}

// addServer write into .toml file
// validate date and write into toml file
func addServer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Server", "Go Server")
	w.Header().Set("Content-type", "text")
	name := r.URL.Query().Get(`name`)
	host := r.URL.Query().Get(`host`)
	port := r.URL.Query().Get(`port`)
	auth := r.URL.Query().Get(`auth`)
	io.WriteString(w, writeServerToml(name, host, port, auth))
}

// getServerKeyNameHtml get html
func getServerKeyNameHtml(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Server", "Go Server")
	w.Header().Set("Content-type", "text/html")
	html, err := Asset("static/redb.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Write(html)
}

// getInfoByServerAndDb connect to redis-server accroding db serial and server address
func getInfoByServerAndDb(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Server", "Go Server")
	server := r.URL.Query().Get(`server`)
	db_serial := r.URL.Query().Get(`db`)
	var showmore = r.URL.Query().Get(`showmore`)
	if showmore == `` {
		showmore = `0`
	}
	if server == `` || db_serial == `` {
		io.WriteString(w, `{}`)
		return
	}
	io.WriteString(w, getTypeNameAndKeyByDb(server, db_serial, showmore))
}

// getInfoBySearchKey search redis key according search key
func getInfoBySearchKey(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Server", "Go Server")
	server := r.URL.Query().Get(`server`)
	db_serial := r.URL.Query().Get(`db`)
	search_key := r.URL.Query().Get(`search`)
	var showmore = r.URL.Query().Get(`showmore`)
	if showmore == `` {
		showmore = `0`
	}
	if server == `` || db_serial == `` {
		io.WriteString(w, `{}`)
		return
	}
	io.WriteString(w, getTypeNameAndKeyBySearchKey(server, db_serial, search_key, showmore))
}

// getInfoByTypeNameorKey get field value from type and name
func getInfoByTypeNameorKey(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Server", "Go Server")
	server := r.URL.Query().Get(`server`)
	db_serial := r.URL.Query().Get(`db`)
	style := r.URL.Query().Get(`style`)
	name := r.URL.Query().Get(`name`)
	key_name := r.URL.Query().Get(`key_name`)
	var showmore = r.URL.Query().Get(`showmore`)
	if server == `` || db_serial == `` {
		io.WriteString(w, `{}`)
		return
	}
	if showmore == `` {
		showmore = `0`
	}
	if key_name == `` {
		io.WriteString(w, getKeyContentByTypeNameorKey(server, db_serial, style, name, showmore))
		return
	}
	io.WriteString(w, getContentByTypeNameAndKey(server, db_serial, style, name, key_name))
}

// addKeysByTypeAndName add key by type and name
func addFieldsByTypeAndName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Server", "Go Server")
	server := r.URL.Query().Get(`server`)
	db_serial := r.URL.Query().Get(`db`)
	style := r.URL.Query().Get(`style`)
	name := r.URL.Query().Get(`name`)
	field := r.URL.Query().Get(`field`)
	if server == `` || db_serial == `` {
		io.WriteString(w, `{}`)
		return
	}
	io.WriteString(w, addKeyOrField(server, db_serial, style, name, field))
}

// deleteTypeName key by type and name
func deleteTypeName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Server", "Go Server")
	server := r.URL.Query().Get(`server`)
	db_serial := r.URL.Query().Get(`db`)
	style := r.URL.Query().Get(`style`)
	name := r.URL.Query().Get(`name`)
	if server == `` || db_serial == `` {
		io.WriteString(w, `{}`)
		return
	}
	fmt.Println(style, name)
	io.WriteString(w, delete(server, db_serial, style, name))
}

// modify key name
func modifyTypeName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Server", "Go Server")
	server := r.URL.Query().Get(`server`)
	db_serial := r.URL.Query().Get(`db`)
	style := r.URL.Query().Get(`style`)
	oldname := r.URL.Query().Get(`oldname`)
	newname := r.URL.Query().Get(`newname`)
	if server == `` || db_serial == `` {
		io.WriteString(w, `{}`)
		return
	}
	fmt.Println(style, oldname, newname)
	io.WriteString(w, modify(server, db_serial, style, oldname, newname))
}

// saveChangeContent change content
func saveChangeContent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Server", "Go Server")
	server := r.URL.Query().Get(`server`)
	db_serial := r.URL.Query().Get(`db`)
	style := r.URL.Query().Get(`style`)
	name := r.URL.Query().Get(`name`)
	field := r.URL.Query().Get(`field`)
	index := r.URL.Query().Get(`index`)
	content := r.URL.Query().Get(`content`)
	if server == `` || db_serial == `` {
		io.WriteString(w, `{}`)
		return
	}
	io.WriteString(w, changeContent(server, db_serial, style, name, index, field, content))
}

// execInstruction responsible for the front end of the interface to interact with the terminal
func execInstruction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Server", "Go Server")
	server := r.URL.Query().Get(`server`)
	command := r.URL.Query().Get(`command`)
	io.WriteString(w, executeCommand(server, command))
}

// login validate
func login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Server", "Go Server")
	admin := r.URL.Query().Get(`admin`)
	password := r.URL.Query().Get(`password`)
	io.WriteString(w, loginCheck(admin, password))
}

// get config from server
func getConfigInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Server", "Go Server")
	server := r.URL.Query().Get(`server`)
	operation := r.URL.Query().Get(`operation`)
	switch operation {
	case `showconfig`:
		io.WriteString(w, getConfigInfoFromServer(server))
		break
	case `saveconfig`:
		body, _ := ioutil.ReadAll(r.Body)
		io.WriteString(w, saveConfigInfo(server, string(body)))
		break
	}

}

// redisMonitorRealTime websocket connect for redis command of monitor
func redisMonitorRealTime(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(`Upgrade fail ... `, err)
		return
	}
	// 获取数据
	go readMessage(conn)
	// 推送数据
	go realMonitorTime()
	for {
		select {
		case monitor := <-MonitorMessage:
			if err = conn.WriteMessage(websocket.TextMessage, []byte(monitor.Message)); err != nil {
				fmt.Println(`WriteMessage fail ... `, err)
				return
			}
		case <-ctrlcMessage:
			return
		}
	}
}

// middleware to protect private pages
func validate(page http.HandlerFunc) http.HandlerFunc {
	if !isValidate() {
		return func(w http.ResponseWriter, r *http.Request) {
			page(w, r)
			return
		}
	}
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie("Auth")
		if err != nil {
			w.Header().Set("Auth", `Restricted`)
		}
		page(w, r)
		return
	}
}
