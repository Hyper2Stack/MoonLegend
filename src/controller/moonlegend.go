package main

import (
    "flag"
    "os"
    "net/http"
    "runtime"

    "controller/handler"
    "controller/migrate"
    "controller/model"

    "github.com/op/go-logging"
)

var version string = "DEV_VERSION"

type DB struct {
    Driver     string `json:"driver"`
    Endpoint   string `json:"endpoint"`
}

type Logging struct {
    File   string `json:"file"`
    Stdout bool   `json:"stdout"`
    Level  string `json:"level"`
    Format string `json:"format"`
}

type Config struct {
    Db            DB      `json:"db"`
    Logging       Logging `json:"logging"`
    Listenaddr    string  `json:"listenaddr"`
    Listenport    string  `json:"listenport"`
}

var (
    config    Config
    logLevels = map[string]logging.Level{
        "debug":    logging.DEBUG,
        "info":     logging.INFO,
        "notice":   logging.NOTICE,
        "warning":  logging.WARNING,
        "error":    logging.ERROR,
        "critical": logging.CRITICAL,
    }
    log, _ = logging.GetLogger("moonlegend")
)

func initLogger() {
    backends := []logging.Backend{}
    if config.Logging.Stdout {
        backends = append(backends, logging.NewLogBackend(os.Stdout, "", 0))
    }
    if len(config.Logging.File) > 0 {
        file, err := os.OpenFile(config.Logging.File,
            os.O_RDWR|os.O_CREATE|os.O_APPEND, 0660)
        if err != nil {
            panic(err)
        }
        backends = append(backends, logging.NewLogBackend(file, "", 0))
    }
    level, ok := logLevels[config.Logging.Level]
    if !ok {
        panic("Log level not found")
    }

    var format = logging.MustStringFormatter(config.Logging.Format)
    logging.SetBackend(backends...)
    logging.SetFormatter(format)
    logging.SetLevel(level, "")
}

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())

    configPath := flag.String("conf",
        "./config/moonlegend.json", "Path of config file")
    flag.Parse()

    if err := LoadJsonConfig(*configPath, &config); err != nil {
        log.Error("Fail to parse config file: " + *configPath)
        os.Exit(1)
    }
    initLogger()
    log.Info("moonlegend(" + version + ") started")

    if len(config.Db.Driver)==0 || len(config.Db.Endpoint) == 0 {
        log.Error("DB driver or endpoint not found")
        os.Exit(1)
    }

    if err := model.Initialize(config.Db.Driver, config.Db.Endpoint); err != nil {
        log.Error("Init db failed: " + err.Error())
        os.Exit(1)
    }

    if err := handler.Initialize(); err != nil {
        log.Error("Init handler failed: " + err.Error())
        os.Exit(1)
    }

    if err := migrate.Run(config.Db.Endpoint); err != nil {
        log.Error("Schema migration failed: " + err.Error())
        os.Exit(1)
    }

    log.Info("Listening on " + config.Listenaddr + ":" + config.Listenport + "...")
    log.Error(http.ListenAndServe(":"+config.Listenport, handler.NewRouter()))
    os.Exit(1)
}
