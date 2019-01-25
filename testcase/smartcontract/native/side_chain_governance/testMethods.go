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

package side_chain_governance

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/ontio/ontology-crypto/keypair"
	sdk "github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology-tool/testframework"
)

type RegisterSideChainParam struct {
	Path        string
	SideChainID uint32
	Ratio       uint32
	Deposit     uint64
	OngPool     uint64
}

func RegisterSideChain(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./side_chain_params/RegisterSideChain.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	registerSideChainParam := new(RegisterSideChainParam)
	err = json.Unmarshal(data, registerSideChainParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}
	time.Sleep(1 * time.Second)
	user, ok := getAccountByPassword(ctx, registerSideChainParam.Path)
	if !ok {
		return false
	}
	ok = registerSideChain(ctx, user, registerSideChainParam.Ratio, registerSideChainParam.Deposit, registerSideChainParam.OngPool)
	if !ok {
		return false
	}
	return true
}

type ApproveSideChainParam struct {
	Path        []string
	SideChainID uint32
}

func ApproveSideChain(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./params/ApproveSideChain.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	approveSideChainParam := new(ApproveSideChainParam)
	err = json.Unmarshal(data, approveSideChainParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}
	var users []*sdk.Account
	var pubKeys []keypair.PublicKey
	time.Sleep(1 * time.Second)
	for _, path := range approveSideChainParam.Path {
		user, ok := getAccountByPassword(ctx, path)
		if !ok {
			return false
		}
		users = append(users, user)
		pubKeys = append(pubKeys, user.PublicKey)
	}
	ok := approveSideChainMultiSign(ctx, pubKeys, users, approveSideChainParam.SideChainID)
	if !ok {
		return false
	}
	waitForBlock(ctx)
	return true
}

type RejectSideChainParam struct {
	Path        []string
	SideChainID uint32
}

func RejectSideChain(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./params/RejectSideChain.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	rejectSideChainParam := new(RejectSideChainParam)
	err = json.Unmarshal(data, rejectSideChainParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}
	var users []*sdk.Account
	var pubKeys []keypair.PublicKey
	time.Sleep(1 * time.Second)
	for _, path := range rejectSideChainParam.Path {
		user, ok := getAccountByPassword(ctx, path)
		if !ok {
			return false
		}
		users = append(users, user)
		pubKeys = append(pubKeys, user.PublicKey)
	}
	ok := rejectSideChainMultiSign(ctx, pubKeys, users, rejectSideChainParam.SideChainID)
	if !ok {
		return false
	}
	waitForBlock(ctx)
	return true
}

type RegisterNodeToSideChainParam struct {
	SideChainID uint32
	PeerPubkey  []string
	Path        []string
}

func RegisterNodeToSideChain(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./side_chain_params/RegisterNodeToSideChain.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	registerNodeToSideChainParam := new(RegisterNodeToSideChainParam)
	err = json.Unmarshal(data, registerNodeToSideChainParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}
	time.Sleep(1 * time.Second)
	for index, peerPubkey := range registerNodeToSideChainParam.PeerPubkey {
		user, ok := getAccountByPassword(ctx, registerNodeToSideChainParam.Path[index])
		if !ok {
			return false
		}
		ok = registerNodeToSideChain(ctx, user, registerNodeToSideChainParam.SideChainID, peerPubkey)
		if !ok {
			return false
		}
	}
	return true
}

type SideChainID struct {
	SideChainID uint32
}

func GetSideChain(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./side_chain_params/GetSideChain.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	sideChainID := new(SideChainID)
	err = json.Unmarshal(data, sideChainID)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}
	sideChain, err := getSideChain(ctx, sideChainID.SideChainID)
	if err != nil {
		ctx.LogError("getSideChain error %v", err)
		return false
	}
	fmt.Println("SideChainID of side chain is: ", sideChain.SideChainID)
	fmt.Println("Address of side chain is: ", sideChain.Address)
	fmt.Println("Status of side chain is: ", sideChain.Status)
	fmt.Println("OngNum of side chain is: ", sideChain.OngNum)
	fmt.Println("Deposit of side chain is: ", sideChain.Deposit)
	fmt.Println("OngPool of side chain is: ", sideChain.OngPool)
	fmt.Println("Ratio of side chain is: ", sideChain.Ratio)
	return true
}

func GetSideChainNodeInfo(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./side_chain_params/GetSideChainNodeInfo.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	sideChainID := new(SideChainID)
	err = json.Unmarshal(data, sideChainID)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}
	sideChainNodeInfo, err := getSideChainNodeInfo(ctx, sideChainID.SideChainID)
	if err != nil {
		ctx.LogError("getSideChain error %v", err)
		return false
	}
	fmt.Println("SideChainID of side chain is: ", sideChainNodeInfo.SideChainID)
	fmt.Println("NodeInfoMap of side chain is: ", sideChainNodeInfo.NodeInfoMap)
	return true
}

func GetSideChainPeerPoolMap(ctx *testframework.TestFrameworkContext) bool {
	peerPoolMap, err := getSideChainPeerPoolMap(ctx)
	if err != nil {
		ctx.LogError("getPeerPoolMap failed %v", err)
		return false
	}

	for _, v := range peerPoolMap.PeerPoolMap {
		fmt.Println("###########################################")
		fmt.Println("peerPoolItem.Index is:", v.Index)
		fmt.Println("peerPoolItem.PeerPubkey is:", v.PeerPubkey)
		fmt.Println("peerPoolItem.Address is:", v.Address.ToBase58())
		fmt.Println("peerPoolItem.Status is:", v.Status)
		fmt.Println("peerPoolItem.InitPos is:", v.InitPos)
		fmt.Println("peerPoolItem.TotalPos is:", v.TotalPos)
	}
	return true
}
