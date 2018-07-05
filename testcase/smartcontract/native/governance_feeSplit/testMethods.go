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
	"time"

	"github.com/ontio/ontology-crypto/keypair"
	"github.com/ontio/ontology-crypto/vrf"
	"github.com/ontio/ontology-tool/testframework"
	"github.com/ontio/ontology/account"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/core/types"
	"github.com/ontio/ontology/smartcontract/service/native/governance"
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
	ok = assignFuncsToRole(ctx, user)
	if !ok {
		return false
	}
	waitForBlock(ctx)
	return true
}

type AssignOntIDsToRoleParam struct {
	Path1 string
	Path2 string
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
	user2, ok := getAccountByPassword(ctx, assignOntIDsToRoleParam.Path2)
	if !ok {
		return false
	}
	ok = assignOntIDsToRole(ctx, user1, []*account.Account{user2})
	if !ok {
		return false
	}
	waitForBlock(ctx)
	return true
}

type RegisterCandidateParam struct {
	Path       []string
	PeerPubkey []string
	InitPos    []uint32
}

func RegisterCandidate(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./params/RegisterCandidate.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	registerCandidateParam := new(RegisterCandidateParam)
	err = json.Unmarshal(data, registerCandidateParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}
	for i := 0; i < len(registerCandidateParam.PeerPubkey); i++ {
		user, ok := getAccountByPassword(ctx, registerCandidateParam.Path[i])
		if !ok {
			return false
		}
		ok = registerCandidate(ctx, user, registerCandidateParam.PeerPubkey[i], registerCandidateParam.InitPos[i])
		if !ok {
			return false
		}
	}
	waitForBlock(ctx)
	return true
}

type UnRegisterCandidateParam struct {
	Path       string
	PeerPubkey []string
}

func UnRegisterCandidate(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./params/UnRegisterCandidate.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	unRegisterCandidateParam := new(UnRegisterCandidateParam)
	err = json.Unmarshal(data, unRegisterCandidateParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}
	user, ok := getAccount(ctx, unRegisterCandidateParam.Path)
	if !ok {
		return false
	}
	for i := 0; i < len(unRegisterCandidateParam.PeerPubkey); i++ {
		ok = unRegisterCandidate(ctx, user, unRegisterCandidateParam.PeerPubkey[i])
		if !ok {
			return false
		}
	}
	waitForBlock(ctx)
	return true
}

type ApproveCandidateParam struct {
	Path       []string
	PeerPubkey []string
}

func ApproveCandidate(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./params/ApproveCandidate.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	approveCandidateParam := new(ApproveCandidateParam)
	err = json.Unmarshal(data, approveCandidateParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}
	var users []*account.Account
	var pubKeys []keypair.PublicKey
	time.Sleep(1 * time.Second)
	for _, path := range approveCandidateParam.Path {
		user, ok := getAccountByPassword(ctx, path)
		if !ok {
			return false
		}
		users = append(users, user)
		pubKeys = append(pubKeys, user.PublicKey)
	}
	for _, peerPubkey := range approveCandidateParam.PeerPubkey {
		ok := approveCandidateMultiSign(ctx, pubKeys, users, peerPubkey)
		if !ok {
			return false
		}
	}
	waitForBlock(ctx)
	return true
}

type RejectCandidateParam struct {
	Path       []string
	PeerPubkey string
}

func RejectCandidate(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./params/RejectCandidate.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	rejectCandidateParam := new(RejectCandidateParam)
	err = json.Unmarshal(data, rejectCandidateParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}
	var users []*account.Account
	var pubKeys []keypair.PublicKey
	time.Sleep(1 * time.Second)
	for _, path := range rejectCandidateParam.Path {
		user, ok := getAccountByPassword(ctx, path)
		if !ok {
			return false
		}
		users = append(users, user)
		pubKeys = append(pubKeys, user.PublicKey)
	}
	ok := rejectCandidateMultiSign(ctx, pubKeys, users, rejectCandidateParam.PeerPubkey)
	if !ok {
		return false
	}
	waitForBlock(ctx)
	return true
}

type VoteForPeerParam struct {
	Path           string
	PeerPubkeyList []string
	PosList        []uint32
}

func VoteForPeer(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./params/VoteForPeer.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	voteForPeerParam := new(VoteForPeerParam)
	err = json.Unmarshal(data, voteForPeerParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}
	user, ok := getAccount(ctx, voteForPeerParam.Path)
	if !ok {
		return false
	}
	ok = voteForPeer(ctx, user, voteForPeerParam.PeerPubkeyList, voteForPeerParam.PosList)
	if !ok {
		return false
	}
	waitForBlock(ctx)
	return true
}

func UnVoteForPeer(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./params/UnVoteForPeer.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	voteForPeerParam := new(VoteForPeerParam)
	err = json.Unmarshal(data, voteForPeerParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}
	user, ok := getAccount(ctx, voteForPeerParam.Path)
	if !ok {
		return false
	}
	ok = unVoteForPeer(ctx, user, voteForPeerParam.PeerPubkeyList, voteForPeerParam.PosList)
	if !ok {
		return false
	}
	waitForBlock(ctx)
	return true
}

type WithdrawParam struct {
	Path           string
	PeerPubkeyList []string
	WithdrawList   []uint32
}

func Withdraw(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./params/Withdraw.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	withdrawParam := new(WithdrawParam)
	err = json.Unmarshal(data, withdrawParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}
	user, ok := getAccount(ctx, withdrawParam.Path)
	if !ok {
		return false
	}
	ok = withdraw(ctx, user, withdrawParam.PeerPubkeyList, withdrawParam.WithdrawList)
	if !ok {
		return false
	}
	waitForBlock(ctx)
	return true
}

type QuitNodeParam struct {
	Path       []string
	PeerPubkey []string
}

func QuitNode(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./params/QuitNode.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	quitNodeParam := new(QuitNodeParam)
	err = json.Unmarshal(data, quitNodeParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}
	time.Sleep(1 * time.Second)
	for i := 0; i < len(quitNodeParam.Path); i++ {
		user, ok := getAccountByPassword(ctx, quitNodeParam.Path[i])
		if !ok {
			return false
		}
		ok = quitNode(ctx, user, quitNodeParam.PeerPubkey[i])
		if !ok {
			return false
		}
	}
	waitForBlock(ctx)
	return true
}

type BlackNodeParam struct {
	Path           []string
	PeerPubkeyList []string
}

func BlackNode(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./params/BlackNode.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	blackNodeParam := new(BlackNodeParam)
	err = json.Unmarshal(data, blackNodeParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}
	var users []*account.Account
	var pubKeys []keypair.PublicKey
	time.Sleep(1 * time.Second)
	for _, path := range blackNodeParam.Path {
		user, ok := getAccountByPassword(ctx, path)
		if !ok {
			return false
		}
		users = append(users, user)
		pubKeys = append(pubKeys, user.PublicKey)
	}
	ok := blackNodeMultiSign(ctx, pubKeys, users, blackNodeParam.PeerPubkeyList)
	if !ok {
		return false
	}
	waitForBlock(ctx)
	return true
}

type WhiteNodeParam struct {
	Path       []string
	PeerPubkey string
}

func WhiteNode(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./params/WhiteNode.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	whiteNodeParam := new(WhiteNodeParam)
	err = json.Unmarshal(data, whiteNodeParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}
	var users []*account.Account
	var pubKeys []keypair.PublicKey
	time.Sleep(1 * time.Second)
	for _, path := range whiteNodeParam.Path {
		user, ok := getAccountByPassword(ctx, path)
		if !ok {
			return false
		}
		users = append(users, user)
		pubKeys = append(pubKeys, user.PublicKey)
	}
	ok := whiteNodeMultiSign(ctx, pubKeys, users, whiteNodeParam.PeerPubkey)
	if !ok {
		return false
	}
	waitForBlock(ctx)
	return true
}

type MultiAccount struct {
	Path []string
}

func CommitDpos(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./params/CommitDpos.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	multiAccount := new(MultiAccount)
	err = json.Unmarshal(data, multiAccount)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}
	var users []*account.Account
	var pubKeys []keypair.PublicKey
	time.Sleep(1 * time.Second)
	for _, path := range multiAccount.Path {
		user, ok := getAccountByPassword(ctx, path)
		if !ok {
			return false
		}
		users = append(users, user)
		pubKeys = append(pubKeys, user.PublicKey)
	}
	ok := commitDposMultiSign(ctx, pubKeys, users)
	if !ok {
		return false
	}
	waitForBlock(ctx)
	return true
}

func CallSplit(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./params/CallSplit.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	multiAccount := new(MultiAccount)
	err = json.Unmarshal(data, multiAccount)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}
	var users []*account.Account
	var pubKeys []keypair.PublicKey
	time.Sleep(1 * time.Second)
	for _, path := range multiAccount.Path {
		user, ok := getAccountByPassword(ctx, path)
		if !ok {
			return false
		}
		users = append(users, user)
		pubKeys = append(pubKeys, user.PublicKey)
	}
	ok := callSplitMultiSign(ctx, pubKeys, users)
	if !ok {
		return false
	}
	waitForBlock(ctx)
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

func UpdateConfig(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./params/UpdateConfig.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	updateConfigParam := new(UpdateConfigParam)
	err = json.Unmarshal(data, updateConfigParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}
	var users []*account.Account
	var pubKeys []keypair.PublicKey
	time.Sleep(1 * time.Second)
	for _, path := range updateConfigParam.Path {
		user, ok := getAccountByPassword(ctx, path)
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
	ok := updateConfigMultiSign(ctx, pubKeys, users, config)
	if !ok {
		return false
	}
	waitForBlock(ctx)
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

func UpdateGlobalParam(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./params/UpdateGlobalParam.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	updateGlobalParamParam := new(UpdateGlobalParamParam)
	err = json.Unmarshal(data, updateGlobalParamParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}
	var users []*account.Account
	var pubKeys []keypair.PublicKey
	time.Sleep(1 * time.Second)
	for _, path := range updateGlobalParamParam.Path {
		user, ok := getAccountByPassword(ctx, path)
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
	ok := updateGlobalParamMultiSign(ctx, pubKeys, users, globalParam)
	if !ok {
		return false
	}
	waitForBlock(ctx)
	return true
}

type UpdateSplitCurveParam struct {
	Path []string
	Yi   []uint32
}

func UpdateSplitCurve(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./params/UpdateSplitCurve.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	updateSplitCurveParam := new(UpdateSplitCurveParam)
	err = json.Unmarshal(data, updateSplitCurveParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}
	var users []*account.Account
	var pubKeys []keypair.PublicKey
	time.Sleep(1 * time.Second)
	for _, path := range updateSplitCurveParam.Path {
		user, ok := getAccountByPassword(ctx, path)
		if !ok {
			return false
		}
		users = append(users, user)
		pubKeys = append(pubKeys, user.PublicKey)
	}
	splitCurve := &governance.SplitCurve{
		Yi: updateSplitCurveParam.Yi,
	}
	ok := updateSplitCurveMultiSign(ctx, pubKeys, users, splitCurve)
	if !ok {
		return false
	}
	waitForBlock(ctx)
	return true
}

type TransferPenaltyParam struct {
	Path1      []string
	PeerPubkey string
	Path2      string
}

func TransferPenalty(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./params/TransferPenalty.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	transferPenaltyParam := new(TransferPenaltyParam)
	err = json.Unmarshal(data, transferPenaltyParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}
	var users []*account.Account
	var pubKeys []keypair.PublicKey
	time.Sleep(1 * time.Second)
	for _, path := range transferPenaltyParam.Path1 {
		user, ok := getAccountByPassword(ctx, path)
		if !ok {
			return false
		}
		users = append(users, user)
		pubKeys = append(pubKeys, user.PublicKey)
	}
	user1, ok := getAccount(ctx, transferPenaltyParam.Path2)
	if !ok {
		return false
	}
	ok = transferPenaltyMultiSign(ctx, pubKeys, users, transferPenaltyParam.PeerPubkey, user1.Address)
	if !ok {
		return false
	}
	waitForBlock(ctx)
	return true
}

func GetVbftConfig(ctx *testframework.TestFrameworkContext) bool {
	config, err := getVbftConfig(ctx)
	if err != nil {
		ctx.LogError("getVbftConfig failed %v", err)
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

func GetGlobalParam(ctx *testframework.TestFrameworkContext) bool {
	globalParam, err := getGlobalParam(ctx)
	if err != nil {
		ctx.LogError("getGlobalParam failed %v", err)
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

func GetSplitCurve(ctx *testframework.TestFrameworkContext) bool {
	splitCurve, err := getSplitCurve(ctx)
	if err != nil {
		ctx.LogError("getSplitCurve failed %v", err)
		return false
	}
	fmt.Println("splitCurve.Yi is", splitCurve.Yi)
	return true
}

func GetGovernanceView(ctx *testframework.TestFrameworkContext) bool {
	governanceView, err := getGovernanceView(ctx)
	if err != nil {
		ctx.LogError("getGovernanceView failed %v", err)
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

func GetPeerPoolItem(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./params/GetPeerPoolItem.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	getPeerPoolItemParam := new(GetPeerPoolItemParam)
	err = json.Unmarshal(data, getPeerPoolItemParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}

	peerPoolMap, err := getPeerPoolMap(ctx)
	if err != nil {
		ctx.LogError("getPeerPoolMap failed %v", err)
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

type GetVoteInfoParam struct {
	Path       string
	PeerPubkey string
}

func GetVoteInfo(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./params/GetVoteInfo.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	getVoteInfoParam := new(GetVoteInfoParam)
	err = json.Unmarshal(data, getVoteInfoParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}
	user, ok := getAccount(ctx, getVoteInfoParam.Path)
	if !ok {
		return false
	}

	voteInfo, err := getVoteInfo(ctx, getVoteInfoParam.PeerPubkey, user.Address)
	if err != nil {
		ctx.LogError("getVoteInfo failed %v", err)
		return false
	}

	fmt.Println("voteInfo.PeerPubkey is:", voteInfo.PeerPubkey)
	fmt.Println("voteInfo.Address is:", voteInfo.Address.ToBase58())
	fmt.Println("voteInfo.ConsensusPos is:", voteInfo.ConsensusPos)
	fmt.Println("voteInfo.FreezePos is:", voteInfo.FreezePos)
	fmt.Println("voteInfo.NewPos is:", voteInfo.NewPos)
	fmt.Println("voteInfo.WithdrawPos is:", voteInfo.WithdrawPos)
	fmt.Println("voteInfo.WithdrawFreezePos is:", voteInfo.WithdrawFreezePos)
	fmt.Println("voteInfo.WithdrawUnfreezePos is:", voteInfo.WithdrawUnfreezePos)
	return true
}

type GetTotalStakeParam struct {
	Path string
}

func GetTotalStake(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./params/GetTotalStake.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	getTotalStakeParam := new(GetTotalStakeParam)
	err = json.Unmarshal(data, getTotalStakeParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}
	user, ok := getAccount(ctx, getTotalStakeParam.Path)
	if !ok {
		return false
	}

	totalStake, err := getTotalStake(ctx, user.Address)
	if err != nil {
		ctx.LogError("getTotalStake failed %v", err)
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

func GetPenaltyStake(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./params/GetPenaltyStake.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	getPenaltyStakeParam := new(GetPenaltyStakeParam)
	err = json.Unmarshal(data, getPenaltyStakeParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}

	penaltyStake, err := getPenaltyStake(ctx, getPenaltyStakeParam.PeerPubkey)
	if err != nil {
		ctx.LogError("getPenaltyStake failed %v", err)
		return false
	}

	fmt.Println("penaltyStake.PeerPubkey is:", penaltyStake.PeerPubkey)
	fmt.Println("penaltyStake.InitPos is:", penaltyStake.InitPos)
	fmt.Println("penaltyStake.VotePos is:", penaltyStake.VotePos)
	fmt.Println("penaltyStake.TimeOffset is:", penaltyStake.TimeOffset)
	fmt.Println("penaltyStake.Amount is:", penaltyStake.Amount)
	return true
}

type InBlackListParam struct {
	PeerPubkey string
}

func InBlackList(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./params/InBlackList.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	inBlackListParam := new(InBlackListParam)
	err = json.Unmarshal(data, inBlackListParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}

	inBlackList, err := inBlackList(ctx, inBlackListParam.PeerPubkey)
	if err != nil {
		ctx.LogError("getPenaltyStake failed %v", err)
		return false
	}

	fmt.Println("result is:", inBlackList)
	return true
}

type WithdrawOngParam struct {
	Path       string
	PeerPubkey string
}

func WithdrawOng(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./params/WithdrawOng.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	withdrawOngParam := new(WithdrawOngParam)
	err = json.Unmarshal(data, withdrawOngParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}
	user, ok := getAccount(ctx, withdrawOngParam.Path)
	if !ok {
		return false
	}
	ok = withdrawOng(ctx, user)
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

func TransferOntMultiSign(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./params/TransferOntMultiSign.json")
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
	var users []*account.Account
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
		ok = transferOntMultiSign(ctx, pubKeys, users, user2.Address, transferMultiSignParam.Amount[index])
		if !ok {
			return false
		}
	}
	waitForBlock(ctx)
	return true
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
	var users []*account.Account
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

type TransferFromMultiSignParam struct {
	Path1  []string
	Path2  []string
	Amount []uint64
}

func TransferFromOngMultiSign(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./params/TransferFromOngMultiSign.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	transferFromMultiSignParam := new(TransferFromMultiSignParam)
	err = json.Unmarshal(data, transferFromMultiSignParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}
	var users []*account.Account
	var pubKeys []keypair.PublicKey
	time.Sleep(1 * time.Second)
	for _, path := range transferFromMultiSignParam.Path1 {
		user, ok := getAccountByPassword(ctx, path)
		if !ok {
			return false
		}
		users = append(users, user)
		pubKeys = append(pubKeys, user.PublicKey)
	}
	time.Sleep(1 * time.Second)
	for index, path2 := range transferFromMultiSignParam.Path2 {
		user2, ok := getAccountByPassword(ctx, path2)
		if !ok {
			return false
		}
		ok = transferFromOngMultiSign(ctx, pubKeys, users, user2.Address, transferFromMultiSignParam.Amount[index])
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

type TransferMultiSignToMultiSignParam struct {
	Path1   []string
	PubKeys []string
	Amount  uint64
}

func TransferOntMultiSignToMultiSign(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./params/TransferOntMultiSignToMultiSign.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	transferMultiSignToMultiSignParam := new(TransferMultiSignToMultiSignParam)
	err = json.Unmarshal(data, transferMultiSignToMultiSignParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}
	var users []*account.Account
	var pubKeys []keypair.PublicKey
	var pubKeysTo []keypair.PublicKey
	time.Sleep(1 * time.Second)
	for _, path := range transferMultiSignToMultiSignParam.Path1 {
		user, ok := getAccountByPassword(ctx, path)
		if !ok {
			return false
		}
		users = append(users, user)
		pubKeys = append(pubKeys, user.PublicKey)
	}
	for _, v := range transferMultiSignToMultiSignParam.PubKeys {
		vByte, err := hex.DecodeString(v)
		if err != nil {
			ctx.LogError("hex.DecodeString failed %v", err)
		}
		k, err := keypair.DeserializePublicKey(vByte)
		if err != nil {
			ctx.LogError("keypair.DeserializePublicKey failed %v", err)
		}
		pubKeysTo = append(pubKeysTo, k)
	}
	to, err := types.AddressFromMultiPubKeys(pubKeysTo, int((5*len(pubKeysTo)+6)/7))
	if err != nil {
		ctx.LogError("types.AddressFromMultiPubKeys error", err)
	}
	ok := transferOntMultiSignToMultiSign(ctx, pubKeys, users, to, transferMultiSignToMultiSignParam.Amount)
	if !ok {
		return false
	}
	waitForBlock(ctx)
	return true
}

func TransferOngMultiSignToMultiSign(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./params/TransferOngMultiSignToMultiSign.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	transferMultiSignToMultiSignParam := new(TransferMultiSignToMultiSignParam)
	err = json.Unmarshal(data, transferMultiSignToMultiSignParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}
	var users []*account.Account
	var pubKeys []keypair.PublicKey
	var pubKeysTo []keypair.PublicKey
	time.Sleep(1 * time.Second)
	for _, path := range transferMultiSignToMultiSignParam.Path1 {
		user, ok := getAccountByPassword(ctx, path)
		if !ok {
			return false
		}
		users = append(users, user)
		pubKeys = append(pubKeys, user.PublicKey)
	}
	for _, v := range transferMultiSignToMultiSignParam.PubKeys {
		vByte, err := hex.DecodeString(v)
		if err != nil {
			ctx.LogError("hex.DecodeString failed %v", err)
		}
		k, err := keypair.DeserializePublicKey(vByte)
		if err != nil {
			ctx.LogError("keypair.DeserializePublicKey failed %v", err)
		}
		pubKeysTo = append(pubKeysTo, k)
	}
	to, err := types.AddressFromMultiPubKeys(pubKeysTo, int((5*len(pubKeysTo)+6)/7))
	if err != nil {
		ctx.LogError("types.AddressFromMultiPubKeys error", err)
	}
	ok := transferOngMultiSignToMultiSign(ctx, pubKeys, users, to, transferMultiSignToMultiSignParam.Amount)
	if !ok {
		return false
	}
	waitForBlock(ctx)
	return true
}

type TransferFromMultiSignToMultiSignParam struct {
	Path1   []string
	PubKeys []string
	Amount  uint64
}

func TransferFromOngMultiSignToMultiSign(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./params/TransferFromOngMultiSignToMultiSign.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	transferFromMultiSignToMultiSignParam := new(TransferFromMultiSignToMultiSignParam)
	err = json.Unmarshal(data, transferFromMultiSignToMultiSignParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}
	var users []*account.Account
	var pubKeys []keypair.PublicKey
	var pubKeysTo []keypair.PublicKey
	time.Sleep(1 * time.Second)
	for _, path := range transferFromMultiSignToMultiSignParam.Path1 {
		user, ok := getAccountByPassword(ctx, path)
		if !ok {
			return false
		}
		users = append(users, user)
		pubKeys = append(pubKeys, user.PublicKey)
	}
	for _, v := range transferFromMultiSignToMultiSignParam.PubKeys {
		vByte, err := hex.DecodeString(v)
		if err != nil {
			ctx.LogError("hex.DecodeString failed %v", err)
		}
		k, err := keypair.DeserializePublicKey(vByte)
		if err != nil {
			ctx.LogError("keypair.DeserializePublicKey failed %v", err)
		}
		pubKeysTo = append(pubKeysTo, k)
	}
	to, err := types.AddressFromMultiPubKeys(pubKeysTo, int((5*len(pubKeysTo)+6)/7))
	if err != nil {
		ctx.LogError("types.AddressFromMultiPubKeys error", err)
	}
	ok := transferFromOngMultiSignToMultiSign(ctx, pubKeys, users, to, transferFromMultiSignToMultiSignParam.Amount)
	if !ok {
		return false
	}
	waitForBlock(ctx)
	return true
}

type TransferMultiSignAddressParam struct {
	Path1   []string
	Address []string
	Amount  []uint64
}

func TransferOntMultiSignAddress(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./params/TransferOntMultiSignAddress.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	transferMultiSignAddressParam := new(TransferMultiSignAddressParam)
	err = json.Unmarshal(data, transferMultiSignAddressParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}
	var users []*account.Account
	var pubKeys []keypair.PublicKey
	time.Sleep(1 * time.Second)
	for _, path := range transferMultiSignAddressParam.Path1 {
		user, ok := getAccountByPassword(ctx, path)
		if !ok {
			return false
		}
		users = append(users, user)
		pubKeys = append(pubKeys, user.PublicKey)
	}
	time.Sleep(1 * time.Second)
	for index, address := range transferMultiSignAddressParam.Address {
		addr, err := common.AddressFromBase58(address)
		if err != nil {
			ctx.LogError("common.AddressFromBase58 failed %v", err)
			return false
		}
		ok := transferOntMultiSign(ctx, pubKeys, users, addr, transferMultiSignAddressParam.Amount[index])
		if !ok {
			return false
		}
	}
	waitForBlock(ctx)
	return true
}

func TransferOngMultiSignAddress(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./params/TransferOngMultiSignAddress.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	transferMultiSignAddressParam := new(TransferMultiSignAddressParam)
	err = json.Unmarshal(data, transferMultiSignAddressParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}
	var users []*account.Account
	var pubKeys []keypair.PublicKey
	time.Sleep(1 * time.Second)
	for _, path := range transferMultiSignAddressParam.Path1 {
		user, ok := getAccountByPassword(ctx, path)
		if !ok {
			return false
		}
		users = append(users, user)
		pubKeys = append(pubKeys, user.PublicKey)
	}
	time.Sleep(1 * time.Second)
	for index, address := range transferMultiSignAddressParam.Address {
		addr, err := common.AddressFromBase58(address)
		if err != nil {
			ctx.LogError("common.AddressFromBase58 failed %v", err)
			return false
		}
		ok := transferOngMultiSign(ctx, pubKeys, users, addr, transferMultiSignAddressParam.Amount[index])
		if !ok {
			return false
		}
	}
	waitForBlock(ctx)
	return true
}

type TransferFromMultiSignAddressParam struct {
	Path1   []string
	Address []string
	Amount  []uint64
}

func TransferFromOngMultiSignAddress(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./params/TransferFromOngMultiSignAddress.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	transferFromMultiSignAddressParam := new(TransferFromMultiSignAddressParam)
	err = json.Unmarshal(data, transferFromMultiSignAddressParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}
	var users []*account.Account
	var pubKeys []keypair.PublicKey
	time.Sleep(1 * time.Second)
	for _, path := range transferFromMultiSignAddressParam.Path1 {
		user, ok := getAccountByPassword(ctx, path)
		if !ok {
			return false
		}
		users = append(users, user)
		pubKeys = append(pubKeys, user.PublicKey)
	}
	time.Sleep(1 * time.Second)
	for index, address := range transferFromMultiSignAddressParam.Address {
		addr, err := common.AddressFromBase58(address)
		if err != nil {
			ctx.LogError("common.AddressFromBase58 failed %v", err)
			return false
		}
		ok := transferFromOngMultiSign(ctx, pubKeys, users, addr, transferFromMultiSignAddressParam.Amount[index])
		if !ok {
			return false
		}
	}
	waitForBlock(ctx)
	return true
}
