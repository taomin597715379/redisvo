package main

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/websocket"
	_ "github.com/taomin597715379/redisvo/splitecommand"
	"github.com/xuyu/goredis"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

// readMessage read from websocket request, request body include
// redis's monitor command and clientActiveClose chan is for notice
// goroutine realMonitorTime
func readMessage(conn *websocket.Conn) {
	for {
		_, server_info, err := conn.ReadMessage()
		if err != nil {
			clientActiveClose <- true
			return
		}
		if string(server_info) == `Ctrl+C` {
			ctrlcMessage <- string(server_info)
			receiveOver <- string(server_info)
			return
		} else {
			serverMessage <- string(server_info)
		}
	}
}

// realMonitorTime is for Redis the monitor data in real time to the front
func realMonitorTime() {
	var client *goredis.Redis
	var monitor *goredis.MonitorCommand
	var info_slice []string
	var server_info string
	server := <-serverMessage
	if strings.Contains(server, "_") {
		info_slice = strings.Split(server, "_")
		server_info = info_slice[0]
	} else {
		server_info = server
	}
	if _, ok := goRedisMap[server_info]; ok {
		client = goRedisMap[server_info]
	} else {
		client, _ = Dial(server_info, "")
		goRedisMap[server_info] = client
	}
	monitor, err := client.Monitor()
	if err != nil {
		return
	}
	t := time.NewTicker(50 * time.Millisecond)
	for {
		select {
		case <-receiveOver:
			return
		case <-clientActiveClose:
			return
		case <-t.C:
			h, err := monitor.Receive()
			if len(info_slice) == 0 {
				MonitorMessage <- MonitorInfo{Message: h, Err: err}
			}
			if len(info_slice) > 0 {
				pipe := strings.Split(info_slice[1], " ")
				if pipe[0] == `findstr` || pipe[0] == `grep` {
					up_stirng := strings.ToUpper(strings.Trim(pipe[1], `"`))
					low_stirng := strings.ToLower(strings.Trim(pipe[1], `"`))
					if strings.Contains(h, up_stirng) || strings.Contains(h, low_stirng) {
						MonitorMessage <- MonitorInfo{Message: h, Err: err}
					}
				} else {
					return
				}
			}
		}
	}
}

// saveConfigInfo save changed config info to redis-config.conf file
// according to redis config command
func saveConfigInfo(server, value string) string {
	var conf map[string]string
	rdsConn := getConnectFromPool(server)
	defer rdsConn.Close()
	err := json.Unmarshal([]byte(value), &conf)
	if err != nil {
		fmt.Println("error:", err)
		return `Error`
	}
	fmt.Println(conf)
	cache := cacheRedisConf[server]
	fmt.Println(cache)
	for k, v := range conf {
		if v != cache[k] {
			_, err := redis.String(rdsConn.Do("config", "set", k, v))
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
	}
	return `OK`
}

// getConfigInfo from server and Let us sort out the current redis.conf file
// can be obtained in the configuration items. now The page only shows the
// configuration in the redis.conf file.
func getConfigInfoFromServer(server string) string {
	var buf []byte
	var redisConf RedisConfig
	rdsConn := getConnectFromPool(server)
	defer rdsConn.Close()
	t := reflect.TypeOf(&redisConf).Elem()
	for i := 0; i < t.NumField(); i++ {
		name := t.Field(i).Tag.Get("json")
		redisConfMap[name] = ""
	}
	for k, _ := range redisConfMap {
		conf, err := redis.StringMap(rdsConn.Do("config", "get", k))
		if err != nil {
			fmt.Println(err)
			break
		}
		redisConfMap[k] = conf[k]
	}
	cacheRedisConf[server] = redisConfMap
	buf, _ = json.Marshal(redisConfMap)
	return string(buf)
}

// executeCommand analysis the ajax request of command
func executeCommand(server, command string) string {
	var str string
	rdsConn := getConnectFromPool(server)
	defer rdsConn.Close()
	commArg, err := splitecommand.AnalyzeCommand(command)
	if err != nil {
		return err.Error()
	}
	if len(commArg) == 2 && (strings.ToUpper(commArg[0]) == `DEL` ||
		strings.ToLower(commArg[0]) == `del`) &&
		commArg[1] != `*` && strings.Contains(commArg[1], "*") {
		r := batchDeletion(rdsConn, commArg[1])
		return strconv.FormatInt(r.(int64), 10)
	}
	switch len(commArg) {
	case 1:
		rdsConn.Send(commArg[0])
	case 2:
		rdsConn.Send(commArg[0], commArg[1])
	case 3:
		rdsConn.Send(commArg[0], commArg[1], commArg[2])
	case 4:
		rdsConn.Send(commArg[0], commArg[1], commArg[2], commArg[3])
	case 5:
		rdsConn.Send(commArg[0], commArg[1], commArg[2], commArg[3], commArg[4])
	default:
		return fmt.Sprintf("ERR wrong number of arguments for %s command", commArg[0])
	}
	rdsConn.Flush()
	r, err := rdsConn.Receive()
	if err != nil {
		return fmt.Sprintf("%s", err)
	}
	if r == nil {
		return fmt.Sprintf("%s", "")
	}
	switch t := r.(type) {
	case int64:
		return strconv.FormatInt(r.(int64), 10)
	case []uint8:
		return fmt.Sprintf("%s", t)
	case []interface{}:
		for _, v := range t {
			str += (fmt.Sprintf("%s", v) + "\n")
		}
		return strings.TrimRight(str, "\n")
	default:
		break
	}
	return fmt.Sprintf("%s", r)
}

// changeContent save content changed information
func changeContent(server, db, style, name, index, field, content string) string {
	dbno, _ := strconv.ParseInt(db, 10, 0)
	idx, _ := strconv.ParseInt(index, 10, 0)
	rdsConn := getConnectFromPool(server)
	defer rdsConn.Close()
	if _, err := redis.String(rdsConn.Do("select", dbno)); err != nil {
		fmt.Println(err)
		return "false"
	}
	if name == `` {
		return "false"
	}
	err := saveContent(rdsConn, style, name, field, content, idx)
	if err != nil {
		fmt.Println(err)
		return "false"
	}
	return "true"
}

// addKeyOrField save the modified redis key or field
func addKeyOrField(server, db, style, name, field string) string {
	dbno, _ := strconv.ParseInt(db, 10, 0)
	rdsConn := getConnectFromPool(server)
	defer rdsConn.Close()
	if _, err := redis.String(rdsConn.Do("select", dbno)); err != nil {
		fmt.Println(err)
		return "false"
	}
	if name == `` {
		return "false"
	}
	ok, err := setKeyOrFieldByStyle(rdsConn, style, name, field)
	if ok != `OK` {
		fmt.Println(err)
		return "false"
	}
	return "true"
}

// delete key or field
func delete(server, db, style, name string) string {
	dbno, _ := strconv.ParseInt(db, 10, 0)
	rdsConn := getConnectFromPool(server)
	defer rdsConn.Close()
	if _, err := redis.String(rdsConn.Do("select", dbno)); err != nil {
		fmt.Println(err)
		return "false"
	}
	if name == `` {
		return "false"
	}
	_, err := redis.Int64(rdsConn.Do("del", name))
	if err != nil {
		fmt.Println(err)
		return "false"
	}
	return "true"
}

// modify key name
func modify(server, db, style, oldname, newname string) string {
	dbno, _ := strconv.ParseInt(db, 10, 0)
	rdsConn := getConnectFromPool(server)
	defer rdsConn.Close()
	if _, err := redis.String(rdsConn.Do("select", dbno)); err != nil {
		fmt.Println(err)
		return "false"
	}
	if newname == `` {
		return "false"
	}
	_, err := redis.String(rdsConn.Do("rename", oldname, newname))
	if err != nil {
		fmt.Println(err)
		return "false"
	}
	return "true"
}

// getTypeNameAndKeyByDb the database to select when dialing a connection
func getTypeNameAndKeyByDb(server, db, showmore string) string {
	var buf []byte
	var rdsTypeName RedisTypeName
	var typeNames TypeNames
	var keyNames []KeyName
	var content string
	dbno, _ := strconv.ParseInt(db, 10, 0)
	more, _ := strconv.ParseInt(showmore, 10, 0)
	rdsConn := getConnectFromPool(server)
	defer rdsConn.Close()
	if _, err := redis.String(rdsConn.Do("select", dbno)); err != nil {
		fmt.Println(err)
		return `{"typename":[],"keysnameswithtype":{"keysname":[]},"content":""}`
	}
	redisReply, err := redis.Strings(rdsConn.Do("keys", "*"))
	if err != nil {
		fmt.Println(err)
		return `{"typename":[],"keysnameswithtype":{"keysname":[]},"content":""}`
	}
	if len(redisReply) <= 0 {
		return `{"typename":[],"keysnameswithtype":{"keysname":[]},"content":""}`
	}
	sumShowNumber := getInfo(rdsConn, "Keyspace", db)
	typeNames = moreTypeName(sumShowNumber, more, redisReply, rdsConn)
	rdsTypeName.TypeNames = typeNames
	if len(rdsTypeName.TypeNames) > 0 {
		tye := rdsTypeName.TypeNames[0].Type
		name := rdsTypeName.TypeNames[0].Name
		keyNames, content = getKeysByTypeName(rdsConn, tye, name)
		sumShowNumber := int64(len(keyNames))
		keyNames := moreFieldName(sumShowNumber, more, keyNames)
		rdsTypeName.KeysNamesWithType.KeysNames = keyNames
		rdsTypeName.KeysNamesWithType.SelfTypeName = rdsTypeName.TypeNames[0]
		rdsTypeName.KeysNamesWithType.KeysNames = keyNames
		rdsTypeName.Contents = content
		buf, _ = json.Marshal(rdsTypeName)
		return string(buf)
	}
	return `{"typename":[],"keysnameswithtype":{"keysname":[]},"content":""}`
}

// getTypeNameAndKeyBySearchKey search for redis key
func getTypeNameAndKeyBySearchKey(server, db, search_key, showmore string) string {
	var buf []byte
	var rdsTypeName RedisTypeName
	var typeNames TypeNames
	var keyNames []KeyName
	var content string
	dbno, _ := strconv.ParseInt(db, 10, 0)
	more, _ := strconv.ParseInt(showmore, 10, 0)
	rdsConn := getConnectFromPool(server)
	defer rdsConn.Close()
	if _, err := redis.String(rdsConn.Do("select", dbno)); err != nil {
		fmt.Println(err)
		return `{"typename":[],"keysnameswithtype":{"keysname":[]},"content":""}`
	}
	redisReply, err := redis.Strings(rdsConn.Do("keys", search_key+"*"))
	if err != nil {
		fmt.Println(err)
		return `{"typename":[],"keysnameswithtype":{"keysname":[]},"content":""}`
	}
	if len(redisReply) <= 0 {
		return `{"typename":[],"keysnameswithtype":{"keysname":[]},"content":""}`
	}
	sumShowNumber := int64(len(redisReply))
	typeNames = moreTypeName(sumShowNumber, more, redisReply, rdsConn)
	rdsTypeName.TypeNames = typeNames
	if len(rdsTypeName.TypeNames) > 0 {
		tye := rdsTypeName.TypeNames[0].Type
		name := rdsTypeName.TypeNames[0].Name
		keyNames, content = getKeysByTypeName(rdsConn, tye, name)
		sumShowNumber := int64(len(keyNames))
		keyNames := moreFieldName(sumShowNumber, more, keyNames)
		rdsTypeName.KeysNamesWithType.KeysNames = keyNames
		rdsTypeName.KeysNamesWithType.SelfTypeName = rdsTypeName.TypeNames[0]
		rdsTypeName.KeysNamesWithType.KeysNames = keyNames
		rdsTypeName.Contents = content
		buf, _ = json.Marshal(rdsTypeName)
		return string(buf)
	}
	return `{"typename":[],"keysnameswithtype":{"keysname":[]},"content":""}`
}

// getKeyContentByTypeNameorKey according redis one key to get field or value
func getKeyContentByTypeNameorKey(server, db, style, name, showmore string) string {
	var buf []byte
	var rdsKeyName RedisKeyName
	var keyNames []KeyName
	var content string
	dbno, _ := strconv.ParseInt(db, 10, 0)
	more, _ := strconv.ParseInt(showmore, 10, 0)
	rdsConn := getConnectFromPool(server)
	defer rdsConn.Close()
	if _, err := redis.String(rdsConn.Do("select", dbno)); err != nil {
		fmt.Println(err)
		return `{"typename":[],"keysnameswithtype":{"keysname":[]},"content":""}`
	}
	keyNames, content = getKeysByTypeName(rdsConn, style, name)

	sumShowNumber := int64(len(keyNames))
	keyNames = moreFieldName(sumShowNumber, more, keyNames)
	rdsKeyName.KeysNamesWithType.KeysNames = keyNames
	rdsKeyName.KeysNamesWithType.SelfTypeName = TypeName{Type: style, Name: name}
	rdsKeyName.Contents = content
	buf, _ = json.Marshal(rdsKeyName)
	return string(buf)
}

// getContentByTypeNameAndKey according redis one key or one field to get value
func getContentByTypeNameAndKey(server, db, style, name, key_name string) string {
	var buf []byte
	var rdsKeyName RedisKeyName
	var keyNames []KeyName
	var content string
	dbno, _ := strconv.ParseInt(db, 10, 0)
	rdsConn := getConnectFromPool(server)
	defer rdsConn.Close()
	if _, err := redis.String(rdsConn.Do("select", dbno)); err != nil {
		fmt.Println(err)
		return `{"typename":[],"keysnameswithtype":{"keysname":[]},"content":""}`
	}
	content = getContentByTypeNameAnd(rdsConn, style, name, key_name)
	rdsKeyName.KeysNamesWithType.KeysNames = keyNames
	rdsKeyName.KeysNamesWithType.SelfTypeName = TypeName{Type: style, Name: name}
	rdsKeyName.Contents = content
	buf, _ = json.Marshal(rdsKeyName)
	return string(buf)
}

func getInfo(c redis.Conn, info_key, dbno string) int64 {
	var cnt int64 = 0
	value, err := redis.String(c.Do("info", info_key))
	if err != nil {
		return int64(0)
	}
	for _, sub_string := range strings.Split(value, ",") {
		if sub_string != `` && strings.Contains(sub_string, "db"+dbno+":keys=") {
			cnt, _ = strconv.ParseInt(strings.Split(sub_string, "=")[1], 10, 32)
			break
		}
	}
	return cnt
}

func moreTypeName(sumShowNumber, more int64, redisReply []string, rdsConn redis.Conn) (typeNames TypeNames) {
	var i, moreFlag, natureShowNumber int64
	if sumShowNumber <= SHOWMAXROW {
		for _, name := range redisReply {
			tp, _ := redis.String(rdsConn.Do("type", name))
			typeNames = append(typeNames, TypeName{Type: tp, Name: name})
		}
	} else {
		var needShowNumber = SHOWMAXROW * (more + 1)
		if needShowNumber < sumShowNumber {
			natureShowNumber = needShowNumber
			moreFlag = 1
		} else {
			moreFlag = 0
			natureShowNumber = sumShowNumber
		}
		for _, name := range redisReply {
			i++
			tp, _ := redis.String(rdsConn.Do("type", name))
			typeNames = append(typeNames, TypeName{Type: tp, Name: name})
			if i >= natureShowNumber {
				break
			}
		}
	}
	sort.Sort(typeNames)
	if moreFlag == 1 {
		typeNames = append(typeNames, TypeName{Type: "", Name: "More"})
	}
	return typeNames
}

func moreFieldName(sumShowNumber, more int64, keyNames []KeyName) []KeyName {
	var moreFlag, natureShowNumber int64
	if sumShowNumber <= SHOWMAXROW {
		return keyNames
	} else {
		var needShowNumber = SHOWMAXROW * (more + 1)
		if needShowNumber < sumShowNumber {
			natureShowNumber = needShowNumber
			moreFlag = 1
		} else {
			moreFlag = 0
			natureShowNumber = sumShowNumber
		}
		keyNames = keyNames[:natureShowNumber]
		if moreFlag == 1 {
			keyNames = append(keyNames, KeyName{Name: "More"})
		}
		return keyNames
	}
}

// batchDeletion suuport redis batch delete keys
func batchDeletion(rdsConn redis.Conn, commArg string) interface{} {
	var count int64
	redisReply, err := redis.Strings(rdsConn.Do("keys", commArg))
	if err != nil {
		return int64(0)
	}
	if len(redisReply) <= 0 {
		return int64(0)
	}
	for _, name := range redisReply {
		count++
		_, err := redis.Int64(rdsConn.Do("del", name))
		if err != nil {
			fmt.Println(err)
		}
	}
	return count
}
