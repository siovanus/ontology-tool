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
	"github.com/ontio/ontology/common"
	"io/ioutil"
	"time"

	sdk "github.com/ontio/multi-chain-go-sdk"
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
	Router       uint64
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
	txHash, err := ctx.Ont.Native.Scm.RegisterSideChain(user.Address.ToBase58(), registerSideChainParam.Chainid,
		registerSideChainParam.Router, registerSideChainParam.Name, registerSideChainParam.BlocksToWait, user)
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
		ctx.LogError("ctx.Ont.Native.Ccm.InitRedeemScript error: %v", err)
		return false
	}
	ctx.LogInfo("InitRedeemScript txHash is: %v", txHash.ToHexString())
	waitForBlock(ctx)
	return true
}

type RegisterPeerParam struct {
	PeerPubkey string
	Path       string
}

func RegisterCandidate(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./side_chain_params/RegisterCandidate.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	registerPeerParam := new(RegisterPeerParam)
	err = json.Unmarshal(data, registerPeerParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}

	user, ok := getAccountByPassword(ctx, registerPeerParam.Path)
	if !ok {
		return false
	}
	txHash, err := ctx.Ont.Native.Nm.RegisterCandidate(registerPeerParam.PeerPubkey, user.Address[:], user)
	if err != nil {
		ctx.LogError("ctx.Ont.Native.Nm.RegisterCandidate error: %v", err)
		return false
	}
	ctx.LogInfo("RegisterCandidate txHash is: %v", txHash.ToHexString())
	waitForBlock(ctx)
	return true
}

func UnRegisterCandidate(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./side_chain_params/UnRegisterCandidate.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	registerPeerParam := new(RegisterPeerParam)
	err = json.Unmarshal(data, registerPeerParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}

	user, ok := getAccountByPassword(ctx, registerPeerParam.Path)
	if !ok {
		return false
	}
	txHash, err := ctx.Ont.Native.Nm.UnRegisterCandidate(registerPeerParam.PeerPubkey, user.Address[:], user)
	if err != nil {
		ctx.LogError("ctx.Ont.Native.Nm.UnRegisterCandidate error: %v", err)
		return false
	}
	ctx.LogInfo("UnRegisterCandidate txHash is: %v", txHash.ToHexString())
	waitForBlock(ctx)
	return true
}

func QuitNode(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./side_chain_params/QuitNode.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	registerPeerParam := new(RegisterPeerParam)
	err = json.Unmarshal(data, registerPeerParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}

	user, ok := getAccountByPassword(ctx, registerPeerParam.Path)
	if !ok {
		return false
	}
	txHash, err := ctx.Ont.Native.Nm.QuitNode(registerPeerParam.PeerPubkey, user.Address[:], user)
	if err != nil {
		ctx.LogError("ctx.Ont.Native.Nm.QuitNode error: %v", err)
		return false
	}
	ctx.LogInfo("QuitNode txHash is: %v", txHash.ToHexString())
	waitForBlock(ctx)
	return true
}

type PeerParam struct {
	PeerPubkey string
	Path       []string
}

func ApproveCandidate(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./side_chain_params/ApproveCandidate.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	peerParam := new(PeerParam)
	err = json.Unmarshal(data, peerParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}

	var users []*sdk.Account
	time.Sleep(1 * time.Second)
	for _, path := range peerParam.Path {
		user, ok := getAccountByPassword(ctx, path)
		if !ok {
			return false
		}
		users = append(users, user)
	}

	txHash, err := ctx.Ont.Native.Nm.ApproveCandidate(peerParam.PeerPubkey, users)
	if err != nil {
		ctx.LogError("ctx.Ont.Native.Nm.ApproveCandidate error: %v", err)
		return false
	}
	ctx.LogInfo("ApproveCandidate txHash is: %v", txHash.ToHexString())
	waitForBlock(ctx)
	return true
}

func RejectCandidate(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./side_chain_params/RejectCandidate.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	peerParam := new(PeerParam)
	err = json.Unmarshal(data, peerParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}

	var users []*sdk.Account
	time.Sleep(1 * time.Second)
	for _, path := range peerParam.Path {
		user, ok := getAccountByPassword(ctx, path)
		if !ok {
			return false
		}
		users = append(users, user)
	}

	txHash, err := ctx.Ont.Native.Nm.RejectCandidate(peerParam.PeerPubkey, users)
	if err != nil {
		ctx.LogError("ctx.Ont.Native.Nm.RejectCandidate error: %v", err)
		return false
	}
	ctx.LogInfo("RejectCandidate txHash is: %v", txHash.ToHexString())
	waitForBlock(ctx)
	return true
}

type PeerListParam struct {
	PeerPubkeyList []string
	Path           []string
}

func BlackNode(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./side_chain_params/BlackNode.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	peerListParam := new(PeerListParam)
	err = json.Unmarshal(data, peerListParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}

	var users []*sdk.Account
	time.Sleep(1 * time.Second)
	for _, path := range peerListParam.Path {
		user, ok := getAccountByPassword(ctx, path)
		if !ok {
			return false
		}
		users = append(users, user)
	}

	txHash, err := ctx.Ont.Native.Nm.BlackNode(peerListParam.PeerPubkeyList, users)
	if err != nil {
		ctx.LogError("ctx.Ont.Native.Nm.BlackNode error: %v", err)
		return false
	}
	ctx.LogInfo("BlackNode txHash is: %v", txHash.ToHexString())
	waitForBlock(ctx)
	return true
}

func WhiteNode(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./side_chain_params/WhiteNode.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	peerParam := new(PeerParam)
	err = json.Unmarshal(data, peerParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}

	var users []*sdk.Account
	time.Sleep(1 * time.Second)
	for _, path := range peerParam.Path {
		user, ok := getAccountByPassword(ctx, path)
		if !ok {
			return false
		}
		users = append(users, user)
	}

	txHash, err := ctx.Ont.Native.Nm.WhiteNode(peerParam.PeerPubkey, users)
	if err != nil {
		ctx.LogError("ctx.Ont.Native.Nm.WhiteNode error: %v", err)
		return false
	}
	ctx.LogInfo("WhiteNode txHash is: %v", txHash.ToHexString())
	waitForBlock(ctx)
	return true
}

type Configuration struct {
	BlockMsgDelay        uint32
	HashMsgDelay         uint32
	PeerHandshakeTimeout uint32
	Path                 []string
}

func UpdateConfig(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./side_chain_params/UpdateConfig.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	configuration := new(Configuration)
	err = json.Unmarshal(data, configuration)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}

	var users []*sdk.Account
	time.Sleep(1 * time.Second)
	for _, path := range configuration.Path {
		user, ok := getAccountByPassword(ctx, path)
		if !ok {
			return false
		}
		users = append(users, user)
	}

	txHash, err := ctx.Ont.Native.Nm.UpdateConfig(configuration.BlockMsgDelay, configuration.HashMsgDelay,
		configuration.PeerHandshakeTimeout, users)
	if err != nil {
		ctx.LogError("ctx.Ont.Native.Nm.UpdateConfig error: %v", err)
		return false
	}
	ctx.LogInfo("UpdateConfig txHash is: %v", txHash.ToHexString())
	waitForBlock(ctx)
	return true
}

type RelayerParam struct {
	Path string
}

func RegisterRelayer(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./side_chain_params/RegisterRelayer.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	relayerParam := new(RelayerParam)
	err = json.Unmarshal(data, relayerParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}

	user, ok := getAccountByPassword(ctx, relayerParam.Path)
	if !ok {
		return false
	}
	txHash, err := ctx.Ont.Native.Rm.RegisterRelayer(user.Address[:], user)
	if err != nil {
		ctx.LogError("ctx.Ont.Native.Rm.RegisterRelayer error: %v", err)
		return false
	}
	ctx.LogInfo("RegisterRelayer txHash is: %v", txHash.ToHexString())
	waitForBlock(ctx)
	return true
}

func UnRegisterRelayer(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./side_chain_params/UnRegisterRelayer.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	relayerParam := new(RelayerParam)
	err = json.Unmarshal(data, relayerParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}

	user, ok := getAccountByPassword(ctx, relayerParam.Path)
	if !ok {
		return false
	}
	txHash, err := ctx.Ont.Native.Rm.UnRegisterRelayer(user.Address[:], user)
	if err != nil {
		ctx.LogError("ctx.Ont.Native.Rm.UnRegisterRelayer error: %v", err)
		return false
	}
	ctx.LogInfo("UnRegisterRelayer txHash is: %v", txHash.ToHexString())
	waitForBlock(ctx)
	return true
}

func QuitRelayer(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./side_chain_params/QuitRelayer.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	relayerParam := new(RelayerParam)
	err = json.Unmarshal(data, relayerParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}

	user, ok := getAccountByPassword(ctx, relayerParam.Path)
	if !ok {
		return false
	}
	txHash, err := ctx.Ont.Native.Rm.QuitRelayer(user.Address[:], user)
	if err != nil {
		ctx.LogError("ctx.Ont.Native.Rm.QuitRelayer error: %v", err)
		return false
	}
	ctx.LogInfo("QuitRelayer txHash is: %v", txHash.ToHexString())
	waitForBlock(ctx)
	return true
}

type RelayerParam2 struct {
	Address string
	Path    []string
}

func ApproveRelayer(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./side_chain_params/ApproveRelayer.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	relayerParam2 := new(RelayerParam2)
	err = json.Unmarshal(data, relayerParam2)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}

	var users []*sdk.Account
	time.Sleep(1 * time.Second)
	for _, path := range relayerParam2.Path {
		user, ok := getAccountByPassword(ctx, path)
		if !ok {
			return false
		}
		users = append(users, user)
	}

	address, err := common.AddressFromBase58(relayerParam2.Address)
	if err != nil {
		ctx.LogError("common.AddressFromBase58 failed %v", err)
		return false
	}
	txHash, err := ctx.Ont.Native.Rm.ApproveRelayer(address[:], users)
	if err != nil {
		ctx.LogError("ctx.Ont.Native.Rm.ApproveRelayer error: %v", err)
		return false
	}
	ctx.LogInfo("ApproveRelayer txHash is: %v", txHash.ToHexString())
	waitForBlock(ctx)
	return true
}

func RejectRelayer(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./side_chain_params/RejectRelayer.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	relayerParam2 := new(RelayerParam2)
	err = json.Unmarshal(data, relayerParam2)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}

	var users []*sdk.Account
	time.Sleep(1 * time.Second)
	for _, path := range relayerParam2.Path {
		user, ok := getAccountByPassword(ctx, path)
		if !ok {
			return false
		}
		users = append(users, user)
	}

	address, err := common.AddressFromBase58(relayerParam2.Address)
	if err != nil {
		ctx.LogError("common.AddressFromBase58 failed %v", err)
		return false
	}
	txHash, err := ctx.Ont.Native.Rm.RejectRelayer(address[:], users)
	if err != nil {
		ctx.LogError("ctx.Ont.Native.Rm.RejectRelayer error: %v", err)
		return false
	}
	ctx.LogInfo("RejectRelayer txHash is: %v", txHash.ToHexString())
	waitForBlock(ctx)
	return true
}

type RelayerListParam struct {
	AddressList []string
	Path        []string
}

func BlackRelayer(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./side_chain_params/BlackRelayer.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	relayerListParam := new(RelayerListParam)
	err = json.Unmarshal(data, relayerListParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}

	var users []*sdk.Account
	time.Sleep(1 * time.Second)
	for _, path := range relayerListParam.Path {
		user, ok := getAccountByPassword(ctx, path)
		if !ok {
			return false
		}
		users = append(users, user)
	}
	addresses := make([][]byte, 0)
	for _, address := range relayerListParam.AddressList {
		address, err := common.AddressFromBase58(address)
		if err != nil {
			ctx.LogError("common.AddressFromBase58 failed %v", err)
			return false
		}
		addresses = append(addresses, address[:])
	}

	txHash, err := ctx.Ont.Native.Rm.BlackRelayer(addresses, users)
	if err != nil {
		ctx.LogError("ctx.Ont.Native.Rm.BlackRelayer error: %v", err)
		return false
	}
	ctx.LogInfo("BlackRelayer txHash is: %v", txHash.ToHexString())
	waitForBlock(ctx)
	return true
}

func WhiteRelayer(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./side_chain_params/WhiteRelayer.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	relayerParam2 := new(RelayerParam2)
	err = json.Unmarshal(data, relayerParam2)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}

	var users []*sdk.Account
	time.Sleep(1 * time.Second)
	for _, path := range relayerParam2.Path {
		user, ok := getAccountByPassword(ctx, path)
		if !ok {
			return false
		}
		users = append(users, user)
	}

	address, err := common.AddressFromBase58(relayerParam2.Address)
	if err != nil {
		ctx.LogError("common.AddressFromBase58 failed %v", err)
		return false
	}
	txHash, err := ctx.Ont.Native.Rm.WhiteRelayer(address[:], users)
	if err != nil {
		ctx.LogError("ctx.Ont.Native.Rm.WhiteRelayer error: %v", err)
		return false
	}
	ctx.LogInfo("WhiteRelayer txHash is: %v", txHash.ToHexString())
	waitForBlock(ctx)
	return true
}
