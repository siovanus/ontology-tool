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
	"io/ioutil"
	"time"

	sdk "github.com/ontio/multi-chain-go-sdk"
	"github.com/ontio/multi-chain/native/service/side_chain_manager"
	osdk "github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology-tool/testframework"
)

type SyncGenesisHeaderParam struct {
	Path     []string
	ChainID  uint64
	ChainRpc string
}

func SyncGenesisHeader(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./side_chain_params/SyncGenesisHeader.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	syncGenesisHeaderParam := new(SyncGenesisHeaderParam)
	err = json.Unmarshal(data, syncGenesisHeaderParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}

	var users []*sdk.Account
	time.Sleep(1 * time.Second)
	for _, path := range syncGenesisHeaderParam.Path {
		user, ok := getAccountByPassword(ctx, path)
		if !ok {
			return false
		}
		users = append(users, user)
	}

	sideSdk := osdk.NewOntologySdk()
	sideSdk.NewRpcClient().SetAddress(syncGenesisHeaderParam.ChainRpc)
	genesisBlock, err := sideSdk.GetBlockByHeight(0)
	if err != nil {
		ctx.LogError("get side chain genesis block error: %s", err)
		return false
	}
	genesisBlockHeader := genesisBlock.Header.ToArray()

	txHash, err := ctx.Ont.Native.Hs.SyncGenesisHeader(syncGenesisHeaderParam.ChainID, genesisBlockHeader, users)
	if err != nil {
		ctx.LogError("ctx.Ont.Native.Scm.RegisterSideChain error: %v", err)
		return false
	}
	ctx.LogInfo("RegisterSideChain txHash is: %v", txHash.ToHexString())
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
	txHash, err := ctx.Ont.Native.Scm.RegisterSideChain(user.Address.ToBase58(), registerSideChainParam.Chainid, registerSideChainParam.Name, registerSideChainParam.BlocksToWait, user)
	if err != nil {
		ctx.LogError("ctx.Ont.Native.Scm.RegisterSideChain error: %v", err)
		return false
	}
	ctx.LogInfo("RegisterSideChain txHash is: %v", txHash.ToHexString())
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
	time.Sleep(1 * time.Second)
	for _, path := range approveRegisterSideChainParam.Path {
		user, ok := getAccountByPassword(ctx, path)
		if !ok {
			return false
		}
		users = append(users, user)
	}

	txHash, err := ctx.Ont.Native.Scm.ApproveRegisterSideChain(approveRegisterSideChainParam.Chainid, users)
	if err != nil {
		ctx.LogError("ctx.Ont.Native.Scm.ApproveRegisterSideChain error: %v", err)
		return false
	}
	ctx.LogInfo("ApproveRegisterSideChain txHash is: %v", txHash.ToHexString())
	waitForBlock(ctx)
	return true
}

type AssetMappingParam struct {
	Path      string
	AssetName string
	AssetList []*side_chain_manager.CrossChainContract
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
	txHash, err := ctx.Ont.Native.Scm.AssetMapping(user.Address.ToBase58(), assetMappingParam.AssetName, assetMappingParam.AssetList, user)
	if err != nil {
		ctx.LogError("ctx.Ont.Native.Scm.ApproveRegisterSideChain error: %v", err)
		return false
	}
	ctx.LogInfo("AssetMapping txHash is: %v", txHash.ToHexString())
	waitForBlock(ctx)
	return true
}

type ApproveAssetMappingParam struct {
	Path      []string
	AssetName string
}

func ApproveAssetMapping(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./side_chain_params/ApproveAssetMapping.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	approveAssetMappingParam := new(ApproveAssetMappingParam)
	err = json.Unmarshal(data, approveAssetMappingParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}

	var users []*sdk.Account
	time.Sleep(1 * time.Second)
	for _, path := range approveAssetMappingParam.Path {
		user, ok := getAccountByPassword(ctx, path)
		if !ok {
			return false
		}
		users = append(users, user)
	}

	txHash, err := ctx.Ont.Native.Scm.ApproveAssetMapping(approveAssetMappingParam.AssetName, users)
	if err != nil {
		ctx.LogError("ctx.Ont.Native.Scm.ApproveAssetMapping error: %v", err)
		return false
	}
	ctx.LogInfo("ApproveAssetMapping txHash is: %v", txHash.ToHexString())
	waitForBlock(ctx)
	return true
}

type InitRedeemScriptParam struct {
	Path         []string
	RedeemScript string
}

func InitRedeemScript(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./side_chain_params/InitRedeemScript.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	initRedeemScriptParam := new(InitRedeemScriptParam)
	err = json.Unmarshal(data, initRedeemScriptParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}

	var users []*sdk.Account
	time.Sleep(1 * time.Second)
	for _, path := range initRedeemScriptParam.Path {
		user, ok := getAccountByPassword(ctx, path)
		if !ok {
			return false
		}
		users = append(users, user)
	}

	txHash, err := ctx.Ont.Native.Ccm.InitRedeemScript(initRedeemScriptParam.RedeemScript, users)
	if err != nil {
		ctx.LogError("ctx.Ont.Native.Scm.ApproveAssetMapping error: %v", err)
		return false
	}
	ctx.LogInfo("ApproveAssetMapping txHash is: %v", txHash.ToHexString())
	waitForBlock(ctx)
	return true
}
