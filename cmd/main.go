package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/MeNoln/go-bus-ws-notifier/pkg/busclient"
	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gopkg.in/olahol/melody.v1"
)

func main() {
	err := loadConfig()
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(-1)
	}

	r := gin.Default()
	m := melody.New()

	r.LoadHTMLGlob("*.html")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.GET("/ws", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	r.GET("/app", func(c *gin.Context) {
		c.JSON(200, struct {
			message string
		}{
			message: "im good",
		})
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		m.Broadcast(msg)
	})

	go busclient.ProcessMessages(m)

	r.Run(fmt.Sprintf(":5100"))
}

func loadConfig() error {
	viper.SetConfigName(getRunningEnv())
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../config/")

	return viper.ReadInConfig()
}

func getRunningEnv() string {
	const localCfg string = "local"
	if cfgEnv := os.Getenv("APP_ENV"); len(cfgEnv) != 0 {
		return cfgEnv
	}

	return localCfg
}
