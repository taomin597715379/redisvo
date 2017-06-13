package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/garyburd/redigo/redis"
	"github.com/xuyu/goredis"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

var serverPoolMap = make(map[string]*redis.Pool, 0)
var goRedisMap = make(map[string]*goredis.Redis, 0)
var MonitorMessage = make(chan MonitorInfo, 0)
var serverMessage = make(chan string, 0)
var ctrlcMessage = make(chan string, 0)
var receiveOver = make(chan string, 0)
var clientActiveClose = make(chan bool, 0)
var redisConfMap = make(map[string]string, 0)
var cacheRedisConf = make(map[string]map[string]string, 0)

const (
	MAXFIELDNUMBER = 100000
	SELF_CONF_FILE = `.` + string(os.PathSeparator) + `redisvo.toml`
)

func (c Int64Slice) Len() int {
	return len(c)
}
func (c Int64Slice) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c Int64Slice) Less(i, j int) bool {
	return c[i] < c[j]
}
func (u TypeNames) Len() int {
	return len(u)
}
func (u TypeNames) Less(i, j int) bool {
	return u[i].Name < u[j].Name
}
func (u TypeNames) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}

// getConnectFromPool Each redis-server: ip, port and a redis-connect mapping
func getConnectFromPool(server string) redis.Conn {
	if _, ok := serverPoolMap[server]; !ok {
		rdsPool := newRdsPool(server, "")
		serverPoolMap[server] = rdsPool
	}
	return serverPoolMap[server].Get()
}

// newRdsPool connection estabish
func newRdsPool(server, auth string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     100,
		MaxActive:   30,
		IdleTimeout: 60 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if auth == "" {
				return c, err
			}
			if _, err := c.Do("AUTH", auth); err != nil {
				c.Close()
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

// Dial connection estabish
func Dial(server, auth string) (*goredis.Redis, error) {
	connect, err := goredis.Dial(&goredis.DialConfig{
		Network: "tcp",
		Address: server,
		Timeout: 10 * time.Second,
		MaxIdle: 10,
	})
	return connect, err
}

// saveContent
func saveContent(rdsConn redis.Conn, style, name, field, content string, idx int64) (err error) {
	switch style {
	case `string`:
		_, err = redis.String(rdsConn.Do("set", name, content))
		return err
	case `hash`:
		_, err = redis.String(rdsConn.Do("hset", name, field, content))
		return err
	case `list`:
		_, err = redis.String(rdsConn.Do("lset", name, idx, content))
		return err
	case `set`:
		_, err = redis.String(rdsConn.Do("srem", name, field))
		_, err = redis.String(rdsConn.Do("sadd", name, content))
		return err
	case `zset`:
		_, err = redis.String(rdsConn.Do("zrem", name, field))
		_, err = redis.String(rdsConn.Do("zadd", name, idx, content))
		return err
	default:
		return nil
	}
	return nil
}

// setKeyOrFieldByStyle
func setKeyOrFieldByStyle(rdsConn redis.Conn, style, name, field string) (content string, err error) {
	var value string
	switch style {
	case `string`:
		content, err = redis.String(rdsConn.Do("set", name, field))
		return content, err
	case `hash`:
		if field == `` {
			field = `New Key`
		}
		content, err = redis.String(rdsConn.Do("hset", name, field, "New Member"))
		return content, err
	case `list`:
		if field == `` {
			field = `New Item`
			content, err = redis.String(rdsConn.Do("lpush", name, field))
			return content, err
		}
		splitField := strings.Split(field, `_`)
		value = splitField[1]
		if splitField[0] == `` {
			value = `New Item`
		}
		if splitField[0] == `head` {
			content, err = redis.String(rdsConn.Do("lpush", name, value))
		}
		if splitField[0] == `tail` {
			content, err = redis.String(rdsConn.Do("rpush", name, value))
		}
		return content, err
	case `set`:
		if field == `` {
			field = `New Member`
		}
		content, err = redis.String(rdsConn.Do("sadd", name, field))
		return content, err
	case `zset`:
		if field == `` {
			field = `New Zmember`
			content, err = redis.String(rdsConn.Do("zadd", name, 0, field))
			return ``, nil
		}
		splitField := strings.Split(field, `_`)
		value = splitField[1]
		if splitField[1] == `` {
			value = `New Item`
		}
		score, _ := strconv.ParseInt(splitField[0], 10, 0)
		content, err = redis.String(rdsConn.Do("zadd", name, score, value))
		return content, err
	default:
		return ``, nil
	}
	return ``, nil
}

// getContentByTypeNameAnd according type, name, key_name to get content
func getContentByTypeNameAnd(rdsConn redis.Conn, typ, name, key_name string) (content string) {
	switch typ {
	case `string`:
		content, _ = redis.String(rdsConn.Do("get", name))
		return content
	case `hash`:
		content, _ = redis.String(rdsConn.Do("hget", name, key_name))
		return content
	case `list`:
		content = key_name
		return content
	case `set`:
		content = key_name
		return content
	case `zset`:
		content = key_name
		return content
	default:
		break
	}
	return
}

// getKeysByTypeName according type, name to key and content
func getKeysByTypeName(rdsConn redis.Conn, typ, name string) (keyNames []KeyName, content string) {
	var i, l int
	switch typ {
	case `string`:
		keyNames = []KeyName{}
		content, _ = redis.String(rdsConn.Do("get", name))
		return keyNames, content
	case `hash`:
		fields, _ := redis.Strings(rdsConn.Do("hkeys", name))
		if len(fields) > 0 {
			for k, field := range fields {
				if i < MAXFIELDNUMBER {
					keyNames = append(keyNames, KeyName{Name: field, Score: k})
					i++
				}
			}
		}
		if len(keyNames) > 0 {
			content, _ = redis.String(rdsConn.Do("hget", name, keyNames[0].Name))
		}
		return keyNames, content
	case `list`:
		fields, _ := redis.Strings(rdsConn.Do("lrange", name, 0, -1))
		if len(fields) > 0 {
			for k, field := range fields {
				keyNames = append(keyNames, KeyName{Name: field, Score: k})
			}
		}
		if len(keyNames) > 0 {
			content = keyNames[0].Name
		}
		return keyNames, content
	case `set`:
		fields, _ := redis.Strings(rdsConn.Do("smembers", name))
		if len(fields) > 0 {
			for k, field := range fields {
				keyNames = append(keyNames, KeyName{Name: field, Score: k})
			}
		}
		if len(keyNames) > 0 {
			content = keyNames[0].Name
		}
		return keyNames, content
	case `zset`:
		fields, _ := redis.Strings(rdsConn.Do("zrange", name, 0, -1, "WITHSCORES"))
		if len(fields) > 0 {
			for i := 0; i < len(fields); i = i + 2 {
				j, _ := strconv.Atoi(fields[i+1])
				keyNames = append(keyNames, KeyName{Name: fields[i], Index: l, Score: j})
				l++
			}
		}
		if len(keyNames) > 0 {
			content = keyNames[0].Name
		}
	default:
		break
	}
	return
}

// getServers from toml file server list according parameter to get status
// func getServers(isIncludeStatus string) string {
// 	var conf ConfigInfo
// 	var includeStatus []ServerStatus
// 	var s ServerStatus
// 	var buf []byte
// 	_, err := toml.DecodeFile(SELF_CONF_FILE, &conf)
// 	if err != nil {
// 		fmt.Println(err)
// 		return ""
// 	}
// 	if isIncludeStatus == `exclude` {
// 		for _, serverInfo := range conf.ServerInfo {
// 			s = ServerStatus{ServerAddr: serverInfo.Host + ":" + serverInfo.Port, Status: 1}
// 			includeStatus = append(includeStatus, s)
// 		}
// 		buf, _ = json.Marshal(includeStatus)
// 	}
// 	if isIncludeStatus == `include` {
// 		timeout := time.Duration(1) * time.Second
// 		for _, serverInfo := range conf.ServerInfo {
// 			if c, err := redis.DialTimeout("tcp", serverInfo.Host+":"+serverInfo.Port, timeout, timeout, timeout); err == nil {
// 				s = ServerStatus{ServerAddr: serverInfo.Host + ":" + serverInfo.Port, Status: 1}
// 				c.Close()
// 			} else {
// 				s = ServerStatus{ServerAddr: serverInfo.Host + ":" + serverInfo.Port, Status: 0}
// 			}
// 			includeStatus = append(includeStatus, s)
// 		}
// 		buf, _ = json.Marshal(includeStatus)
// 	}
// 	return string(buf)
// }

// getServerInfos get server info for example redis-version, clients and so on
// others server info will be sort by ip
func getServerInfos() string {
	var buf []byte
	var conf ConfigInfo
	var serverInfos []ServerExtInfo
	var serverOnline = make(map[int64]ServerExtInfo, 0)
	var serverNoOnline = make(map[int64]ServerExtInfo, 0)
	var s ServerExtInfo
	_, err := toml.DecodeFile(SELF_CONF_FILE, &conf)
	if err != nil {
		return ""
	}
	for _, serverInfo := range conf.ServerInfo {
		timeout := time.Duration(100) * time.Millisecond
		if c, err := redis.DialTimeout("tcp", serverInfo.Host+":"+serverInfo.Port, timeout, timeout, timeout); err == nil {
			version, memory, clients, commands, count := getInfoByField(c)
			s = ServerExtInfo{ServerAddr: serverInfo.Host + ":" + serverInfo.Port,
				UserMemory:   memory,
				ClientOnline: clients,
				ExeCommand:   commands,
				RedisVer:     version,
				KeyNumber:    count}
			c.Close()
			serverOnline[ipToInteger(serverInfo.Host)+portToInteger(serverInfo.Port)] = s
		} else {
			s = ServerExtInfo{ServerAddr: serverInfo.Host + ":" + serverInfo.Port,
				UserMemory:   "-",
				ClientOnline: "-",
				ExeCommand:   "-",
				RedisVer:     "-",
				KeyNumber:    "-"}
			serverNoOnline[ipToInteger(serverInfo.Host)+portToInteger(serverInfo.Port)] = s
		}
	}
	tmpInfo := sortServerInfo(serverOnline)
	serverInfos = append(serverInfos, tmpInfo...)
	tmpInfo = sortServerInfo(serverNoOnline)
	serverInfos = append(serverInfos, tmpInfo...)
	buf, _ = json.Marshal(serverInfos)
	return string(buf)
}

// removeServerInfo delete server info from toml file
func removeServerInfo(name string) string {
	var conf ConfigInfo
	var sBuffer bytes.Buffer
	_, err := toml.DecodeFile(SELF_CONF_FILE, &conf)
	if err != nil {
		return "Error"
	}
	for i, serverInfo := range conf.ServerInfo {
		if name == serverInfo.Host+":"+serverInfo.Port {
			conf.ServerInfo = append(conf.ServerInfo[:i], conf.ServerInfo[i+1:]...)
			break
		}
	}
	err = toml.NewEncoder(&sBuffer).Encode(conf)
	if err != nil {
		fmt.Println(err)
		return "Error"
	}
	err = ioutil.WriteFile(SELF_CONF_FILE, []byte(sBuffer.String()), 0644)
	if err != nil {
		fmt.Println(err)
		return "Error"
	}
	return "OK"
}

func writeServerToml(name, host, port, auth string) string {
	var conf ConfigInfo
	var flag int = 0
	var sBuffer bytes.Buffer
	if _, err := os.Stat(SELF_CONF_FILE); os.IsNotExist(err) {
		f, err := os.Create(SELF_CONF_FILE)
		if err != nil {
			return `false`
		}
		defer f.Close()
	}
	_, err := toml.DecodeFile(SELF_CONF_FILE, &conf)
	if err != nil {
		return `false`
	}
	for _, serverInfo := range conf.ServerInfo {
		if host == serverInfo.Host && port == serverInfo.Port {
			flag = 1
			break
		}
	}
	if flag == 0 {
		c := Info{Name: name, Host: host, Port: port, Auth: auth}
		conf.ServerInfo = append(conf.ServerInfo, c)
	}
	if conf.ServerAddress == `` {
		conf.ServerAddress = `127.0.0.1:7000`
	}
	if conf.AuthInfo == (Auth{}) {
		conf.AuthInfo = Auth{Admin: "", Password: "", Enable: 0}
	}
	err = toml.NewEncoder(&sBuffer).Encode(conf)
	if err != nil {
		fmt.Println(err)
		return `false`
	}
	err = ioutil.WriteFile(SELF_CONF_FILE, []byte(sBuffer.String()), 0644)
	if err != nil {
		fmt.Println(err)
		return `false`
	}
	return `OK`
}

func getInfoByField(c redis.Conn) (string, string, string, string, string) {
	var count int64
	var field = []string{"server", "memory", "clients", "stats"}
	var fieldRet = make(map[string]interface{}, 0)
	var fieldMap = map[string]string{"server": "redis_version", "memory": "used_memory", "clients": "connected_clients", "stats": "total_commands_processed"}
	for _, v := range field {
		value, err := redis.String(c.Do("info", v))
		if err != nil {
			return "-", "-", "-", "-", "-"
		}
		for _, sub_string := range strings.Split(value, "\n") {
			if strings.Contains(sub_string, fieldMap[v]) {
				fieldRet[v] = strings.Split(sub_string, ":")[1]
				break
			}
		}
	}
	for i := 0; i <= 15; i++ {
		count += getInfo(c, "Keyspace", strconv.Itoa(i))
	}
	return fieldRet[field[0]].(string), fieldRet[field[1]].(string),
		fieldRet[field[2]].(string), fieldRet[field[3]].(string), strconv.Itoa(int(count))
}

// ipToInteger ip convert to int
func ipToInteger(ip string) int64 {
	var result int64 = 0
	var ipArr []string = strings.Split(ip, ".")
	for k, v := range ipArr {
		j, _ := strconv.ParseInt(v, 10, 0)
		j = j << uint16((3-k)*8)
		result |= j
	}
	return result
}

// portToInteger port convert to int
func portToInteger(port string) int64 {
	j, _ := strconv.ParseInt(port, 10, 0)
	return j
}

// sortServerInfo
func sortServerInfo(info map[int64]ServerExtInfo) []ServerExtInfo {
	var ipInt64 Int64Slice
	var serverInfos []ServerExtInfo
	for k, _ := range info {
		ipInt64 = append(ipInt64, k)
	}
	sort.Sort(ipInt64)
	for _, v := range ipInt64 {
		serverInfos = append(serverInfos, info[v])
	}
	return serverInfos
}

// get serverAddress from config
func getServerAddress() string {
	var conf ConfigInfo
	if _, err := os.Stat(SELF_CONF_FILE); os.IsNotExist(err) {
		return `127.0.0.1:7000`
	}
	_, err := toml.DecodeFile(SELF_CONF_FILE, &conf)
	if err != nil {
		return `127.0.0.1:7000`
	}
	return conf.ServerAddress
}

// get serverAddress from config
func loginCheck(admin, password string) string {
	var conf ConfigInfo
	if _, err := os.Stat(SELF_CONF_FILE); os.IsNotExist(err) {
		return `false`
	}
	_, err := toml.DecodeFile(SELF_CONF_FILE, &conf)
	if err != nil {
		return `false`
	}
	if conf.AuthInfo.Admin == admin &&
		conf.AuthInfo.Password == password {
		return `true`
	}
	return `false`
}

// get validate info from config file
func isValidate() bool {
	var conf ConfigInfo
	if _, err := os.Stat(SELF_CONF_FILE); os.IsNotExist(err) {
		return false
	}
	_, err := toml.DecodeFile(SELF_CONF_FILE, &conf)
	if err != nil {
		return false
	}
	if conf.AuthInfo.Enable == 1 {
		return true
	}
	return false
}
