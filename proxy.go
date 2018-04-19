// Copyright 2018 Ivan Korostelev. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
    "github.com/jessevdk/go-flags"
    "os"
    "github.com/armon/go-socks5"
    "fmt"
    "log"
    "github.com/ayvan/godaemon"
    "runtime/debug"
    "path/filepath"
)

var options struct {
    ConfigPath string `short:"c" long:"config" default:"gocks5.yml" description:"config file path"`
    LogPath    string `short:"l" long:"log" default:"stdout" description:"log file path"`
    Daemon     bool   `short:"d" long:"daemon" description:"start as daemon (need to set config path)"`
}

func main() {
    start()
}

func start() {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("%s", r)
            log.Printf("%s", string(debug.Stack()))
        }
    }()

    parser := flags.NewParser(&options, flags.Default)

    if _, err := parser.ParseArgs(os.Args); err != nil {
        log.Printf("parse args error: %s", err)
        os.Exit(1)
    }

    if options.Daemon {
        wd, err := os.Getwd()
        if err != nil {
            log.Printf("getwd error: %s", err)
        }

        // set working dir for simple load config if --config argument not absolute path
        godaemon.MakeDaemon(&godaemon.DaemonAttr{WorkingDir: wd})

        if parser.FindOptionByLongName("config").IsSetDefault() {
            // if config path not set and start as daemon - need to set path in /etc
            options.ConfigPath = "/etc/gocks5.yml"
        }
    }

    config := ProxyConfig{
        Host:    "localhost",
        Port:    1080,
        LogPath: options.LogPath,
    }

    if err := LoadConfig(options.ConfigPath, &config, &config); err != nil {
        log.Printf("load config error: %s", err)
        os.Exit(1)
    }

    // set high priority for log path from command-line arguments
    if !parser.FindOptionByLongName("log").IsSetDefault() {
        config.LogPath = options.LogPath
    }

    if err := initLog(config.LogPath); err != nil {
        log.Printf("init log error: %s", err)
        os.Exit(1)
    }

    if err := startProxy(config); err != nil {
        log.Print(err)
        os.Exit(1)
    }
}

type credentials struct {
    User     string
    Password string
}

func (c credentials) Valid(user, password string) bool {
    if c.User == user && c.Password == password {
        return true
    }
    log.Printf("failed to authenticate, bad username or password: %s %s", user, password)

    return false
}

func startProxy(config ProxyConfig) error {

    // Create a SOCKS5 server
    conf := &socks5.Config{
        AuthMethods: []socks5.Authenticator{
            socks5.UserPassAuthenticator{
                Credentials: credentials{User: config.User, Password: config.Password},
                //Credentials: socks5.StaticCredentials{config.User:config.Password},
            },
        },
    }
    server, err := socks5.New(conf)
    if err != nil {
        return err
    }

    address := fmt.Sprintf("%s:%d", config.Host, config.Port)

    log.Printf("start socks 5 proxy on %s", address)

    // Create SOCKS5 proxy on localhost port 8000
    if err := server.ListenAndServe("tcp", address); err != nil {
        return err
    }

    return nil
}

func initLog(logPath string) error {
    if logPath == "" || logPath == "stdout" {
        log.SetOutput(os.Stdout)
        return nil
    }

    logPath, err := filepath.Abs(logPath)
    if err != nil {
        return fmt.Errorf("log file path error: %s %s", logPath, err)
    }
    log.SetOutput(os.Stdout)
    if logFile, err := os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664); err != nil {
        return err
    } else {
        log.SetOutput(logFile)
    }

    return nil
}
