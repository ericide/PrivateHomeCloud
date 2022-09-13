package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"net/http"
	"service.chat/internal/config"
	"service.chat/internal/confx"
	"service.chat/internal/handler"
	"service.chat/internal/svc"
	"service.chat/internal/ws"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/chat-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	confx.MustLoad(*configFile, &c)

	ctx := svc.NewServiceContext(c)
	migrateMysql(c)
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	handler.RegisterHandlers(server, ctx)

	server.AddRoute(rest.Route{
		Method: http.MethodGet,
		Path:   "/",
		Handler: func(writer http.ResponseWriter, request *http.Request) {
			http.ServeFile(writer, request, "../test/index.html")
		},
	})

	hub := ws.NewHub(ctx)
	go hub.Run()
	server.AddRoute(rest.Route{
		Method: http.MethodGet,
		Path:   "/ws",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			ws.ServeWs(hub, ctx, w, r)
		},
	})

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

func migrateMysql(c config.Config) {
	db, err := sql.Open("mysql", c.Mysql.DataSource)
	if err != nil {
		panic(err)
	}
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	migrateFolder := "file://./migrate"
	if confx.MatchEnv(confx.EnvTypeLocal) {
		migrateFolder = "file://./migrate"
	}
	m, err := migrate.NewWithDatabaseInstance(
		migrateFolder,
		"mysql", driver)
	if err != nil {
		panic(err)
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		panic(err)
	}
}
