package main

type RedisConfig struct {
	Port           string `json:"port"`
	Bind           string `json:"bind"`
	Unixsocket     string `json:"unixsocket"`
	Unixsocketperm string `json:"unixsocketperm"`
	Daemonize      string `json:"daemonize"`
	Pidfile        string `json:"pidfile"`
	TcpBacklog     string `json:"tcp-backlog"`
	TcpKeepalive   string `json:"tcp-keepalive"`
	Timeout        string `json:"timeout"`
	Databases      string `json:"databases"`

	Loglevel string `json:"loglevel"`
	Logfile  string `json:"logfile"`

	Dbfilename              string `json:"dbfilename"`
	Dir                     string `json:"dir"`
	Save                    string `json:"save"`
	StopWritesOnBgSaveError string `json:"stop-writes-on-bgsave-error"`
	Rdbcompression          string `json:"rdbcompression"`
	Rdbchecksum             string `json:"rdbchecksum"`

	Slaveof               string `json:"slaveof"`
	Masterauth            string `json:"masterauth"`
	SlaveServeStaleData   string `json:"slave-serve-stale-data"`
	SlaveReadOnly         string `json:"slave-read-only"`
	ReplDiskLessSync      string `json:"repl-diskless-sync"`
	ReplDisklessSyncDelay string `json:"repl-diskless-sync-delay"`
	ReplPingSlavePeriod   string `json:"repl-ping-slave-period"`
	ReplTimeout           string `json:"repl-timeout"`
	ReplDisableTcpNodelay string `json:"repl-disable-tcp-nodelay"`
	ReplBacklogSize       string `json:"repl-backlog-size"`
	ReplBacklogTtl        string `json:"repl-backlog-ttl"`
	SlavePriority         string `json:"slave-priority"`
	MinSlavesToWrite      string `json:"min-slaves-to-write"`
	MinSlavesMaxLag       string `json:"min-slaves-max-lag"`

	Requirepass string `json:"requirepass"`

	MaxClients       string `json:"maxclients"`
	Maxmemory        string `json:"maxmemory"`
	MaxmemoryPolicy  string `json:"maxmemory-policy"`
	MaxmemorySamples string `json:"maxmemory-samples"`

	AppendonlySamples        string `json:"appendonly-samples"`
	AppendfsyncSamples       string `json:"appendfsync-samples"`
	NoAppendFsyncOnRewrite   string `json:"no-appendfsync-on-rewrite"`
	AutoAofRewritePercentage string `json:"auto-aof-rewrite-percentage"`
	AutoAofRewriteMinSize    string `json:"auto-aof-rewrite-min-size"`
	AofLoadTruncated         string `json:"aof-load-truncated"`

	LuaTimeLimit string `json:"lua-time-limit"`

	ClusterNodeTimeout         string `json:"cluster-node-timeout"`
	ClusterSlaveValidityFactor string `json:"cluster-slave-validity-factor"`
	ClusterMigrationBarrier    string `json:"cluster-migration-barrier"`
	ClusterRequireFullCoverage string `json:"cluster-require-full-coverage"`

	SlowlogLogSlowerThan string `json:"slowlog-log-slower-than"`
	SlowlogMaxLen        string `json:"slowlog-max-len"`

	LatencyMonitorThreshold string `json:"latency-monitor-threshold"`

	NotifyKeyspaceEvents string `json:"notify-keyspace-events"`

	HashMaxZiplistEntries      string `json:"hash-max-ziplist-entries"`
	HashMaxZiplistValue        string `json:"hash-max-ziplist-value"`
	SetMaxIntsetEntries        string `json:"set-max-intset-entries"`
	ZsetMaxZiplistEntries      string `json:"zset-max-ziplist-entries"`
	ZsetMaxZiplistValue        string `json:"zset-max-ziplist-value"`
	HllSparseMaxBytes          string `json:"hll-sparse-max-bytes"`
	Activerehashing            string `json:"activerehashing"`
	ClientOutputBufferLimit    string `json:"client-output-buffer-limit"`
	Hz                         string `json:"hz"`
	AofRewriteIncrementalFsync string `json:"aof-rewrite-incremental-fsync"`
}

type ConfigInfo struct {
	ServerAddress string `toml:"server_address"`
	AuthInfo      Auth   `toml:"auth_info"`
	ServerInfo    []Info `toml:"server_info"`
}

type Auth struct {
	Admin    string `toml:"admin"`
	Password string `toml:"password"`
	Enable   int    `toml:"enable"`
}

type Info struct {
	Name string `toml:"name"`
	Host string `toml:"host"`
	Port string `toml:"port"`
	Auth string `toml:"auth"`
}

type ServerExtInfo struct {
	ServerAddr   string `json:"serveraddr"`
	UserMemory   string `json:"user_memory"`
	ClientOnline string `json:"client_online"`
	ExeCommand   string `json:"exe_command"`
	RedisVer     string `json:"redis_verion"`
	KeyNumber    string `json:"key_number"`
}

type RedisTypeName struct {
	TypeNames         []TypeName `json:"typename"`
	KeysNamesWithType struct {
		KeysNames    []KeyName `json:"keysname"`
		SelfTypeName TypeName  `json:"selftypename"`
	} `json:"keysnameswithtype"`
	Contents string `json:"content"`
}

type RedisKeyName struct {
	KeysNamesWithType struct {
		KeysNames    []KeyName `json:"keysname"`
		SelfTypeName TypeName  `json:"selftypename"`
	} `json:"keysnameswithtype"`
	Contents string `json:"content"`
}

type TypeName struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type KeyName struct {
	Name  string `json:"name"`
	Index int    `json:"index"`
	Score int    `json:"score"`
}

type MonitorInfo struct {
	Message string `json:"message"`
	Err     error  `json:"err"`
}
type Int64Slice []int64

type TypeNames []TypeName
