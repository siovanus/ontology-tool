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

package core

import (
	log4 "github.com/alecthomas/log4go"
	sdk "github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology-tool/config"
)

var OntTool = NewOntologyTool()

type Method func(sdk *sdk.OntologySdk) bool

type OntologyTool struct {
	//Map name to method
	methodsMap map[string]Method
	//Map method result
	methodsRes map[string]bool
}

func NewOntologyTool() *OntologyTool {
	return &OntologyTool{
		methodsMap: make(map[string]Method, 0),
		methodsRes: make(map[string]bool, 0),
	}
}

func (this *OntologyTool) RegMethod(name string, method Method) {
	this.methodsMap[name] = method
}

//Start run
func (this *OntologyTool) Start(methodsList []string) {
	if len(methodsList) > 0 {
		this.runMethodList(methodsList)
		return
	}
	log4.Info("No method to run")
	return
}

func (this *OntologyTool) runMethodList(methodsList []string) {
	this.onStart()
	defer this.onFinish(methodsList)
	ontSdk := sdk.NewOntologySdk()
	ontSdk.NewRpcClient().SetAddress(config.DefConfig.JsonRpcAddress)
	for i, method := range methodsList {
		this.runMethod(i+1, ontSdk, method)
	}
}

func (this *OntologyTool) runMethod(index int, sdk *sdk.OntologySdk, methodName string) {
	this.onBeforeMethodStart(index, methodName)
	method := this.getMethodByName(methodName)
	if method != nil {
		ok := method(sdk)
		this.onAfterMethodFinish(index, methodName, ok)
		this.methodsRes[methodName] = ok
	}
}

func (this *OntologyTool) onStart() {
	log4.Info("===============================================================")
	log4.Info("-------Ontology Tool Start-------")
	log4.Info("===============================================================")
	log4.Info("")
}

func (this *OntologyTool) onFinish(methodsList []string) {
	failedList := make([]string, 0)
	successList := make([]string, 0)
	for methodName, ok := range this.methodsRes {
		if ok {
			successList = append(successList, methodName)
		} else {
			failedList = append(failedList, methodName)
		}
	}

	skipList := make([]string, 0)
	for _, method := range methodsList {
		_, ok := this.methodsRes[method]
		if !ok {
			skipList = append(skipList, method)
		}
	}

	succCount := len(successList)
	failedCount := len(failedList)

	log4.Info("===============================================================")
	log4.Info("Ontology Tool Finish Total:%v Success:%v Failed:%v Skip:%v",
		len(methodsList),
		succCount,
		failedCount,
		len(methodsList)-succCount-failedCount)
	if succCount > 0 {
		log4.Info("---------------------------------------------------------------")
		log4.Info("Success list:")
		for i, succ := range successList {
			log4.Info("%d.\t%s", i+1, succ)
		}
	}
	if failedCount > 0 {
		log4.Info("---------------------------------------------------------------")
		log4.Info("Fail list:")
		for i, fail := range failedList {
			log4.Info("%d.\t%s", i+1, fail)
		}
	}
	if len(skipList) > 0 {
		log4.Info("---------------------------------------------------------------")
		log4.Info("Skip list:")
		for i, skip := range skipList {
			log4.Info("%d.\t%s", i+1, skip)
		}
	}
	log4.Info("===============================================================")
}

func (this *OntologyTool) onBeforeMethodStart(index int, methodName string) {
	log4.Info("===============================================================")
	log4.Info("%d. Start Method:%s", index, methodName)
	log4.Info("---------------------------------------------------------------")
}

func (this *OntologyTool) onAfterMethodFinish(index int, methodName string, res bool) {
	if res {
		log4.Info("Run Method:%s success.", methodName)
	} else {
		log4.Info("Run Method:%s failed.", methodName)
	}
	log4.Info("---------------------------------------------------------------")
	log4.Info("")
}

func (this *OntologyTool) getMethodByName(name string) Method {
	return this.methodsMap[name]
}
