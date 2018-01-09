package main

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func init() {
	viper.BindEnv("listen", "LOCATOR_LISTEN")
	viper.SetDefault("listen", "127.0.0.1:8000")

	viper.BindEnv("debug", "LOCATOR_DEBUG")
	viper.SetDefault("debug", "false")
}

func main() {
	if !viper.GetBool("debug") {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	r.GET("/", handleClient)
	r.GET("/:addr", handleIP)
	if err := r.Run(viper.GetString("listen")); err != nil {
		fmt.Fprint(os.Stderr, err.Error())
	}
}

func handleClient(c *gin.Context) {
	handle(c, c.ClientIP())
}

func handleIP(c *gin.Context) {
	handle(c, c.Param("addr"))
}

func handle(c *gin.Context, addr string) {
	ip := net.ParseIP(addr)
	if ip == nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	res, err := geoipAddr(ip)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, res)
}
