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

package governance_feeSplit

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"time"

	"github.com/ontio/ontology-crypto/keypair"
	"github.com/ontio/ontology-crypto/vrf"
	sdk "github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology-tool/testframework"
	"github.com/ontio/ontology/consensus/vbft/config"
	"github.com/ontio/ontology/core/types"
	"github.com/ontio/ontology/smartcontract/service/native/utils"
)

type Account struct {
	Path string
}

func RegIdWithPublicKey(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./params/RegIdWithPublicKey.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	account := new(Account)
	err = json.Unmarshal(data, account)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}
	time.Sleep(1 * time.Second)
	user, ok := getAccountByPassword(ctx, account.Path)
	if !ok {
		return false
	}
	ok = regIdWithPublicKey(ctx, user)
	if !ok {
		return false
	}
	waitForBlock(ctx)
	return true
}

func AssignFuncsToRole(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./params/AssignFuncsToRole.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	account := new(Account)
	err = json.Unmarshal(data, account)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}
	time.Sleep(1 * time.Second)
	user, ok := getAccountByPassword(ctx, account.Path)
	if !ok {
		return false
	}
	ok = assignFuncsToRole(ctx, user, utils.GovernanceContractAddress, "TrionesCandidatePeerOwner", "registerCandidate")
	if !ok {
		return false
	}
	waitForBlock(ctx)
	return true
}

type AssignFuncsToRoleAnyParam struct {
	Path            string
	ContractAddress string
	Role            string
	Function        string
}

func AssignFuncsToRoleAny(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./params/AssignFuncsToRoleAny.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	assignFuncsToRoleAnyParam := new(AssignFuncsToRoleAnyParam)
	err = json.Unmarshal(data, assignFuncsToRoleAnyParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}
	time.Sleep(1 * time.Second)
	user, ok := getAccountByPassword(ctx, assignFuncsToRoleAnyParam.Path)
	if !ok {
		return false
	}
	contractAddress, err := getAddressByHexString(assignFuncsToRoleAnyParam.ContractAddress)
	if err != nil {
		ctx.LogError("getAddressByHexString failed %v", err)
		return false
	}
	ok = assignFuncsToRole(ctx, user, contractAddress, assignFuncsToRoleAnyParam.Role, assignFuncsToRoleAnyParam.Function)
	if !ok {
		return false
	}
	waitForBlock(ctx)
	return true
}

type AssignOntIDsToRoleParam struct {
	Path1 string
	Ontid []string
}

func AssignOntIDsToRole(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./params/AssignOntIDsToRole.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	assignOntIDsToRoleParam := new(AssignOntIDsToRoleParam)
	err = json.Unmarshal(data, assignOntIDsToRoleParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}
	time.Sleep(1 * time.Second)
	user1, ok := getAccountByPassword(ctx, assignOntIDsToRoleParam.Path1)
	if !ok {
		return false
	}
	ok = assignOntIDsToRole(ctx, user1, utils.GovernanceContractAddress, "TrionesCandidatePeerOwner", assignOntIDsToRoleParam.Ontid)
	if !ok {
		return false
	}
	waitForBlock(ctx)
	return true
}

type AssignOntIDsToRoleAnyParam struct {
	Path1           string
	ContractAddress string
	Role            string
	Ontid           []string
}

func AssignOntIDsToRoleAny(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./params/AssignOntIDsToRoleAny.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	assignOntIDsToRoleAnyParam := new(AssignOntIDsToRoleAnyParam)
	err = json.Unmarshal(data, assignOntIDsToRoleAnyParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}
	time.Sleep(1 * time.Second)
	user1, ok := getAccountByPassword(ctx, assignOntIDsToRoleAnyParam.Path1)
	if !ok {
		return false
	}
	contractAddress, err := getAddressByHexString(assignOntIDsToRoleAnyParam.ContractAddress)
	if err != nil {
		ctx.LogError("getAddressByHexString failed %v", err)
		return false
	}
	ok = assignOntIDsToRole(ctx, user1, contractAddress, assignOntIDsToRoleAnyParam.Role, assignOntIDsToRoleAnyParam.Ontid)
	if !ok {
		return false
	}
	waitForBlock(ctx)
	return true
}

type VrfParam struct {
	Path string
}

type vrfData struct {
	BlockNum uint32 `json:"block_num"`
	PrevVrf  []byte `json:"prev_vrf"`
}

func Vrf(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./params/Vrf.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	vrfParam := new(VrfParam)
	err = json.Unmarshal(data, vrfParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}
	user, ok := getAccount(ctx, vrfParam.Path)
	if !ok {
		return false
	}

	data, err = json.Marshal(&vrfData{
		BlockNum: 0,
		PrevVrf:  keypair.SerializePublicKey(user.PublicKey),
	})
	if err != nil {
		ctx.LogError("json.Unmarshal vrf payload failed %v", err)
		return false
	}

	value, proof, err := vrf.Vrf(user.PrivateKey, data)
	if err != nil {
		ctx.LogError("vrf computation failed %v", err)
		return false
	}

	if ok, err := vrf.Verify(user.PublicKey, data, value, proof); err != nil || !ok {
		ctx.LogError("vrf verify failed: %v", err)
		return false
	}

	ctx.LogInfo("vrf value: %s", hex.EncodeToString(value))
	ctx.LogInfo("vrf proof: %s", hex.EncodeToString(proof))

	return true
}

type TransferMultiSignParam struct {
	Path1  []string
	Path2  []string
	Amount []uint64
}

func TransferOngMultiSign(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./params/TransferOngMultiSign.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	transferMultiSignParam := new(TransferMultiSignParam)
	err = json.Unmarshal(data, transferMultiSignParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}
	var users []*sdk.Account
	var pubKeys []keypair.PublicKey
	time.Sleep(1 * time.Second)
	for _, path := range transferMultiSignParam.Path1 {
		user, ok := getAccountByPassword(ctx, path)
		if !ok {
			return false
		}
		users = append(users, user)
		pubKeys = append(pubKeys, user.PublicKey)
	}
	time.Sleep(1 * time.Second)
	for index, path2 := range transferMultiSignParam.Path2 {
		user2, ok := getAccountByPassword(ctx, path2)
		if !ok {
			return false
		}
		ok = transferOngMultiSign(ctx, pubKeys, users, user2.Address, transferMultiSignParam.Amount[index])
		if !ok {
			return false
		}
	}
	waitForBlock(ctx)
	return true
}

type GetAddressMultiSignParam struct {
	PubKeys []string
}

func GetAddressMultiSign(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./params/GetAddressMultiSign.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	getAddressMultiSignParam := new(GetAddressMultiSignParam)
	err = json.Unmarshal(data, getAddressMultiSignParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}
	var pubKeys []keypair.PublicKey
	time.Sleep(1 * time.Second)
	for _, v := range getAddressMultiSignParam.PubKeys {
		vByte, err := hex.DecodeString(v)
		if err != nil {
			ctx.LogError("hex.DecodeString failed %v", err)
		}
		k, err := keypair.DeserializePublicKey(vByte)
		if err != nil {
			ctx.LogError("keypair.DeserializePublicKey failed %v", err)
		}
		pubKeys = append(pubKeys, k)
	}
	from, err := types.AddressFromMultiPubKeys(pubKeys, int((5*len(pubKeys)+6)/7))
	if err != nil {
		ctx.LogError("types.AddressFromMultiPubKeys error", err)
	}
	fmt.Println("address is:", from.ToBase58())
	return true
}

func GetVbftInfo(ctx *testframework.TestFrameworkContext) bool {
	blkNum, err := ctx.Ont.GetCurrentBlockHeight()
	if err != nil {
		ctx.LogError("TestGetVbftInfo GetBlockCount error:%s", err)
		return false
	}
	blk, err := ctx.Ont.GetBlockByHeight(blkNum - 1)
	if err != nil {
		ctx.LogError("TestGetVbftInfo GetBlockByHeight error:%s", err)
		return false
	}
	block, err := initVbftBlock(blk)
	if err != nil {
		ctx.LogError("TestGetVbftInfo initVbftBlock error:%s", err)
		return false
	}

	var cfg vconfig.ChainConfig
	if block.Info.NewChainConfig != nil {
		cfg = *block.Info.NewChainConfig
	} else {
		var cfgBlock *types.Block
		if block.Info.LastConfigBlockNum != math.MaxUint32 {
			cfgBlock, err = ctx.Ont.GetBlockByHeight(block.Info.LastConfigBlockNum)
			if err != nil {
				ctx.LogError("TestGetVbftInfo chainconfig GetBlockByHeight error:%s", err)
				return false
			}
		}
		blk, err := initVbftBlock(cfgBlock)
		if err != nil {
			ctx.LogError("TestGetVbftInfo initVbftBlock error:%s", err)
			return false
		}
		if blk.Info.NewChainConfig == nil {
			ctx.LogError("TestGetVbftInfo newchainconfig error:%s", err)
			return false
		}
		cfg = *blk.Info.NewChainConfig
	}
	fmt.Printf("block vbft chainConfig, View:%d, N:%d, C:%d, BlockMsgDelay:%v, HashMsgDelay:%v, PeerHandshakeTimeout:%v, MaxBlockChangeView:%d, PosTable:%v\n",
		cfg.View, cfg.N, cfg.C, cfg.BlockMsgDelay, cfg.HashMsgDelay, cfg.PeerHandshakeTimeout, cfg.MaxBlockChangeView, cfg.PosTable)
	for _, p := range cfg.Peers {
		fmt.Printf("peerInfo Index: %d, ID:%v\n", p.Index, p.ID)
	}
	return true
}

type MultiTransferParam struct {
	FromPath  []string
	ToAddress []string
	Amount    []uint64
}

func MultiTransferOng(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./params/MultiTransferOng.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	multiTransferParam := new(MultiTransferParam)
	err = json.Unmarshal(data, multiTransferParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}
	var users []*sdk.Account
	time.Sleep(1 * time.Second)
	for _, path := range multiTransferParam.FromPath {
		user, ok := getAccountByPassword(ctx, path)
		if !ok {
			return false
		}
		users = append(users, user)
	}
	ok := multiTransfer(ctx, utils.OngContractAddress, users, multiTransferParam.ToAddress, multiTransferParam.Amount)
	if !ok {
		return false
	}
	waitForBlock(ctx)
	return true
}
