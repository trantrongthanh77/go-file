package main

import (
	"fmt"
	"go-file/common"
	"go-file/common/config"
	"go-file/model"
	"go-file/router"
	"html/template"
	"log"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func loadTemplate() *template.Template {
	var funcMap = template.FuncMap{
		"unescape": common.UnescapeHTML,
	}
	t := template.Must(template.New("").Funcs(funcMap).ParseFS(common.FS, "public/*.html"))
	return t
}

func main() {
	conf, err := buildConfig()
	if err != nil {
		log.Fatalf("failed to build config: %v", err)
	}
	common.InitConfig(conf)
	common.SetupGinLog(conf)
	common.SysLog(fmt.Sprintf("Go File %s started at port %d", common.Version, conf.Port))
	if conf.GinMode != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}
	// Initialize SQL Database
	db, err := model.InitDB()
	if err != nil {
		common.FatalLog(err)
	}
	defer func(db *gorm.DB) {
		err := db.Close()
		if err != nil {
			common.FatalLog("failed to close database: " + err.Error())
		}
	}(db)

	// Initialize Redis
	err = common.InitRedisClient(conf)
	if err != nil {
		common.FatalLog(err)
	}

	// Initialize options
	model.InitOptionMap()

	// Initialize HTTP server
	server := gin.Default()
	server.SetHTMLTemplate(loadTemplate())

	// Initialize session store
	var store sessions.Store
	if common.RedisEnabled {
		opt := common.ParseRedisOption(conf)
		store, _ = redis.NewStore(opt.MinIdleConns, opt.Network, opt.Addr, opt.Password, []byte(common.SessionSecret))
	} else {
		store = cookie.NewStore([]byte(common.SessionSecret))
	}
	store.Options(sessions.Options{
		HttpOnly: true,
	})
	server.Use(sessions.Sessions("session", store))

	router.SetRouter(server, conf)
	if conf.Host == "" {
		ip := common.GetIp()
		if ip != "" {
			conf.Host = ip
		}
	}
	strPort := strconv.Itoa(conf.Port)
	serverUrl := fmt.Sprintf("http://%s:%s/", conf.Host, strPort)
	if !conf.Nobrowser {
		common.OpenBrowser(serverUrl)
	}
	if conf.Enablep2p {
		go common.StartP2PServer(conf)
	}
	err = server.Run(":" + strPort)
	if err != nil {
		common.FatalLog(err)
	}
}

func buildConfig() (*config.Config, error) {
	c := &config.Config{}
	return c, c.ApplyConfig(
		config.FromEnv,
	)
}
