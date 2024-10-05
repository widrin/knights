// Copyright 2015 Daniel Theophanes.
// Use of this source code is governed by a zlib-style
// license that can be found in the LICENSE file.

// Simple service that only works by printing a log message every few seconds.
package main

import (
	"fmt"
	"gopkg.in/ini.v1"
	"os"
	"tzgit.kaixinxiyou.com/utils/tzid/tzIdService/workid"
)

func main() {
	work := workid.Create()
	var err error
	cfg, err := ini.LoadSources(ini.LoadOptions{AllowNonUniqueSections: true}, "config.conf")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	redisAddr := cfg.Section("").Key("redis_addr").Value()
	redisPass := cfg.Section("").Key("redis_pass").Value()
	redisUser := cfg.Section("").Key("redis_user").Value()
	work.Init(redisAddr, redisUser, redisPass, 0)
	work.Start()

}
