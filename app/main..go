package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pi-sin/go-repo-structure/config"
	"github.com/pi-sin/go-repo-structure/middleware"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"log"
	"net/url"
	"time"

	_articleHttpDelivery "github.com/pi-sin/go-repo-structure/article/delivery/http"
	_articleManager "github.com/pi-sin/go-repo-structure/article/manager"
	_articleRepo "github.com/pi-sin/go-repo-structure/article/repository/mysql"
	_authorRepo "github.com/pi-sin/go-repo-structure/author/repository/mysql"
)

func main() {
	dbHost := config.GetConfig().GetString(`database.host`)
	dbPort := config.GetConfig().GetString(`database.port`)
	dbUser := config.GetConfig().GetString(`database.user`)
	dbPass := config.GetConfig().GetString(`database.pass`)
	dbName := config.GetConfig().GetString(`database.name`)
	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)

	val := url.Values{}
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Jakarta")
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	dbConn, err := sql.Open(`mysql`, dsn)

	if err != nil {
		//log error if required
		panic(err)
	}

	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	server := gin.New()

	server.Use(gin.Recovery())
	server.Use(middleware.ApmMiddleware())
	server.Use(middleware.CorsMiddleware())

	authorRepo := _authorRepo.NewMysqlAuthorRepository(dbConn)
	ar := _articleRepo.NewMysqlArticleRepository(dbConn)

	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second
	am := _articleManager.NewArticleManager(ar, authorRepo, timeoutContext)
	_articleHttpDelivery.NewHttpArticleHandler(server, am)

	if err := server.Run(config.GetConfig().GetString("http.port")); err != nil {
		logrus.WithFields(logrus.Fields{"mod": "server", "evn": "listen"}).Error(err)
	}
}
