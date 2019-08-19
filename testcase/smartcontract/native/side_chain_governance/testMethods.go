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
	"github.com/ontio/multi-chain/smartcontract/service/native/side_chain_manager"
	"io/ioutil"
	"time"

	"github.com/ontio/ontology-crypto/keypair"
	sdk "github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology-tool/testframework"
)

type SyncGenesisHeaderParam struct {
	Path     []string
	ChainRpc string
}

func SyncGenesisHeader(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./side_chain_params/SyncGenesisHeader.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	SyncGenesisHeaderParam := new(SyncGenesisHeaderParam)
	err = json.Unmarshal(data, SyncGenesisHeaderParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}

	var users []*sdk.Account
	var pubKeys []keypair.PublicKey
	time.Sleep(1 * time.Second)
	for _, path := range SyncGenesisHeaderParam.Path {
		user, ok := getAccountByPassword(ctx, path)
		if !ok {
			return false
		}
		users = append(users, user)
		pubKeys = append(pubKeys, user.PublicKey)
	}

	sideSdk := sdk.NewOntologySdk()
	sideSdk.NewRpcClient().SetAddress(SyncGenesisHeaderParam.ChainRpc)
	genesisBlock, err := sideSdk.GetBlockByHeight(0)
	if err != nil {
		ctx.LogError("get side chain genesis block error: %s", err)
		return false
	}
	genesisBlockHeader := genesisBlock.Header.ToArray()

	ok := syncGenesisHeader(ctx, pubKeys, users, genesisBlockHeader)
	if !ok {
		return false
	}
	waitForBlock(ctx)
	return true
}

type RegisterSideChainParam struct {
	Path         string
	Chainid      uint64
	Name         string
	BlocksToWait uint64
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

	user, ok := getAccountByPassword(ctx, registerSideChainParam.Path)
	if !ok {
		return false
	}
	ok = registerSideChain(ctx, user, registerSideChainParam.Chainid, registerSideChainParam.Name, registerSideChainParam.BlocksToWait)
	if !ok {
		return false
	}
	waitForBlock(ctx)
	return true
}

type ApproveRegisterSideChainParam struct {
	Path    []string
	Chainid uint64
}

func ApproveRegisterSideChain(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./side_chain_params/ApproveRegisterSideChain.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	approveRegisterSideChainParam := new(ApproveRegisterSideChainParam)
	err = json.Unmarshal(data, approveRegisterSideChainParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}

	var users []*sdk.Account
	var pubKeys []keypair.PublicKey
	time.Sleep(1 * time.Second)
	for _, path := range approveRegisterSideChainParam.Path {
		user, ok := getAccountByPassword(ctx, path)
		if !ok {
			return false
		}
		users = append(users, user)
		pubKeys = append(pubKeys, user.PublicKey)
	}

	ok := approveRegisterSideChain(ctx, pubKeys, users, approveRegisterSideChainParam.Chainid)
	if !ok {
		return false
	}
	waitForBlock(ctx)
	return true
}

type AssetMappingParam struct {
	Path      string
	AssetName string
	AssetList []*side_chain_manager.Asset
}

func AssetMapping(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./side_chain_params/AssetMapping.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	assetMappingParam := new(AssetMappingParam)
	err = json.Unmarshal(data, assetMappingParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}

	user, ok := getAccountByPassword(ctx, assetMappingParam.Path)
	if !ok {
		return false
	}
	ok = assetMapping(ctx, user, assetMappingParam.AssetName, assetMappingParam.AssetList)
	if !ok {
		return false
	}
	waitForBlock(ctx)
	return true
}
