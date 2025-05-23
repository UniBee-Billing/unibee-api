package config

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"
	"os"
	"strconv"
	"strings"
	"unibee/utility"
)

const DefaultConfigFileName = "config.yaml"

var (
	env                              string
	mode                             string
	unibeeApiUrl                     string
	serverAddress                    string
	serverJwtKey                     string
	swaggerPath                      string
	redisAddress                     string
	redisPass                        string
	redisDatabase                    string
	redisMaxIdle                     string
	redisMinIdle                     string
	redisIdleTimeout                 string
	databaseLink                     string
	databaseDebug                    string
	databaseCharset                  string
	authLoginExpire                  *int64
	authLoginExpireStr               string
	loggerLevel                      string
	nacosIpArg                       string
	nacosPortArg                     string
	nacosNamespaceArg                string
	nacosGroupArg                    string
	nacosDataIdArg                   string
	VatNumberUnExemptionCountryCodes string
)

func Init() {

	flag.StringVar(&env, "env", utility.GetEnvParam("env"), "local|daily|prod")
	flag.StringVar(&mode, "mode", utility.GetEnvParam("mode"), "stand-alone|cloud")
	flag.StringVar(&unibeeApiUrl, "unibee-api-url", utility.GetEnvParam("unibee.api.url"), "url, default http://127.0.0.1:8088")
	flag.StringVar(&serverAddress, "server-address", utility.GetEnvParam("server.address"), "server address, default :8088")
	flag.StringVar(&serverJwtKey, "server-jwtKey", utility.GetEnvParam("server.jwtKey"), "jwtKey to encrypt")
	flag.StringVar(&swaggerPath, "server-swaggerPath", utility.GetEnvParam("server.swaggerPath"), "swaggerPath, default /swagger")
	flag.StringVar(&redisAddress, "redis-address", utility.GetEnvParam("redis.address"), "redis address, require")
	flag.StringVar(&redisPass, "redis-password", utility.GetEnvParam("redis.password"), "redis password, require")
	flag.StringVar(&redisDatabase, "redis-database", utility.GetEnvParam("redis.database"), "redis database, default 0")
	flag.StringVar(&redisMaxIdle, "redis-maxIdle", utility.GetEnvParam("redis.maxIdle"), "redis maxIdle, default 500")
	flag.StringVar(&redisMinIdle, "redis-minIdle", utility.GetEnvParam("redis.minIdle"), "redis minIdle, default 10")
	flag.StringVar(&redisIdleTimeout, "redis-idleTimeout", utility.GetEnvParam("redis.idleTimeout"), "redis idleTimeout, default 1d")
	flag.StringVar(&databaseLink, "database-link", utility.GetEnvParam("database.link"), "database link, require")
	flag.StringVar(&databaseDebug, "database-debug", utility.GetEnvParam("database.debug"), "database debug, default false")
	flag.StringVar(&databaseCharset, "database-charset", utility.GetEnvParam("database.charset"), "database charset, default utf8mb4")
	flag.StringVar(&loggerLevel, "logger-level", utility.GetEnvParam("logger.level"), "logger level, default all")
	flag.StringVar(&authLoginExpireStr, "auth-login-expire", utility.GetEnvParam("auth.login.expire"), "login token expire time, default 600")
	flag.StringVar(&nacosIpArg, "nacos-ip", utility.GetEnvParam("nacos.ip"), "ip or domain, env params will replaced if nacos used")
	flag.StringVar(&nacosPortArg, "nacos-port", utility.GetEnvParam("nacos.port"), "nacos port, 8848")
	flag.StringVar(&nacosNamespaceArg, "nacos-namespace", utility.GetEnvParam("nacos.namespace"), "nacos namespace, default")
	flag.StringVar(&nacosGroupArg, "nacos-group", utility.GetEnvParam("nacos.group"), "nacos group")
	flag.StringVar(&nacosDataIdArg, "nacos-data-id", utility.GetEnvParam("nacos.data.id"), "nacos dataid like unibee-settings.yaml")
	flag.StringVar(&VatNumberUnExemptionCountryCodes, "vat-number-un-exemption-country-codes", utility.GetEnvParam("vat.number.un.exemption.country.codes"), "vat config, vat number not exemption countryCodes")

	var ctx = gctx.New()
	g.Cfg().GetAdapter().(*gcfg.AdapterFile).SetFileName(DefaultConfigFileName)

	if len(authLoginExpireStr) > 0 {
		t, _ := strconv.ParseInt(authLoginExpireStr, 10, 64)
		if t > 0 {
			authLoginExpire = &t
		}
	}

	// Parse Params
	flag.Parse()
	if len(nacosIpArg) > 0 {
		_ = deleteFile(DefaultConfigFileName) //delete old config file
		uPort, err := strconv.ParseUint(nacosPortArg, 10, 64)
		if err != nil {
			fmt.Println("Get Nacos Port:", err)
			panic(err)
		}
		fmt.Printf("Nacos IP:%s \n", nacosIpArg)
		fmt.Printf("Nacos Port:%d \n", uPort)
		fmt.Printf("Nacos Namespace:%s \n", nacosNamespaceArg)
		fmt.Printf("Nacos Group:%s \n", nacosGroupArg)
		fmt.Printf("Nacos DataId:%s \n", nacosDataIdArg)

		_, _ = ReplaceConfigContentUserNacos(strings.Trim(nacosIpArg, " "), uPort, strings.Trim(nacosNamespaceArg, " "), strings.Trim(nacosDataIdArg, " "), strings.Trim(nacosGroupArg, " "))
	} else {
		_, err := os.Stat(DefaultConfigFileName)
		if os.IsNotExist(err) || err != nil {
			if os.IsNotExist(err) {
				g.Log().Warningf(ctx, fmt.Sprintf("%s not found\n", DefaultConfigFileName))
			}
			g.Log().Warningf(ctx, "Get Config File %s Error:%s\n", DefaultConfigFileName, err.Error())
			config := map[string]interface{}{
				"server": map[string]interface{}{},
				"redis": map[string]interface{}{
					"default": map[string]interface{}{},
				},
				"database":  map[string]interface{}{"default": map[string]interface{}{}},
				"logger":    map[string]interface{}{},
				"vatConfig": map[string]interface{}{},
				"auth":      map[string]interface{}{"login": map[string]interface{}{}},
			}
			g.Cfg().GetAdapter().(*gcfg.AdapterFile).SetContent(utility.MarshalToJsonString(config), DefaultConfigFileName)
		}
	}

	SetupDefaultConfigs(ctx)

	// print configs
	fmt.Printf("Env:")
	fmt.Println(gcfg.Instance().Get(ctx, "env"))
	fmt.Printf("mode:")
	fmt.Println(gcfg.Instance().Get(ctx, "mode"))
	fmt.Println("Server Config:")
	fmt.Println(gcfg.Instance().Get(ctx, "server"))
	fmt.Println("Logger Config:")
	fmt.Println(gcfg.Instance().Get(ctx, "logger"))
	fmt.Println("Database Config:")
	fmt.Println(gcfg.Instance().Get(ctx, "database"))
	fmt.Println("Redis Config:")
	fmt.Println(gcfg.Instance().Get(ctx, "redis"))
	fmt.Println("Auth Config:")
	fmt.Println(gcfg.Instance().Get(ctx, "auth"))
}

type Nacos struct {
	ip                                       string
	namespace, dataId, group, configFilePath string
	port                                     uint64
}

func SetupDefaultConfigs(ctx context.Context) {
	// init default configs
	config := g.Cfg().MustGet(ctx, ".").Map()
	utility.Assert(config != nil, "config not found")
	setUpDefaultConfig(config, "env", env, "prod")
	setUpDefaultConfig(config, "mode", mode, "standalone")
	setUpDefaultConfig(config, "logger", map[string]interface{}{}, map[string]interface{}{})
	setUpDefaultConfig(config, "auth", map[string]interface{}{"login": map[string]interface{}{}}, map[string]interface{}{"login": map[string]interface{}{}})
	serverConfig := g.Cfg().MustGet(ctx, "server").Map()
	utility.Assert(serverConfig != nil, "server config not found")
	serverConfig["dumpRouterMap"] = false
	setUpDefaultConfig(serverConfig, "address", serverAddress, ":8088")
	setUpDefaultConfig(serverConfig, "domainPath", unibeeApiUrl, "http://127.0.0.1:8088")
	setUpDefaultConfig(serverConfig, "jwtKey", serverJwtKey, "3^&secret-key-for-UniBee*1!8*")
	serverConfig["openapiPath"] = "/api.json"
	setUpDefaultConfig(serverConfig, "swaggerPath", swaggerPath, "") ///swagger
	if serverConfig["domainPath"] == nil {
		glog.Errorf(ctx, "server.domainPath not set")
	}
	redisConfig := g.Cfg().MustGet(ctx, "redis.default").Map()
	_redisDatabaseInt := 0
	if len(redisDatabase) > 0 {
		_redisDatabaseInt, _ = strconv.Atoi(redisDatabase)
	}
	utility.Assert(redisConfig != nil, "redis config not found")
	setUpDefaultConfig(redisConfig, "address", redisAddress, nil)
	setUpDefaultConfig(redisConfig, "pass", redisPass, nil)
	setUpDefaultConfig(redisConfig, "db", _redisDatabaseInt, 0)
	setUpDefaultConfig(redisConfig, "maxIdle", redisMaxIdle, 500)
	setUpDefaultConfig(redisConfig, "minIdle", redisMinIdle, 10)
	setUpDefaultConfig(redisConfig, "idleTimeout", redisIdleTimeout, "1d")
	databaseConfig := g.Cfg().MustGet(ctx, "database.default").Map()
	utility.Assert(databaseConfig != nil, "database config not found")
	setUpDefaultConfig(databaseConfig, "link", databaseLink, nil)
	setUpDefaultConfig(databaseConfig, "debug", databaseDebug, false)
	setUpDefaultConfig(databaseConfig, "charset", databaseCharset, "utf8mb4")
	loggerConfig := g.Cfg().MustGet(ctx, "logger").Map()
	utility.Assert(loggerConfig != nil, "logger config not found")
	setUpDefaultConfig(loggerConfig, "level", loggerLevel, "all")
	setUpDefaultConfig(loggerConfig, "stdout", true, true)
	authLoginConfig := g.Cfg().MustGet(ctx, "auth.login").Map()
	utility.Assert(authLoginConfig != nil, "auth login config not found")
	if authLoginExpire != nil {
		setUpDefaultConfig(authLoginConfig, "expire", authLoginExpire, 600)
	} else {
		setUpDefaultConfig(authLoginConfig, "expire", 600, 600)
	}

	//vatConfig := g.Cfg().MustGet(ctx, "vatConfig").Map()
	//if vatConfig != nil {
	//	setUpDefaultConfig(vatConfig, "nonEuEnable", VatNonEuEnable, "false")
	//	setUpDefaultConfig(vatConfig, "numberUnExemptionCountryCodes", VatNumberUnExemptionCountryCodes, "")
	//}
	g.Cfg().GetAdapter().(*gcfg.AdapterFile).SetContent(utility.MarshalToJsonString(config), DefaultConfigFileName)
	SetConfig(utility.MarshalToJsonString(config))
	if VatNumberUnExemptionCountryCodes != "" {
		GetConfigInstance().VatConfig.NumberUnExemptionCountryCodes = VatNumberUnExemptionCountryCodes
	}
}

func ReplaceConfigContentUserNacos(ip string, port uint64, namespace, dataId, group string) (n *Nacos, err error) {
	n = &Nacos{
		ip:        ip,
		port:      port,
		namespace: namespace,
		dataId:    dataId,
		group:     group,
	}
	config, err := GetNacosConfig(n.ip, n.port, n.namespace, n.group, n.dataId)
	if err != nil {
		panic(err)
	}
	g.Cfg().GetAdapter().(*gcfg.AdapterFile).SetContent(config, DefaultConfigFileName)
	return
}

func (n Nacos) GetConfigFilePath() string {
	if len(n.configFilePath) == 0 {
		panic("nacos config to save local file is not found!")
	}
	return n.configFilePath
}

func (n *Nacos) syncToFile() (err error) {
	config, err := GetNacosConfig(n.ip, n.port, n.namespace, n.group, n.dataId)
	if err != nil {
		fmt.Println("nacos config load failure")
		panic(err)
	}
	file, err := createFile(DefaultConfigFileName)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic("file close error")
		}
	}(file)

	if file == nil {
		panic("create or read file error")
	}
	writer := bufio.NewWriter(file)
	_, err = writer.WriteString(config)
	err = writer.Flush()
	n.configFilePath = DefaultConfigFileName

	return
}

func createFile(path string) (file *os.File, err error) {
	file, err = os.Create(path)
	if err != nil {
		panic("create file " + err.Error())
	}
	return
}

func deleteFile(path string) (err error) {
	err = os.Remove(path)
	return
}

func setUpDefaultConfig(config map[string]interface{}, key string, flagValue interface{}, defaultValue interface{}) {
	if config[key] == nil {
		if flagValue != nil && flagValue != "" {
			config[key] = flagValue
		} else {
			config[key] = defaultValue
		}
	}
}
