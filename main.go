package main

import (
	"database/sql"
	_ "embed"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/pelletier/go-toml/v2"
	_ "modernc.org/sqlite"
)

type Config struct {
	Twitch struct {
		ClientID     string `toml:"client_id"`
		ClientSecret string `toml:"client_secret"`
		RedirectURI  string `toml:"redirect_uri"`
		ChannelName  string `toml:"channel_name"`
		Reward       struct {
			Name string
			Cost int
		} `toml:"reward"`
	} `toml:"twitch"`
}

var (
	cfg Config
	db  *sql.DB

	//go:embed init.sql
	initSQL string
)

func init() {
	f, err := os.OpenFile("latest.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Panic("log file error:", err)
	}
	defer f.Close()

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(io.MultiWriter(os.Stdout, f))

	f, err = os.Open("./config.toml")
	if err != nil {
		log.Panic("config error:", err)
	}
	defer f.Close()

	b, _ := io.ReadAll(f)
	err = toml.Unmarshal(b, &cfg)
	if err != nil {
		log.Panic(err)
	}

	db, err := sql.Open("sqlite", "./data.db")
	if err != nil {
		log.Panic("invalid db config:", err)
	}
	if err = db.Ping(); err != nil {
		log.Panic("db unreachable:", err)
	}

	stmt, err := db.Prepare(initSQL)
	if err != nil {
		log.Panic("init sql:", err)
	}
	stmt.Exec()

	log.Println("init() executed successfully")
}

func main() {
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, World! :3")
	})

	go func() {
		err := http.ListenAndServe("127.0.0.1:3000", nil)
		if err != nil {
			log.Panic(err)
		}
	}()

	<-exit
}
