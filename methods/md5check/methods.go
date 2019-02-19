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

package md5check

import (
	"encoding/json"
	"io/ioutil"

	log4 "github.com/alecthomas/log4go"
	sdk "github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology-tool/common"
)

type Md5CheckParam struct {
	Path string
}

func Md5Check(ontSdk *sdk.OntologySdk) bool {
	data, err := ioutil.ReadFile("./params/Md5Check.json")
	if err != nil {
		log4.Error("ioutil.ReadFile failed ", err)
		return false
	}
	md5CheckParam := new(Md5CheckParam)
	err = json.Unmarshal(data, md5CheckParam)
	if err != nil {
		log4.Error("json.Unmarshal failed ", err)
		return false
	}
	ok := md5Check(ontSdk, md5CheckParam.Path)
	if !ok {
		return false
	}
	common.WaitForBlock(ontSdk)
	return true
}
