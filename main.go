/*
 * Copyright (C) 2018 The ontology Authors
 * This file is part of The ontology library.
 *
 * The ontology is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The ontology is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The ontology.  If not, see <http://www.gnu.org/licenses/>.
 */
package main

import (
	"flag"
	"github.com/ontio/ontology-tool/config"
	"github.com/ontio/ontology-tool/core"
	"github.com/ontio/ontology-tool/log"
	_ "github.com/ontio/ontology-tool/methods"
	"math/rand"
	"strings"
	"time"
)

var (
	Config  string //config file
	Methods string //Methods list in cmdline
)

func init() {
	flag.StringVar(&Config, "cfg", "./config.json", "Config of ontology-tool")
	flag.StringVar(&Methods, "t", "", "methods to run. use ',' to split methods")
	flag.Parse()
}

func main() {
	rand.Seed(time.Now().UnixNano())
	defer time.Sleep(time.Second)

	err := config.DefConfig.Init(Config)
	if err != nil {
		log.Error("DefConfig.Init error:%s", err)
		return
	}

	methods := make([]string, 0)
	if Methods != "" {
		methods = strings.Split(Methods, ",")
	}

	core.OntTool.Start(methods)
}
