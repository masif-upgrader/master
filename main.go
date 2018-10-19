//go:generate go run vendor/github.com/Al2Klimov/go-gen-source-repos/main.go github.com/masif-upgrader/master
//go:generate go run gen-mysql.go

package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	_ "github.com/Al2Klimov/go-gen-source-repos"
	"github.com/go-ini/ini"
	_ "github.com/masif-upgrader/common"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"strings"
)

type settings struct {
	api struct {
		listen string
	}
	tls struct {
		cert, key, ca, crl string
	}
	db struct {
		typ, dsn string
	}
	log struct {
		level log.Level
	}
}

var logLevels = map[string]log.Level{
	"error":   log.ErrorLevel,
	"err":     log.ErrorLevel,
	"warning": log.WarnLevel,
	"warn":    log.WarnLevel,
	"info":    log.InfoLevel,
	"debug":   log.DebugLevel,
}

var db *sql.DB = nil

func main() {
	if len(os.Args) == 1 && terminal.IsTerminal(int(os.Stdout.Fd())) {
		fmt.Printf(
			"For the terms of use, the source code and the authors\n"+
				"see the projects this program is assembled from:\n\n  %s\n",
			strings.Join(GithubcomAl2klimovGo_gen_source_repos, "\n  "),
		)
		os.Exit(1)
	}

	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	if err := runMaster(); err != nil {
		log.Fatal(err)
	}
}

func runMaster() error {
	cfg, errLC := loadCfg()
	if errLC != nil {
		return errLC
	}

	log.SetLevel(cfg.log.level)

	httpd, errNA := newApi(cfg.api.listen, cfg.tls)
	if errNA != nil {
		return errNA
	}

	log.Info("Loading SQL schema")

	var errDB error
	if db, errDB = sql.Open(cfg.db.typ, cfg.db.dsn); errDB != nil {
		return errDB
	}

	for _, ddl := range mysqlDdls {
		for {
			if _, errExec := db.Exec(ddl); errExec != nil {
				if isRecoverableDbError(errExec) {
					continue
				}

				return errExec
			}

			break
		}
	}

	log.Info("Starting HTTPd")

	return httpd.ListenAndServeTLS("", "")
}

func loadCfg() (config *settings, err error) {
	cfgFile := flag.String("config", "", "config file")
	flag.Parse()

	if *cfgFile == "" {
		return nil, errors.New("config file missing")
	}

	log.WithFields(log.Fields{"file": *cfgFile}).Debug("Loading config")

	cfg, errLI := ini.Load(*cfgFile)
	if errLI != nil {
		return nil, errLI
	}

	cfgTls := cfg.Section("tls")
	cfgDb := cfg.Section("db")
	result := &settings{
		api: struct{ listen string }{
			listen: cfg.Section("api").Key("listen").String(),
		},
		tls: struct{ cert, key, ca, crl string }{
			cert: cfgTls.Key("cert").String(),
			key:  cfgTls.Key("key").String(),
			ca:   cfgTls.Key("ca").String(),
			crl:  cfgTls.Key("crl").String(),
		},
		db: struct{ typ, dsn string }{
			typ: cfgDb.Key("type").String(),
			dsn: cfgDb.Key("dsn").String(),
		},
	}

	if result.api.listen == "" {
		return nil, errors.New("config: api.listen missing")
	}

	if result.tls.cert == "" {
		return nil, errors.New("config: tls.cert missing")
	}

	if result.tls.key == "" {
		return nil, errors.New("config: tls.key missing")
	}

	if result.tls.ca == "" {
		return nil, errors.New("config: tls.ca missing")
	}

	if result.db.typ == "" {
		return nil, errors.New("config: db.type missing")
	}

	if result.db.dsn == "" {
		return nil, errors.New("config: db.dsn missing")
	}

	if rawLogLvl := cfg.Section("log").Key("level").String(); rawLogLvl == "" {
		result.log.level = log.InfoLevel
	} else if logLvl, logLvlValid := logLevels[rawLogLvl]; logLvlValid {
		result.log.level = logLvl
	} else {
		return nil, errors.New("config: bad log.level")
	}

	return result, nil
}
