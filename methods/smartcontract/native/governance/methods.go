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

package governance

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"time"

	log4 "github.com/alecthomas/log4go"
	"github.com/ontio/ontology-crypto/keypair"
	s "github.com/ontio/ontology-crypto/signature"
	"github.com/ontio/ontology-crypto/vrf"
	sdk "github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology-tool/common"
	ocommon "github.com/ontio/ontology/common"
	"github.com/ontio/ontology/common/password"
	"github.com/ontio/ontology/consensus/vbft/config"
	"github.com/ontio/ontology/core/types"
	"github.com/ontio/ontology/smartcontract/service/native/governance"
	"github.com/ontio/ontology/smartcontract/service/native/utils"
)

type Account struct {
	Path string
}

func RegIdWithPublicKey(ontSdk *sdk.OntologySdk) bool {
	data, err := ioutil.ReadFile("./params/RegIdWithPublicKey.json")
	if err != nil {
		log4.Error("ioutil.ReadFile failed ", err)
		return false
	}
	account := new(Account)
	err = json.Unmarshal(data, account)
	if err != nil {
		log4.Error("json.Unmarshal failed ", err)
		return false
	}
	time.Sleep(1 * time.Second)
	user, ok := common.GetAccountByPassword(ontSdk, account.Path)
	if !ok {
		return false
	}
	ok = regIdWithPublicKey(ontSdk, user)
	if !ok {
		return false
	}
	common.WaitForBlock(ontSdk)
	return true
}

func AssignFuncsToRole(ontSdk *sdk.OntologySdk) bool {
	data, err := ioutil.ReadFile("./params/AssignFuncsToRole.json")
	if err != nil {
		log4.Error("ioutil.ReadFile failed ", err)
		return false
	}
	account := new(Account)
	err = json.Unmarshal(data, account)
	if err != nil {
		log4.Error("json.Unmarshal failed ", err)
		return false
	}
	time.Sleep(1 * time.Second)
	user, ok := common.GetAccountByPassword(ontSdk, account.Path)
	if !ok {
		return false
	}
	ok = assignFuncsToRole(ontSdk, user, utils.GovernanceContractAddress, "TrionesCandidatePeerOwner", "registerCandidate")
	if !ok {
		return false
	}
	common.WaitForBlock(ontSdk)
	return true
}

type AssignFuncsToRoleAnyParam struct {
	Path            string
	ContractAddress string
	Role            string
	Function        string
}

func AssignFuncsToRoleAny(ontSdk *sdk.OntologySdk) bool {
	data, err := ioutil.ReadFile("./params/AssignFuncsToRoleAny.json")
	if err != nil {
		log4.Error("ioutil.ReadFile failed ", err)
		return false
	}
	assignFuncsToRoleAnyParam := new(AssignFuncsToRoleAnyParam)
	err = json.Unmarshal(data, assignFuncsToRoleAnyParam)
	if err != nil {
		log4.Error("json.Unmarshal failed ", err)
		return false
	}
	time.Sleep(1 * time.Second)
	user, ok := common.GetAccountByPassword(ontSdk, assignFuncsToRoleAnyParam.Path)
	if !ok {
		return false
	}
	contractAddress, err := common.GetAddressByHexString(assignFuncsToRoleAnyParam.ContractAddress)
	if err != nil {
		log4.Error("getAddressByHexString failed ", err)
		return false
	}
	ok = assignFuncsToRole(ontSdk, user, contractAddress, assignFuncsToRoleAnyParam.Role, assignFuncsToRoleAnyParam.Function)
	if !ok {
		return false
	}
	common.WaitForBlock(ontSdk)
	return true
}

type AssignOntIDsToRoleParam struct {
	Path1 string
	Ontid []string
}

func AssignOntIDsToRole(ontSdk *sdk.OntologySdk) bool {
	data, err := ioutil.ReadFile("./params/AssignOntIDsToRole.json")
	if err != nil {
		log4.Error("ioutil.ReadFile failed ", err)
		return false
	}
	assignOntIDsToRoleParam := new(AssignOntIDsToRoleParam)
	err = json.Unmarshal(data, assignOntIDsToRoleParam)
	if err != nil {
		log4.Error("json.Unmarshal failed ", err)
		return false
	}
	time.Sleep(1 * time.Second)
	user1, ok := common.GetAccountByPassword(ontSdk, assignOntIDsToRoleParam.Path1)
	if !ok {
		return false
	}
	ok = assignOntIDsToRole(ontSdk, user1, utils.GovernanceContractAddress, "TrionesCandidatePeerOwner", assignOntIDsToRoleParam.Ontid)
	if !ok {
		return false
	}
	common.WaitForBlock(ontSdk)
	return true
}

type AssignOntIDsToRoleAnyParam struct {
	Path1           string
	ContractAddress string
	Role            string
	Ontid           []string
}

func AssignOntIDsToRoleAny(ontSdk *sdk.OntologySdk) bool {
	data, err := ioutil.ReadFile("./params/AssignOntIDsToRoleAny.json")
	if err != nil {
		log4.Error("ioutil.ReadFile failed ", err)
		return false
	}
	assignOntIDsToRoleAnyParam := new(AssignOntIDsToRoleAnyParam)
	err = json.Unmarshal(data, assignOntIDsToRoleAnyParam)
	if err != nil {
		log4.Error("json.Unmarshal failed ", err)
		return false
	}
	time.Sleep(1 * time.Second)
	user1, ok := common.GetAccountByPassword(ontSdk, assignOntIDsToRoleAnyParam.Path1)
	if !ok {
		return false
	}
	contractAddress, err := common.GetAddressByHexString(assignOntIDsToRoleAnyParam.ContractAddress)
	if err != nil {
		log4.Error("getAddressByHexString failed ", err)
		return false
	}
	ok = assignOntIDsToRole(ontSdk, user1, contractAddress, assignOntIDsToRoleAnyParam.Role, assignOntIDsToRoleAnyParam.Ontid)
	if !ok {
		return false
	}
	common.WaitForBlock(ontSdk)
	return true
}

type RegisterCandidateParam struct {
	Path       []string
	PeerPubkey []string
	InitPos    []uint32
}

func RegisterCandidate(ontSdk *sdk.OntologySdk) bool {
	data, err := ioutil.ReadFile("./params/RegisterCandidate.json")
	if err != nil {
		log4.Error("ioutil.ReadFile failed ", err)
		return false
	}
	registerCandidateParam := new(RegisterCandidateParam)
	err = json.Unmarshal(data, registerCandidateParam)
	if err != nil {
		log4.Error("json.Unmarshal failed ", err)
		return false
	}
	time.Sleep(1 * time.Second)
	for i := 0; i < len(registerCandidateParam.PeerPubkey); i++ {
		user, ok := common.GetAccountByPassword(ontSdk, registerCandidateParam.Path[i])
		if !ok {
			return false
		}
		ok = registerCandidate(ontSdk, user, registerCandidateParam.PeerPubkey[i], registerCandidateParam.InitPos[i])
		if !ok {
			return false
		}
	}
	common.WaitForBlock(ontSdk)
	return true
}

type RegisterCandidate2SignParam struct {
	Key        string
	Address    string
	Salt       string
	Path       string
	PeerPubkey string
	InitPos    uint32
}

func RegisterCandidate2Sign(ontSdk *sdk.OntologySdk) bool {
	//"+UADcReBcLq0pn/2Grmz+UJsKl3ryop8pgRVHbQVgTBfT0lho06Svh4eQLSmC93j"
	//"AG9W6c7nNhaiywcyVPgW9hQKvUYQr5iLvk"
	//"IfxFV0Fer5LknIyCLP2P2w==2"

	data, err := ioutil.ReadFile("./params/RegisterCandidate2Sign.json")
	if err != nil {
		log4.Error("ioutil.ReadFile failed ", err)
		return false
	}
	registerCandidate2SignParam := new(RegisterCandidate2SignParam)
	err = json.Unmarshal(data, registerCandidate2SignParam)
	if err != nil {
		log4.Error("json.Unmarshal failed ", err)
		return false
	}

	key, _ := base64.StdEncoding.DecodeString(registerCandidate2SignParam.Key)
	salt, _ := base64.StdEncoding.DecodeString(registerCandidate2SignParam.Salt)
	var res = keypair.ProtectedKey{
		Alg:     "ECDSA",
		Address: registerCandidate2SignParam.Address,
		Key:     key,
		Salt:    salt,
		EncAlg:  "aes-256-gcm",
	}
	res.Param = make(map[string]string)
	res.Param["curve"] = "P-256"

	time.Sleep(1 * time.Second)
	pwd, err := password.GetPassword()
	if err != nil {
		log4.Error("getPassword error:%s", err)
		return false
	}
	pri, err := keypair.DecryptWithCustomScrypt(&res, pwd, &keypair.ScryptParam{
		N:     4096,
		R:     keypair.DEFAULT_R,
		P:     keypair.DEFAULT_P,
		DKLen: keypair.DEFAULT_DERIVED_KEY_LENGTH,
	})
	//pri, err := keypair.DecryptPrivateKey(&res, pwd)
	if err != nil {
		log4.Error("error: ", err)
		return false
	}
	address, _ := ocommon.AddressFromBase58(registerCandidate2SignParam.Address)
	account := &sdk.Account{
		PrivateKey: pri,
		PublicKey:  pri.Public(),
		Address:    address,
		SigScheme:  s.SHA256withECDSA,
	}
	user, ok := common.GetAccountByPassword(ontSdk, registerCandidate2SignParam.Path)
	if !ok {
		return false
	}
	ok = registerCandidate2Sign(ontSdk, account, user, registerCandidate2SignParam.PeerPubkey, registerCandidate2SignParam.InitPos)
	if !ok {
		return false
	}
	return true
}

type UnRegisterCandidateParam struct {
	Path       string
	PeerPubkey string
}

func UnRegisterCandidate(ontSdk *sdk.OntologySdk) bool {
	data, err := ioutil.ReadFile("./params/UnRegisterCandidate.json")
	if err != nil {
		log4.Error("ioutil.ReadFile failed ", err)
		return false
	}
	unRegisterCandidateParam := new(UnRegisterCandidateParam)
	err = json.Unmarshal(data, unRegisterCandidateParam)
	if err != nil {
		log4.Error("json.Unmarshal failed ", err)
		return false
	}
	user, ok := common.GetAccountByPassword(ontSdk, unRegisterCandidateParam.Path)
	if !ok {
		return false
	}
	ok = unRegisterCandidate(ontSdk, user, unRegisterCandidateParam.PeerPubkey)
	if !ok {
		return false
	}
	common.WaitForBlock(ontSdk)
	return true
}

type ApproveCandidateParam struct {
	Path       []string
	PeerPubkey []string
}

func ApproveCandidate(ontSdk *sdk.OntologySdk) bool {
	data, err := ioutil.ReadFile("./params/ApproveCandidate.json")
	if err != nil {
		log4.Error("ioutil.ReadFile failed ", err)
		return false
	}
	approveCandidateParam := new(ApproveCandidateParam)
	err = json.Unmarshal(data, approveCandidateParam)
	if err != nil {
		log4.Error("json.Unmarshal failed ", err)
		return false
	}
	var users []*sdk.Account
	var pubKeys []keypair.PublicKey
	time.Sleep(1 * time.Second)
	for _, path := range approveCandidateParam.Path {
		user, ok := common.GetAccountByPassword(ontSdk, path)
		if !ok {
			return false
		}
		users = append(users, user)
		pubKeys = append(pubKeys, user.PublicKey)
	}
	for _, peerPubkey := range approveCandidateParam.PeerPubkey {
		ok := approveCandidateMultiSign(ontSdk, pubKeys, users, peerPubkey)
		if !ok {
			return false
		}
	}
	common.WaitForBlock(ontSdk)
	return true
}

type RejectCandidateParam struct {
	Path       []string
	PeerPubkey string
}

func RejectCandidate(ontSdk *sdk.OntologySdk) bool {
	data, err := ioutil.ReadFile("./params/RejectCandidate.json")
	if err != nil {
		log4.Error("ioutil.ReadFile failed ", err)
		return false
	}
	rejectCandidateParam := new(RejectCandidateParam)
	err = json.Unmarshal(data, rejectCandidateParam)
	if err != nil {
		log4.Error("json.Unmarshal failed ", err)
		return false
	}
	var users []*sdk.Account
	var pubKeys []keypair.PublicKey
	time.Sleep(1 * time.Second)
	for _, path := range rejectCandidateParam.Path {
		user, ok := common.GetAccountByPassword(ontSdk, path)
		if !ok {
			return false
		}
		users = append(users, user)
		pubKeys = append(pubKeys, user.PublicKey)
	}
	ok := rejectCandidateMultiSign(ontSdk, pubKeys, users, rejectCandidateParam.PeerPubkey)
	if !ok {
		return false
	}
	common.WaitForBlock(ontSdk)
	return true
}

type ChangeMaxAuthorizationParam struct {
	PathList         []string
	PeerPubkeyList   []string
	MaxAuthorizeList []uint32
}

func ChangeMaxAuthorization(ontSdk *sdk.OntologySdk) bool {
	data, err := ioutil.ReadFile("./params/ChangeMaxAuthorization.json")
	if err != nil {
		log4.Error("ioutil.ReadFile failed ", err)
		return false
	}
	changeMaxAuthorizationParam := new(ChangeMaxAuthorizationParam)
	err = json.Unmarshal(data, changeMaxAuthorizationParam)
	if err != nil {
		log4.Error("json.Unmarshal failed ", err)
		return false
	}
	time.Sleep(1 * time.Second)
	for index, path := range changeMaxAuthorizationParam.PathList {
		user, ok := common.GetAccountByPassword(ontSdk, path)
		if !ok {
			return false
		}
		ok = changeMaxAuthorization(ontSdk, user, changeMaxAuthorizationParam.PeerPubkeyList[index], changeMaxAuthorizationParam.MaxAuthorizeList[index])
		if !ok {
			return false
		}
	}
	common.WaitForBlock(ontSdk)
	return true
}

type SetPeerCostParam struct {
	PathList       []string
	PeerPubkeyList []string
	PeerCostList   []uint32
}

func SetPeerCost(ontSdk *sdk.OntologySdk) bool {
	data, err := ioutil.ReadFile("./params/SetPeerCost.json")
	if err != nil {
		log4.Error("ioutil.ReadFile failed ", err)
		return false
	}
	setPeerCostParam := new(SetPeerCostParam)
	err = json.Unmarshal(data, setPeerCostParam)
	if err != nil {
		log4.Error("json.Unmarshal failed ", err)
		return false
	}
	time.Sleep(1 * time.Second)
	for index, path := range setPeerCostParam.PathList {
		user, ok := common.GetAccountByPassword(ontSdk, path)
		if !ok {
			return false
		}
		ok = setPeerCost(ontSdk, user, setPeerCostParam.PeerPubkeyList[index], setPeerCostParam.PeerCostList[index])
		if !ok {
			return false
		}
	}
	common.WaitForBlock(ontSdk)
	return true
}

type AddInitPosParam struct {
	Path       string
	PeerPubkey string
	Pos        uint32
}

func AddInitPos(ontSdk *sdk.OntologySdk) bool {
	data, err := ioutil.ReadFile("./params/AddInitPos.json")
	if err != nil {
		log4.Error("ioutil.ReadFile failed ", err)
		return false
	}
	addInitPosParam := new(AddInitPosParam)
	err = json.Unmarshal(data, addInitPosParam)
	if err != nil {
		log4.Error("json.Unmarshal failed ", err)
		return false
	}
	time.Sleep(1 * time.Second)
	user, ok := common.GetAccountByPassword(ontSdk, addInitPosParam.Path)
	if !ok {
		return false
	}
	ok = addInitPos(ontSdk, user, addInitPosParam.PeerPubkey, addInitPosParam.Pos)
	if !ok {
		return false
	}
	common.WaitForBlock(ontSdk)
	return true
}

type ReduceInitPosParam struct {
	Path       string
	PeerPubkey string
	Pos        uint32
}

func ReduceInitPos(ontSdk *sdk.OntologySdk) bool {
	data, err := ioutil.ReadFile("./params/ReduceInitPos.json")
	if err != nil {
		log4.Error("ioutil.ReadFile failed ", err)
		return false
	}
	reduceInitPosParam := new(ReduceInitPosParam)
	err = json.Unmarshal(data, reduceInitPosParam)
	if err != nil {
		log4.Error("json.Unmarshal failed ", err)
		return false
	}
	time.Sleep(1 * time.Second)
	user, ok := common.GetAccountByPassword(ontSdk, reduceInitPosParam.Path)
	if !ok {
		return false
	}
	ok = reduceInitPos(ontSdk, user, reduceInitPosParam.PeerPubkey, reduceInitPosParam.Pos)
	if !ok {
		return false
	}
	common.WaitForBlock(ontSdk)
	return true
}

type AuthorizeForPeerParam struct {
	Path           string
	PeerPubkeyList []string
	PosList        []uint32
}

func AuthorizeForPeer(ontSdk *sdk.OntologySdk) bool {
	data, err := ioutil.ReadFile("./params/AuthorizeForPeer.json")
	if err != nil {
		log4.Error("ioutil.ReadFile failed ", err)
		return false
	}
	authorizeForPeerParam := new(AuthorizeForPeerParam)
	err = json.Unmarshal(data, authorizeForPeerParam)
	if err != nil {
		log4.Error("json.Unmarshal failed ", err)
		return false
	}
	user, ok := common.GetAccountByPassword(ontSdk, authorizeForPeerParam.Path)
	if !ok {
		return false
	}
	ok = authorizeForPeer(ontSdk, user, authorizeForPeerParam.PeerPubkeyList, authorizeForPeerParam.PosList)
	if !ok {
		return false
	}
	common.WaitForBlock(ontSdk)
	return true
}

func UnAuthorizeForPeer(ontSdk *sdk.OntologySdk) bool {
	data, err := ioutil.ReadFile("./params/UnAuthorizeForPeer.json")
	if err != nil {
		log4.Error("ioutil.ReadFile failed ", err)
		return false
	}
	authorizeForPeerParam := new(AuthorizeForPeerParam)
	err = json.Unmarshal(data, authorizeForPeerParam)
	if err != nil {
		log4.Error("json.Unmarshal failed ", err)
		return false
	}
	user, ok := common.GetAccountByPassword(ontSdk, authorizeForPeerParam.Path)
	if !ok {
		return false
	}
	ok = unAuthorizeForPeer(ontSdk, user, authorizeForPeerParam.PeerPubkeyList, authorizeForPeerParam.PosList)
	if !ok {
		return false
	}
	common.WaitForBlock(ontSdk)
	return true
}

type WithdrawParam struct {
	Path           string
	PeerPubkeyList []string
	WithdrawList   []uint32
}

func Withdraw(ontSdk *sdk.OntologySdk) bool {
	data, err := ioutil.ReadFile("./params/Withdraw.json")
	if err != nil {
		log4.Error("ioutil.ReadFile failed ", err)
		return false
	}
	withdrawParam := new(WithdrawParam)
	err = json.Unmarshal(data, withdrawParam)
	if err != nil {
		log4.Error("json.Unmarshal failed ", err)
		return false
	}
	user, ok := common.GetAccountByPassword(ontSdk, withdrawParam.Path)
	if !ok {
		return false
	}
	ok = withdraw(ontSdk, user, withdrawParam.PeerPubkeyList, withdrawParam.WithdrawList)
	if !ok {
		return false
	}
	common.WaitForBlock(ontSdk)
	return true
}

type QuitNodeParam struct {
	Path       []string
	PeerPubkey []string
}

func QuitNode(ontSdk *sdk.OntologySdk) bool {
	data, err := ioutil.ReadFile("./params/QuitNode.json")
	if err != nil {
		log4.Error("ioutil.ReadFile failed ", err)
		return false
	}
	quitNodeParam := new(QuitNodeParam)
	err = json.Unmarshal(data, quitNodeParam)
	if err != nil {
		log4.Error("json.Unmarshal failed ", err)
		return false
	}
	time.Sleep(1 * time.Second)
	for i := 0; i < len(quitNodeParam.Path); i++ {
		user, ok := common.GetAccountByPassword(ontSdk, quitNodeParam.Path[i])
		if !ok {
			return false
		}
		ok = quitNode(ontSdk, user, quitNodeParam.PeerPubkey[i])
		if !ok {
			return false
		}
	}
	common.WaitForBlock(ontSdk)
	return true
}

type BlackNodeParam struct {
	Path           []string
	PeerPubkeyList []string
}

func BlackNode(ontSdk *sdk.OntologySdk) bool {
	data, err := ioutil.ReadFile("./params/BlackNode.json")
	if err != nil {
		log4.Error("ioutil.ReadFile failed ", err)
		return false
	}
	blackNodeParam := new(BlackNodeParam)
	err = json.Unmarshal(data, blackNodeParam)
	if err != nil {
		log4.Error("json.Unmarshal failed ", err)
		return false
	}
	var users []*sdk.Account
	var pubKeys []keypair.PublicKey
	time.Sleep(1 * time.Second)
	for _, path := range blackNodeParam.Path {
		user, ok := common.GetAccountByPassword(ontSdk, path)
		if !ok {
			return false
		}
		users = append(users, user)
		pubKeys = append(pubKeys, user.PublicKey)
	}
	ok := blackNodeMultiSign(ontSdk, pubKeys, users, blackNodeParam.PeerPubkeyList)
	if !ok {
		return false
	}
	common.WaitForBlock(ontSdk)
	return true
}

type WhiteNodeParam struct {
	Path       []string
	PeerPubkey string
}

func WhiteNode(ontSdk *sdk.OntologySdk) bool {
	data, err := ioutil.ReadFile("./params/WhiteNode.json")
	if err != nil {
		log4.Error("ioutil.ReadFile failed ", err)
		return false
	}
	whiteNodeParam := new(WhiteNodeParam)
	err = json.Unmarshal(data, whiteNodeParam)
	if err != nil {
		log4.Error("json.Unmarshal failed ", err)
		return false
	}
	var users []*sdk.Account
	var pubKeys []keypair.PublicKey
	time.Sleep(1 * time.Second)
	for _, path := range whiteNodeParam.Path {
		user, ok := common.GetAccountByPassword(ontSdk, path)
		if !ok {
			return false
		}
		users = append(users, user)
		pubKeys = append(pubKeys, user.PublicKey)
	}
	ok := whiteNodeMultiSign(ontSdk, pubKeys, users, whiteNodeParam.PeerPubkey)
	if !ok {
		return false
	}
	common.WaitForBlock(ontSdk)
	return true
}

type MultiAccount struct {
	Path []string
}

func CommitDpos(ontSdk *sdk.OntologySdk) bool {
	data, err := ioutil.ReadFile("./params/CommitDpos.json")
	if err != nil {
		log4.Error("ioutil.ReadFile failed ", err)
		return false
	}
	multiAccount := new(MultiAccount)
	err = json.Unmarshal(data, multiAccount)
	if err != nil {
		log4.Error("json.Unmarshal failed ", err)
		return false
	}
	var users []*sdk.Account
	var pubKeys []keypair.PublicKey
	time.Sleep(1 * time.Second)
	for _, path := range multiAccount.Path {
		user, ok := common.GetAccountByPassword(ontSdk, path)
		if !ok {
			return false
		}
		users = append(users, user)
		pubKeys = append(pubKeys, user.PublicKey)
	}
	ok := commitDposMultiSign(ontSdk, pubKeys, users)
	if !ok {
		return false
	}
	common.WaitForBlock(ontSdk)
	return true
}

type UpdateConfigParam struct {
	Path                 []string
	N                    uint32
	C                    uint32
	K                    uint32
	L                    uint32
	BlockMsgDelay        uint32
	HashMsgDelay         uint32
	PeerHandshakeTimeout uint32
	MaxBlockChangeView   uint32
}

func UpdateConfig(ontSdk *sdk.OntologySdk) bool {
	data, err := ioutil.ReadFile("./params/UpdateConfig.json")
	if err != nil {
		log4.Error("ioutil.ReadFile failed ", err)
		return false
	}
	updateConfigParam := new(UpdateConfigParam)
	err = json.Unmarshal(data, updateConfigParam)
	if err != nil {
		log4.Error("json.Unmarshal failed ", err)
		return false
	}
	var users []*sdk.Account
	var pubKeys []keypair.PublicKey
	time.Sleep(1 * time.Second)
	for _, path := range updateConfigParam.Path {
		user, ok := common.GetAccountByPassword(ontSdk, path)
		if !ok {
			return false
		}
		users = append(users, user)
		pubKeys = append(pubKeys, user.PublicKey)
	}
	config := &governance.Configuration{
		N:                    updateConfigParam.N,
		C:                    updateConfigParam.C,
		K:                    updateConfigParam.K,
		L:                    updateConfigParam.L,
		BlockMsgDelay:        updateConfigParam.BlockMsgDelay,
		HashMsgDelay:         updateConfigParam.HashMsgDelay,
		PeerHandshakeTimeout: updateConfigParam.PeerHandshakeTimeout,
		MaxBlockChangeView:   updateConfigParam.MaxBlockChangeView,
	}
	ok := updateConfigMultiSign(ontSdk, pubKeys, users, config)
	if !ok {
		return false
	}
	common.WaitForBlock(ontSdk)
	return true
}

type UpdateGlobalParamParam struct {
	Path         []string
	CandidateFee uint64
	MinInitStake uint32
	CandidateNum uint32
	PosLimit     uint32
	A            uint32
	B            uint32
	Yita         uint32
	Penalty      uint32
}

func UpdateGlobalParam(ontSdk *sdk.OntologySdk) bool {
	data, err := ioutil.ReadFile("./params/UpdateGlobalParam.json")
	if err != nil {
		log4.Error("ioutil.ReadFile failed ", err)
		return false
	}
	updateGlobalParamParam := new(UpdateGlobalParamParam)
	err = json.Unmarshal(data, updateGlobalParamParam)
	if err != nil {
		log4.Error("json.Unmarshal failed ", err)
		return false
	}
	var users []*sdk.Account
	var pubKeys []keypair.PublicKey
	time.Sleep(1 * time.Second)
	for _, path := range updateGlobalParamParam.Path {
		user, ok := common.GetAccountByPassword(ontSdk, path)
		if !ok {
			return false
		}
		users = append(users, user)
		pubKeys = append(pubKeys, user.PublicKey)
	}
	globalParam := &governance.GlobalParam{
		CandidateFee: updateGlobalParamParam.CandidateFee,
		MinInitStake: updateGlobalParamParam.MinInitStake,
		CandidateNum: updateGlobalParamParam.CandidateNum,
		PosLimit:     updateGlobalParamParam.PosLimit,
		A:            updateGlobalParamParam.A,
		B:            updateGlobalParamParam.B,
		Yita:         updateGlobalParamParam.Yita,
		Penalty:      updateGlobalParamParam.Penalty,
	}
	ok := updateGlobalParamMultiSign(ontSdk, pubKeys, users, globalParam)
	if !ok {
		return false
	}
	common.WaitForBlock(ontSdk)
	return true
}

type UpdateGlobalParamParam2 struct {
	Path                 []string
	MinAuthorizePos      uint32
	CandidateFeeSplitNum uint32
}

func UpdateGlobalParam2(ontSdk *sdk.OntologySdk) bool {
	data, err := ioutil.ReadFile("./params/UpdateGlobalParam2.json")
	if err != nil {
		log4.Error("ioutil.ReadFile failed ", err)
		return false
	}
	updateGlobalParamParam2 := new(UpdateGlobalParamParam2)
	err = json.Unmarshal(data, updateGlobalParamParam2)
	if err != nil {
		log4.Error("json.Unmarshal failed ", err)
		return false
	}
	var users []*sdk.Account
	var pubKeys []keypair.PublicKey
	time.Sleep(1 * time.Second)
	for _, path := range updateGlobalParamParam2.Path {
		user, ok := common.GetAccountByPassword(ontSdk, path)
		if !ok {
			return false
		}
		users = append(users, user)
		pubKeys = append(pubKeys, user.PublicKey)
	}
	globalParam2 := &governance.GlobalParam2{
		MinAuthorizePos:      updateGlobalParamParam2.MinAuthorizePos,
		CandidateFeeSplitNum: updateGlobalParamParam2.CandidateFeeSplitNum,
	}
	ok := updateGlobalParam2MultiSign(ontSdk, pubKeys, users, globalParam2)
	if !ok {
		return false
	}
	common.WaitForBlock(ontSdk)
	return true
}

type UpdateSplitCurveParam struct {
	Path []string
	Yi   []uint32
}

func UpdateSplitCurve(ontSdk *sdk.OntologySdk) bool {
	data, err := ioutil.ReadFile("./params/UpdateSplitCurve.json")
	if err != nil {
		log4.Error("ioutil.ReadFile failed ", err)
		return false
	}
	updateSplitCurveParam := new(UpdateSplitCurveParam)
	err = json.Unmarshal(data, updateSplitCurveParam)
	if err != nil {
		log4.Error("json.Unmarshal failed ", err)
		return false
	}
	var users []*sdk.Account
	var pubKeys []keypair.PublicKey
	time.Sleep(1 * time.Second)
	for _, path := range updateSplitCurveParam.Path {
		user, ok := common.GetAccountByPassword(ontSdk, path)
		if !ok {
			return false
		}
		users = append(users, user)
		pubKeys = append(pubKeys, user.PublicKey)
	}
	splitCurve := &governance.SplitCurve{
		Yi: updateSplitCurveParam.Yi,
	}
	ok := updateSplitCurveMultiSign(ontSdk, pubKeys, users, splitCurve)
	if !ok {
		return false
	}
	common.WaitForBlock(ontSdk)
	return true
}

type SetPromisePosParam struct {
	Path       []string
	PeerPubkey []string
	PromisePos []uint64
}

func SetPromisePos(ontSdk *sdk.OntologySdk) bool {
	data, err := ioutil.ReadFile("./params/SetPromisePos.json")
	if err != nil {
		log4.Error("ioutil.ReadFile failed ", err)
		return false
	}
	setPromisePosParam := new(SetPromisePosParam)
	err = json.Unmarshal(data, setPromisePosParam)
	if err != nil {
		log4.Error("json.Unmarshal failed ", err)
		return false
	}
	var users []*sdk.Account
	var pubKeys []keypair.PublicKey
	time.Sleep(1 * time.Second)
	for _, path := range setPromisePosParam.Path {
		user, ok := common.GetAccountByPassword(ontSdk, path)
		if !ok {
			return false
		}
		users = append(users, user)
		pubKeys = append(pubKeys, user.PublicKey)
	}
	for index, peerPubkey := range setPromisePosParam.PeerPubkey {
		promisePos := &governance.PromisePos{
			PeerPubkey: peerPubkey,
			PromisePos: setPromisePosParam.PromisePos[index],
		}
		ok := setPromisePosMultiSign(ontSdk, pubKeys, users, promisePos)
		if !ok {
			return false
		}
	}
	common.WaitForBlock(ontSdk)
	return true
}

type TransferPenaltyParam struct {
	Path1      []string
	PeerPubkey string
	Path2      string
}

func TransferPenalty(ontSdk *sdk.OntologySdk) bool {
	data, err := ioutil.ReadFile("./params/TransferPenalty.json")
	if err != nil {
		log4.Error("ioutil.ReadFile failed ", err)
		return false
	}
	transferPenaltyParam := new(TransferPenaltyParam)
	err = json.Unmarshal(data, transferPenaltyParam)
	if err != nil {
		log4.Error("json.Unmarshal failed ", err)
		return false
	}
	var users []*sdk.Account
	var pubKeys []keypair.PublicKey
	time.Sleep(1 * time.Second)
	for _, path := range transferPenaltyParam.Path1 {
		user, ok := common.GetAccountByPassword(ontSdk, path)
		if !ok {
			return false
		}
		users = append(users, user)
		pubKeys = append(pubKeys, user.PublicKey)
	}
	user1, ok := common.GetAccountByPassword(ontSdk, transferPenaltyParam.Path2)
	if !ok {
		return false
	}
	ok = transferPenaltyMultiSign(ontSdk, pubKeys, users, transferPenaltyParam.PeerPubkey, user1.Address)
	if !ok {
		return false
	}
	common.WaitForBlock(ontSdk)
	return true
}

func GetVbftConfig(ontSdk *sdk.OntologySdk) bool {
	config, err := getVbftConfig(ontSdk)
	if err != nil {
		log4.Error("getVbftConfig failed ", err)
		return false
	}
	fmt.Println("config.N is:", config.N)
	fmt.Println("config.C is:", config.C)
	fmt.Println("config.K is:", config.K)
	fmt.Println("config.L is:", config.L)
	fmt.Println("config.BlockMsgDelay is:", config.BlockMsgDelay)
	fmt.Println("config.HashMsgDelay is:", config.HashMsgDelay)
	fmt.Println("config.PeerHandshakeTimeout is:", config.PeerHandshakeTimeout)
	fmt.Println("config.MaxBlockChangeView is:", config.MaxBlockChangeView)
	return true
}

func GetGlobalParam(ontSdk *sdk.OntologySdk) bool {
	globalParam, err := getGlobalParam(ontSdk)
	if err != nil {
		log4.Error("getGlobalParam failed ", err)
		return false
	}
	fmt.Println("globalParam.CandidateFee is:", globalParam.CandidateFee)
	fmt.Println("globalParam.MinInitStake is:", globalParam.MinInitStake)
	fmt.Println("globalParam.CandidateNum is:", globalParam.CandidateNum)
	fmt.Println("globalParam.PosLimit is:", globalParam.PosLimit)
	fmt.Println("globalParam.A is:", globalParam.A)
	fmt.Println("globalParam.B is:", globalParam.B)
	fmt.Println("globalParam.Yita is:", globalParam.Yita)
	fmt.Println("globalParam.Penalty is:", globalParam.Penalty)
	return true
}

func GetGlobalParam2(ontSdk *sdk.OntologySdk) bool {
	globalParam2, err := getGlobalParam2(ontSdk)
	if err != nil {
		log4.Error("getGlobalParam failed ", err)
		return false
	}
	fmt.Println("globalParam2.MinAuthorizePos is:", globalParam2.MinAuthorizePos)
	fmt.Println("globalParam2.CandidateFeeSplitNum is:", globalParam2.CandidateFeeSplitNum)
	return true
}

func GetSplitCurve(ontSdk *sdk.OntologySdk) bool {
	splitCurve, err := getSplitCurve(ontSdk)
	if err != nil {
		log4.Error("getSplitCurve failed ", err)
		return false
	}
	fmt.Println("splitCurve.Yi is", splitCurve.Yi)
	return true
}

func GetGovernanceView(ontSdk *sdk.OntologySdk) bool {
	governanceView, err := getGovernanceView(ontSdk)
	if err != nil {
		log4.Error("getGovernanceView failed ", err)
		return false
	}
	fmt.Println("governanceView.View is:", governanceView.View)
	fmt.Println("governanceView.TxHash is:", governanceView.TxHash)
	fmt.Println("governanceView.Height is:", governanceView.Height)
	return true
}

type GetPeerPoolItemParam struct {
	PeerPubkey string
}

func GetPeerPoolItem(ontSdk *sdk.OntologySdk) bool {
	data, err := ioutil.ReadFile("./params/GetPeerPoolItem.json")
	if err != nil {
		log4.Error("ioutil.ReadFile failed ", err)
		return false
	}
	getPeerPoolItemParam := new(GetPeerPoolItemParam)
	err = json.Unmarshal(data, getPeerPoolItemParam)
	if err != nil {
		log4.Error("json.Unmarshal failed ", err)
		return false
	}

	peerPoolMap, err := getPeerPoolMap(ontSdk)
	if err != nil {
		log4.Error("getPeerPoolMap failed ", err)
		return false
	}

	peerPoolItem, ok := peerPoolMap.PeerPoolMap[getPeerPoolItemParam.PeerPubkey]
	if !ok {
		fmt.Println("Can't find peerPubkey in peerPoolMap")
	}
	fmt.Println("peerPoolItem.Index is:", peerPoolItem.Index)
	fmt.Println("peerPoolItem.PeerPubkey is:", peerPoolItem.PeerPubkey)
	fmt.Println("peerPoolItem.Address is:", peerPoolItem.Address.ToBase58())
	fmt.Println("peerPoolItem.Status is:", peerPoolItem.Status)
	fmt.Println("peerPoolItem.InitPos is:", peerPoolItem.InitPos)
	fmt.Println("peerPoolItem.TotalPos is:", peerPoolItem.TotalPos)
	return true
}

func GetPeerPoolMap(ontSdk *sdk.OntologySdk) bool {
	peerPoolMap, err := getPeerPoolMap(ontSdk)
	if err != nil {
		log4.Error("getPeerPoolMap failed ", err)
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

type GetAuthorizeInfoParam struct {
	Address    string
	PeerPubkey string
}

func GetAuthorizeInfo(ontSdk *sdk.OntologySdk) bool {
	data, err := ioutil.ReadFile("./params/GetAuthorizeInfo.json")
	if err != nil {
		log4.Error("ioutil.ReadFile failed ", err)
		return false
	}
	getAuthorizeInfoParam := new(GetAuthorizeInfoParam)
	err = json.Unmarshal(data, getAuthorizeInfoParam)
	if err != nil {
		log4.Error("json.Unmarshal failed ", err)
		return false
	}

	address, err := ocommon.AddressFromBase58(getAuthorizeInfoParam.Address)
	if err != nil {
		log4.Error("common.AddressFromBase58 failed ", err)
		return false
	}
	authorizeInfo, err := getAuthorizeInfo(ontSdk, getAuthorizeInfoParam.PeerPubkey, address)
	if err != nil {
		log4.Error("getAuthorizeInfo failed ", err)
		return false
	}

	fmt.Println("authorizeInfo.PeerPubkey is:", authorizeInfo.PeerPubkey)
	fmt.Println("authorizeInfo.Address is:", authorizeInfo.Address.ToBase58())
	fmt.Println("authorizeInfo.ConsensusPos is:", authorizeInfo.ConsensusPos)
	fmt.Println("authorizeInfo.CandidatePos is:", authorizeInfo.CandidatePos)
	fmt.Println("authorizeInfo.NewPos is:", authorizeInfo.NewPos)
	fmt.Println("authorizeInfo.WithdrawConsensusPos is:", authorizeInfo.WithdrawConsensusPos)
	fmt.Println("authorizeInfo.WithdrawCandidatePos is:", authorizeInfo.WithdrawCandidatePos)
	fmt.Println("authorizeInfo.WithdrawUnfreezePos is:", authorizeInfo.WithdrawUnfreezePos)
	return true
}

type GetTotalStakeParam struct {
	Address string
}

func GetTotalStake(ontSdk *sdk.OntologySdk) bool {
	data, err := ioutil.ReadFile("./params/GetTotalStake.json")
	if err != nil {
		log4.Error("ioutil.ReadFile failed ", err)
		return false
	}
	getTotalStakeParam := new(GetTotalStakeParam)
	err = json.Unmarshal(data, getTotalStakeParam)
	if err != nil {
		log4.Error("json.Unmarshal failed ", err)
		return false
	}
	address, err := ocommon.AddressFromBase58(getTotalStakeParam.Address)
	if err != nil {
		log4.Error("common.AddressFromBase58 failed ", err)
		return false
	}

	totalStake, err := getTotalStake(ontSdk, address)
	if err != nil {
		log4.Error("getTotalStake failed ", err)
		return false
	}

	fmt.Println("totalStake.Address is:", totalStake.Address.ToBase58())
	fmt.Println("totalStake.Stake is:", totalStake.Stake)
	fmt.Println("totalStake.TimeOffset is:", totalStake.TimeOffset)
	return true
}

type GetPenaltyStakeParam struct {
	PeerPubkey string
}

func GetPenaltyStake(ontSdk *sdk.OntologySdk) bool {
	data, err := ioutil.ReadFile("./params/GetPenaltyStake.json")
	if err != nil {
		log4.Error("ioutil.ReadFile failed ", err)
		return false
	}
	getPenaltyStakeParam := new(GetPenaltyStakeParam)
	err = json.Unmarshal(data, getPenaltyStakeParam)
	if err != nil {
		log4.Error("json.Unmarshal failed ", err)
		return false
	}

	penaltyStake, err := getPenaltyStake(ontSdk, getPenaltyStakeParam.PeerPubkey)
	if err != nil {
		log4.Error("getPenaltyStake failed ", err)
		return false
	}

	fmt.Println("penaltyStake.PeerPubkey is:", penaltyStake.PeerPubkey)
	fmt.Println("penaltyStake.InitPos is:", penaltyStake.InitPos)
	fmt.Println("penaltyStake.AuthorizePos is:", penaltyStake.AuthorizePos)
	fmt.Println("penaltyStake.TimeOffset is:", penaltyStake.TimeOffset)
	fmt.Println("penaltyStake.Amount is:", penaltyStake.Amount)
	return true
}

type InBlackListParam struct {
	PeerPubkey string
}

func InBlackList(ontSdk *sdk.OntologySdk) bool {
	data, err := ioutil.ReadFile("./params/InBlackList.json")
	if err != nil {
		log4.Error("ioutil.ReadFile failed ", err)
		return false
	}
	inBlackListParam := new(InBlackListParam)
	err = json.Unmarshal(data, inBlackListParam)
	if err != nil {
		log4.Error("json.Unmarshal failed ", err)
		return false
	}

	inBlackList, err := inBlackList(ontSdk, inBlackListParam.PeerPubkey)
	if err != nil {
		log4.Error("getPenaltyStake failed ", err)
		return false
	}

	fmt.Println("result is:", inBlackList)
	return true
}

type WithdrawOngParam struct {
	Path       string
	PeerPubkey string
}

func WithdrawOng(ontSdk *sdk.OntologySdk) bool {
	data, err := ioutil.ReadFile("./params/WithdrawOng.json")
	if err != nil {
		log4.Error("ioutil.ReadFile failed ", err)
		return false
	}
	withdrawOngParam := new(WithdrawOngParam)
	err = json.Unmarshal(data, withdrawOngParam)
	if err != nil {
		log4.Error("json.Unmarshal failed ", err)
		return false
	}
	user, ok := common.GetAccountByPassword(ontSdk, withdrawOngParam.Path)
	if !ok {
		return false
	}
	ok = withdrawOng(ontSdk, user)
	if !ok {
		return false
	}
	common.WaitForBlock(ontSdk)
	return true
}

type VrfParam struct {
	Path string
}

type vrfData struct {
	BlockNum uint32 `json:"block_num"`
	PrevVrf  []byte `json:"prev_vrf"`
}

func Vrf(ontSdk *sdk.OntologySdk) bool {
	data, err := ioutil.ReadFile("./params/Vrf.json")
	if err != nil {
		log4.Error("ioutil.ReadFile failed ", err)
		return false
	}
	vrfParam := new(VrfParam)
	err = json.Unmarshal(data, vrfParam)
	if err != nil {
		log4.Error("json.Unmarshal failed ", err)
		return false
	}
	user, ok := common.GetAccountByPassword(ontSdk, vrfParam.Path)
	if !ok {
		return false
	}

	data, err = json.Marshal(&vrfData{
		BlockNum: 0,
		PrevVrf:  keypair.SerializePublicKey(user.PublicKey),
	})
	if err != nil {
		log4.Error("json.Unmarshal vrf payload failed ", err)
		return false
	}

	value, proof, err := vrf.Vrf(user.PrivateKey, data)
	if err != nil {
		log4.Error("vrf computation failed ", err)
		return false
	}

	if ok, err := vrf.Verify(user.PublicKey, data, value, proof); err != nil || !ok {
		log4.Error("vrf verify failed: ", err)
		return false
	}

	log4.Info("vrf value: %s", hex.EncodeToString(value))
	log4.Info("vrf proof: %s", hex.EncodeToString(proof))

	return true
}

type TransferMultiSignParam struct {
	Path1  []string
	Path2  []string
	Amount []uint64
}

func TransferOntMultiSign(ontSdk *sdk.OntologySdk) bool {
	data, err := ioutil.ReadFile("./params/TransferOntMultiSign.json")
	if err != nil {
		log4.Error("ioutil.ReadFile failed ", err)
		return false
	}
	transferMultiSignParam := new(TransferMultiSignParam)
	err = json.Unmarshal(data, transferMultiSignParam)
	if err != nil {
		log4.Error("json.Unmarshal failed ", err)
		return false
	}
	var users []*sdk.Account
	var pubKeys []keypair.PublicKey
	time.Sleep(1 * time.Second)
	for _, path := range transferMultiSignParam.Path1 {
		user, ok := common.GetAccountByPassword(ontSdk, path)
		if !ok {
			return false
		}
		users = append(users, user)
		pubKeys = append(pubKeys, user.PublicKey)
	}
	time.Sleep(1 * time.Second)
	for index, path2 := range transferMultiSignParam.Path2 {
		user2, ok := common.GetAccountByPassword(ontSdk, path2)
		if !ok {
			return false
		}
		ok = transferOntMultiSign(ontSdk, pubKeys, users, user2.Address, transferMultiSignParam.Amount[index])
		if !ok {
			return false
		}
	}
	common.WaitForBlock(ontSdk)
	return true
}

func TransferOngMultiSign(ontSdk *sdk.OntologySdk) bool {
	data, err := ioutil.ReadFile("./params/TransferOngMultiSign.json")
	if err != nil {
		log4.Error("ioutil.ReadFile failed ", err)
		return false
	}
	transferMultiSignParam := new(TransferMultiSignParam)
	err = json.Unmarshal(data, transferMultiSignParam)
	if err != nil {
		log4.Error("json.Unmarshal failed ", err)
		return false
	}
	var users []*sdk.Account
	var pubKeys []keypair.PublicKey
	time.Sleep(1 * time.Second)
	for _, path := range transferMultiSignParam.Path1 {
		user, ok := common.GetAccountByPassword(ontSdk, path)
		if !ok {
			return false
		}
		users = append(users, user)
		pubKeys = append(pubKeys, user.PublicKey)
	}
	time.Sleep(1 * time.Second)
	for index, path2 := range transferMultiSignParam.Path2 {
		user2, ok := common.GetAccountByPassword(ontSdk, path2)
		if !ok {
			return false
		}
		ok = transferOngMultiSign(ontSdk, pubKeys, users, user2.Address, transferMultiSignParam.Amount[index])
		if !ok {
			return false
		}
	}
	common.WaitForBlock(ontSdk)
	return true
}

type TransferFromMultiSignParam struct {
	Path1  []string
	Path2  []string
	Amount []uint64
}

func TransferFromOngMultiSign(ontSdk *sdk.OntologySdk) bool {
	data, err := ioutil.ReadFile("./params/TransferFromOngMultiSign.json")
	if err != nil {
		log4.Error("ioutil.ReadFile failed ", err)
		return false
	}
	transferFromMultiSignParam := new(TransferFromMultiSignParam)
	err = json.Unmarshal(data, transferFromMultiSignParam)
	if err != nil {
		log4.Error("json.Unmarshal failed ", err)
		return false
	}
	var users []*sdk.Account
	var pubKeys []keypair.PublicKey
	time.Sleep(1 * time.Second)
	for _, path := range transferFromMultiSignParam.Path1 {
		user, ok := common.GetAccountByPassword(ontSdk, path)
		if !ok {
			return false
		}
		users = append(users, user)
		pubKeys = append(pubKeys, user.PublicKey)
	}
	time.Sleep(1 * time.Second)
	for index, path2 := range transferFromMultiSignParam.Path2 {
		user2, ok := common.GetAccountByPassword(ontSdk, path2)
		if !ok {
			return false
		}
		ok = transferFromOngMultiSign(ontSdk, pubKeys, users, user2.Address, transferFromMultiSignParam.Amount[index])
		if !ok {
			return false
		}
	}
	common.WaitForBlock(ontSdk)
	return true
}

type GetAddressMultiSignParam struct {
	PubKeys []string
}

func GetAddressMultiSign(ontSdk *sdk.OntologySdk) bool {
	data, err := ioutil.ReadFile("./params/GetAddressMultiSign.json")
	if err != nil {
		log4.Error("ioutil.ReadFile failed ", err)
		return false
	}
	getAddressMultiSignParam := new(GetAddressMultiSignParam)
	err = json.Unmarshal(data, getAddressMultiSignParam)
	if err != nil {
		log4.Error("json.Unmarshal failed ", err)
		return false
	}
	var pubKeys []keypair.PublicKey
	time.Sleep(1 * time.Second)
	for _, v := range getAddressMultiSignParam.PubKeys {
		vByte, err := hex.DecodeString(v)
		if err != nil {
			log4.Error("hex.DecodeString failed ", err)
		}
		k, err := keypair.DeserializePublicKey(vByte)
		if err != nil {
			log4.Error("keypair.DeserializePublicKey failed ", err)
		}
		pubKeys = append(pubKeys, k)
	}
	from, err := types.AddressFromMultiPubKeys(pubKeys, int((5*len(pubKeys)+6)/7))
	if err != nil {
		log4.Error("types.AddressFromMultiPubKeys error", err)
	}
	fmt.Println("address is:", from.ToBase58())
	return true
}

type TransferMultiSignToMultiSignParam struct {
	Path1   []string
	PubKeys []string
	Amount  uint64
}

func TransferOntMultiSignToMultiSign(ontSdk *sdk.OntologySdk) bool {
	data, err := ioutil.ReadFile("./params/TransferOntMultiSignToMultiSign.json")
	if err != nil {
		log4.Error("ioutil.ReadFile failed ", err)
		return false
	}
	transferMultiSignToMultiSignParam := new(TransferMultiSignToMultiSignParam)
	err = json.Unmarshal(data, transferMultiSignToMultiSignParam)
	if err != nil {
		log4.Error("json.Unmarshal failed ", err)
		return false
	}
	var users []*sdk.Account
	var pubKeys []keypair.PublicKey
	var pubKeysTo []keypair.PublicKey
	time.Sleep(1 * time.Second)
	for _, path := range transferMultiSignToMultiSignParam.Path1 {
		user, ok := common.GetAccountByPassword(ontSdk, path)
		if !ok {
			return false
		}
		users = append(users, user)
		pubKeys = append(pubKeys, user.PublicKey)
	}
	for _, v := range transferMultiSignToMultiSignParam.PubKeys {
		vByte, err := hex.DecodeString(v)
		if err != nil {
			log4.Error("hex.DecodeString failed ", err)
		}
		k, err := keypair.DeserializePublicKey(vByte)
		if err != nil {
			log4.Error("keypair.DeserializePublicKey failed ", err)
		}
		pubKeysTo = append(pubKeysTo, k)
	}
	to, err := types.AddressFromMultiPubKeys(pubKeysTo, int((5*len(pubKeysTo)+6)/7))
	if err != nil {
		log4.Error("types.AddressFromMultiPubKeys error", err)
	}
	ok := transferOntMultiSignToMultiSign(ontSdk, pubKeys, users, to, transferMultiSignToMultiSignParam.Amount)
	if !ok {
		return false
	}
	common.WaitForBlock(ontSdk)
	return true
}

func TransferOngMultiSignToMultiSign(ontSdk *sdk.OntologySdk) bool {
	data, err := ioutil.ReadFile("./params/TransferOngMultiSignToMultiSign.json")
	if err != nil {
		log4.Error("ioutil.ReadFile failed ", err)
		return false
	}
	transferMultiSignToMultiSignParam := new(TransferMultiSignToMultiSignParam)
	err = json.Unmarshal(data, transferMultiSignToMultiSignParam)
	if err != nil {
		log4.Error("json.Unmarshal failed ", err)
		return false
	}
	var users []*sdk.Account
	var pubKeys []keypair.PublicKey
	var pubKeysTo []keypair.PublicKey
	time.Sleep(1 * time.Second)
	for _, path := range transferMultiSignToMultiSignParam.Path1 {
		user, ok := common.GetAccountByPassword(ontSdk, path)
		if !ok {
			return false
		}
		users = append(users, user)
		pubKeys = append(pubKeys, user.PublicKey)
	}
	for _, v := range transferMultiSignToMultiSignParam.PubKeys {
		vByte, err := hex.DecodeString(v)
		if err != nil {
			log4.Error("hex.DecodeString failed ", err)
		}
		k, err := keypair.DeserializePublicKey(vByte)
		if err != nil {
			log4.Error("keypair.DeserializePublicKey failed ", err)
		}
		pubKeysTo = append(pubKeysTo, k)
	}
	to, err := types.AddressFromMultiPubKeys(pubKeysTo, int((5*len(pubKeysTo)+6)/7))
	if err != nil {
		log4.Error("types.AddressFromMultiPubKeys error", err)
	}
	ok := transferOngMultiSignToMultiSign(ontSdk, pubKeys, users, to, transferMultiSignToMultiSignParam.Amount)
	if !ok {
		return false
	}
	common.WaitForBlock(ontSdk)
	return true
}

type TransferFromMultiSignToMultiSignParam struct {
	Path1   []string
	PubKeys []string
	Amount  uint64
}

func TransferFromOngMultiSignToMultiSign(ontSdk *sdk.OntologySdk) bool {
	data, err := ioutil.ReadFile("./params/TransferFromOngMultiSignToMultiSign.json")
	if err != nil {
		log4.Error("ioutil.ReadFile failed ", err)
		return false
	}
	transferFromMultiSignToMultiSignParam := new(TransferFromMultiSignToMultiSignParam)
	err = json.Unmarshal(data, transferFromMultiSignToMultiSignParam)
	if err != nil {
		log4.Error("json.Unmarshal failed ", err)
		return false
	}
	var users []*sdk.Account
	var pubKeys []keypair.PublicKey
	var pubKeysTo []keypair.PublicKey
	time.Sleep(1 * time.Second)
	for _, path := range transferFromMultiSignToMultiSignParam.Path1 {
		user, ok := common.GetAccountByPassword(ontSdk, path)
		if !ok {
			return false
		}
		users = append(users, user)
		pubKeys = append(pubKeys, user.PublicKey)
	}
	for _, v := range transferFromMultiSignToMultiSignParam.PubKeys {
		vByte, err := hex.DecodeString(v)
		if err != nil {
			log4.Error("hex.DecodeString failed ", err)
		}
		k, err := keypair.DeserializePublicKey(vByte)
		if err != nil {
			log4.Error("keypair.DeserializePublicKey failed ", err)
		}
		pubKeysTo = append(pubKeysTo, k)
	}
	to, err := types.AddressFromMultiPubKeys(pubKeysTo, int((5*len(pubKeysTo)+6)/7))
	if err != nil {
		log4.Error("types.AddressFromMultiPubKeys error", err)
	}
	ok := transferFromOngMultiSignToMultiSign(ontSdk, pubKeys, users, to, transferFromMultiSignToMultiSignParam.Amount)
	if !ok {
		return false
	}
	common.WaitForBlock(ontSdk)
	return true
}

type TransferMultiSignAddressParam struct {
	Path1   []string
	PubKeys []string
	Address []string
	Amount  []uint64
}

func TransferOntMultiSignAddress(ontSdk *sdk.OntologySdk) bool {
	data, err := ioutil.ReadFile("./params/TransferOntMultiSignAddress.json")
	if err != nil {
		log4.Error("ioutil.ReadFile failed ", err)
		return false
	}
	transferMultiSignAddressParam := new(TransferMultiSignAddressParam)
	err = json.Unmarshal(data, transferMultiSignAddressParam)
	if err != nil {
		log4.Error("json.Unmarshal failed ", err)
		return false
	}
	var users []*sdk.Account
	var pubKeys []keypair.PublicKey
	time.Sleep(1 * time.Second)
	for _, path := range transferMultiSignAddressParam.Path1 {
		user, ok := common.GetAccountByPassword(ontSdk, path)
		if !ok {
			return false
		}
		users = append(users, user)
	}
	for _, v := range transferMultiSignAddressParam.PubKeys {
		vByte, err := hex.DecodeString(v)
		if err != nil {
			log4.Error("hex.DecodeString failed ", err)
		}
		k, err := keypair.DeserializePublicKey(vByte)
		if err != nil {
			log4.Error("keypair.DeserializePublicKey failed ", err)
		}
		pubKeys = append(pubKeys, k)
	}
	for index, address := range transferMultiSignAddressParam.Address {
		addr, err := ocommon.AddressFromBase58(address)
		if err != nil {
			log4.Error("common.AddressFromBase58 failed ", err)
			return false
		}
		ok := transferOntMultiSign(ontSdk, pubKeys, users, addr, transferMultiSignAddressParam.Amount[index])
		if !ok {
			return false
		}
	}
	common.WaitForBlock(ontSdk)
	return true
}

func TransferOngMultiSignAddress(ontSdk *sdk.OntologySdk) bool {
	data, err := ioutil.ReadFile("./params/TransferOngMultiSignAddress.json")
	if err != nil {
		log4.Error("ioutil.ReadFile failed ", err)
		return false
	}
	transferMultiSignAddressParam := new(TransferMultiSignAddressParam)
	err = json.Unmarshal(data, transferMultiSignAddressParam)
	if err != nil {
		log4.Error("json.Unmarshal failed ", err)
		return false
	}
	var users []*sdk.Account
	var pubKeys []keypair.PublicKey
	time.Sleep(1 * time.Second)
	for _, path := range transferMultiSignAddressParam.Path1 {
		user, ok := common.GetAccountByPassword(ontSdk, path)
		if !ok {
			return false
		}
		users = append(users, user)
	}
	for _, v := range transferMultiSignAddressParam.PubKeys {
		vByte, err := hex.DecodeString(v)
		if err != nil {
			log4.Error("hex.DecodeString failed ", err)
		}
		k, err := keypair.DeserializePublicKey(vByte)
		if err != nil {
			log4.Error("keypair.DeserializePublicKey failed ", err)
		}
		pubKeys = append(pubKeys, k)
	}
	for index, address := range transferMultiSignAddressParam.Address {
		addr, err := ocommon.AddressFromBase58(address)
		if err != nil {
			log4.Error("common.AddressFromBase58 failed ", err)
			return false
		}
		ok := transferOngMultiSign(ontSdk, pubKeys, users, addr, transferMultiSignAddressParam.Amount[index])
		if !ok {
			return false
		}
	}
	common.WaitForBlock(ontSdk)
	return true
}

type TransferFromMultiSignAddressParam struct {
	Path1   []string
	Address []string
	Amount  []uint64
}

func TransferFromOngMultiSignAddress(ontSdk *sdk.OntologySdk) bool {
	data, err := ioutil.ReadFile("./params/TransferFromOngMultiSignAddress.json")
	if err != nil {
		log4.Error("ioutil.ReadFile failed ", err)
		return false
	}
	transferFromMultiSignAddressParam := new(TransferFromMultiSignAddressParam)
	err = json.Unmarshal(data, transferFromMultiSignAddressParam)
	if err != nil {
		log4.Error("json.Unmarshal failed ", err)
		return false
	}
	var users []*sdk.Account
	var pubKeys []keypair.PublicKey
	time.Sleep(1 * time.Second)
	for _, path := range transferFromMultiSignAddressParam.Path1 {
		user, ok := common.GetAccountByPassword(ontSdk, path)
		if !ok {
			return false
		}
		users = append(users, user)
		pubKeys = append(pubKeys, user.PublicKey)
	}
	time.Sleep(1 * time.Second)
	for index, address := range transferFromMultiSignAddressParam.Address {
		addr, err := ocommon.AddressFromBase58(address)
		if err != nil {
			log4.Error("common.AddressFromBase58 failed ", err)
			return false
		}
		ok := transferFromOngMultiSign(ontSdk, pubKeys, users, addr, transferFromMultiSignAddressParam.Amount[index])
		if !ok {
			return false
		}
	}
	common.WaitForBlock(ontSdk)
	return true
}
func GetVbftInfo(ontSdk *sdk.OntologySdk) bool {
	blkNum, err := ontSdk.GetCurrentBlockHeight()
	if err != nil {
		log4.Error("TestGetVbftInfo GetBlockCount error:%s", err)
		return false
	}
	blk, err := ontSdk.GetBlockByHeight(blkNum - 1)
	if err != nil {
		log4.Error("TestGetVbftInfo GetBlockByHeight error:%s", err)
		return false
	}
	block, err := common.InitVbftBlock(blk)
	if err != nil {
		log4.Error("TestGetVbftInfo initVbftBlock error:%s", err)
		return false
	}

	var cfg vconfig.ChainConfig
	if block.Info.NewChainConfig != nil {
		cfg = *block.Info.NewChainConfig
	} else {
		var cfgBlock *types.Block
		if block.Info.LastConfigBlockNum != math.MaxUint32 {
			cfgBlock, err = ontSdk.GetBlockByHeight(block.Info.LastConfigBlockNum)
			if err != nil {
				log4.Error("TestGetVbftInfo chainconfig GetBlockByHeight error:%s", err)
				return false
			}
		}
		blk, err := common.InitVbftBlock(cfgBlock)
		if err != nil {
			log4.Error("TestGetVbftInfo initVbftBlock error:%s", err)
			return false
		}
		if blk.Info.NewChainConfig == nil {
			log4.Error("TestGetVbftInfo newchainconfig error:%s", err)
			return false
		}
		cfg = *blk.Info.NewChainConfig
	}
	fmt.Printf("block vbft chainConfig, View:%d, N:%d, C:%d, BlockMsgDelay:, HashMsgDelay:, PeerHandshakeTimeout:, MaxBlockChangeView:%d, PosTable:\n",
		cfg.View, cfg.N, cfg.C, cfg.BlockMsgDelay, cfg.HashMsgDelay, cfg.PeerHandshakeTimeout, cfg.MaxBlockChangeView, cfg.PosTable)
	for _, p := range cfg.Peers {
		fmt.Printf("peerInfo Index: %d, ID:\n", p.Index, p.ID)
	}
	return true
}

type MultiTransferParam struct {
	FromPath  []string
	ToAddress []string
	Amount    []uint64
}

func MultiTransferOnt(ontSdk *sdk.OntologySdk) bool {
	data, err := ioutil.ReadFile("./params/MultiTransferOnt.json")
	if err != nil {
		log4.Error("ioutil.ReadFile failed ", err)
		return false
	}
	multiTransferParam := new(MultiTransferParam)
	err = json.Unmarshal(data, multiTransferParam)
	if err != nil {
		log4.Error("json.Unmarshal failed ", err)
		return false
	}
	var users []*sdk.Account
	time.Sleep(1 * time.Second)
	for _, path := range multiTransferParam.FromPath {
		user, ok := common.GetAccountByPassword(ontSdk, path)
		if !ok {
			return false
		}
		users = append(users, user)
	}
	ok := multiTransfer(ontSdk, utils.OntContractAddress, users, multiTransferParam.ToAddress, multiTransferParam.Amount)
	if !ok {
		return false
	}
	common.WaitForBlock(ontSdk)
	return true
}

func MultiTransferOng(ontSdk *sdk.OntologySdk) bool {
	data, err := ioutil.ReadFile("./params/MultiTransferOng.json")
	if err != nil {
		log4.Error("ioutil.ReadFile failed ", err)
		return false
	}
	multiTransferParam := new(MultiTransferParam)
	err = json.Unmarshal(data, multiTransferParam)
	if err != nil {
		log4.Error("json.Unmarshal failed ", err)
		return false
	}
	var users []*sdk.Account
	time.Sleep(1 * time.Second)
	for _, path := range multiTransferParam.FromPath {
		user, ok := common.GetAccountByPassword(ontSdk, path)
		if !ok {
			return false
		}
		users = append(users, user)
	}
	ok := multiTransfer(ontSdk, utils.OngContractAddress, users, multiTransferParam.ToAddress, multiTransferParam.Amount)
	if !ok {
		return false
	}
	common.WaitForBlock(ontSdk)
	return true
}

type GetAttributesParam struct {
	PeerPubkey string
}

func GetAttributes(ontSdk *sdk.OntologySdk) bool {
	data, err := ioutil.ReadFile("./params/GetAttributes.json")
	if err != nil {
		log4.Error("ioutil.ReadFile failed ", err)
		return false
	}
	getAttributesParam := new(GetAttributesParam)
	err = json.Unmarshal(data, getAttributesParam)
	if err != nil {
		log4.Error("json.Unmarshal failed ", err)
		return false
	}
	peerAttributes, err := getAttributes(ontSdk, getAttributesParam.PeerPubkey)
	if err != nil {
		log4.Error("getAttributes failed ", err)
		return false
	}
	fmt.Println("peerAttributes.PeerPubkey is:", peerAttributes.PeerPubkey)
	fmt.Println("peerAttributes.MaxAuthorize is:", peerAttributes.MaxAuthorize)
	fmt.Println("peerAttributes.T2PeerCost is:", peerAttributes.T2PeerCost)
	fmt.Println("peerAttributes.T1PeerCost is:", peerAttributes.T1PeerCost)
	fmt.Println("peerAttributes.TPeerCost is:", peerAttributes.TPeerCost)

	return true
}

type GetSplitFeeAddressParam struct {
	Address string
}

func GetSplitFeeAddress(ontSdk *sdk.OntologySdk) bool {
	data, err := ioutil.ReadFile("./params/GetSplitFeeAddress.json")
	if err != nil {
		log4.Error("ioutil.ReadFile failed ", err)
		return false
	}
	getSplitFeeAddressParam := new(GetSplitFeeAddressParam)
	err = json.Unmarshal(data, getSplitFeeAddressParam)
	if err != nil {
		log4.Error("json.Unmarshal failed ", err)
		return false
	}
	address, err := ocommon.AddressFromBase58(getSplitFeeAddressParam.Address)
	if err != nil {
		log4.Error("common.AddressFromBase58 failed ", err)
		return false
	}
	splitFeeAddress, err := getSplitFeeAddress(ontSdk, address)
	if err != nil {
		log4.Error("getSplitFeeAddress failed ", err)
		return false
	}
	fmt.Println("splitFeeAddress.Address is:", splitFeeAddress.Address)
	fmt.Println("splitFeeAddress.Amount is:", splitFeeAddress.Amount)

	return true
}

func GetSplitFee(ontSdk *sdk.OntologySdk) bool {
	splitFee, err := getSplitFee(ontSdk)
	if err != nil {
		log4.Error("getSplitFeeAddress failed ", err)
		return false
	}
	fmt.Println("splitFee is:", splitFee)

	return true
}

type GetPromisePosParam struct {
	PeerPubkey string
}

func GetPromisePos(ontSdk *sdk.OntologySdk) bool {
	data, err := ioutil.ReadFile("./params/GetPromisePos.json")
	if err != nil {
		log4.Error("ioutil.ReadFile failed ", err)
		return false
	}
	getPromisePosParam := new(GetPromisePosParam)
	err = json.Unmarshal(data, getPromisePosParam)
	if err != nil {
		log4.Error("json.Unmarshal failed ", err)
		return false
	}
	promisePos, err := getPromisePos(ontSdk, getPromisePosParam.PeerPubkey)
	if err != nil {
		log4.Error("getPromisePos failed ", err)
		return false
	}
	fmt.Println("promisePos.PeerPubkey is:", promisePos.PeerPubkey)
	fmt.Println("promisePos.PromisePos is:", promisePos.PromisePos)

	return true
}
