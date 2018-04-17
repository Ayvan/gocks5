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
    "fmt"
    "path/filepath"
    "io/ioutil"
    "gopkg.in/yaml.v2"
    "os"
)

type ProxyConfig struct {
    Host     string `yaml:"host"`
    Port     int    `yaml:"port"`
    User     string `yaml:"user"`
    Password string `yaml:"password"`
    LogPath  string `yaml:"log_path"`
}

func LoadConfig(path string, to, base interface{}) error {

    configPath, err := filepath.Abs(path)
    if err != nil {
        return fmt.Errorf("file path error: %s %s", configPath, err)
    }

    content, err := ioutil.ReadFile(configPath)
    if err != nil && !os.IsNotExist(err) {
        return fmt.Errorf("config file (%s) read error: %v\n", configPath, err)
    } else if os.IsNotExist(err) {
        return nil
    }

    err = yaml.Unmarshal(content, to)
    if err != nil {
        return fmt.Errorf("config file %s parsing error: %v", configPath, err)
    }

    err = yaml.Unmarshal(content, base)
    if err != nil {
        return fmt.Errorf("config file %s parsing error: %v", configPath, err)
    }

    return nil
}
