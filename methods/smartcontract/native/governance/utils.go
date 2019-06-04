package governance

import (
	"bytes"
	"encoding/hex"

	log4 "github.com/alecthomas/log4go"
	"github.com/ontio/ontology-crypto/keypair"
	sdk "github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology-tool/common"
	"github.com/ontio/ontology-tool/config"
	ontcommon "github.com/ontio/ontology/common"
	"github.com/ontio/ontology/common/serialization"
	"github.com/ontio/ontology/core/types"
	"github.com/ontio/ontology/errors"
	"github.com/ontio/ontology/smartcontract/service/native/auth"
	"github.com/ontio/ontology/smartcontract/service/native/governance"
	"github.com/ontio/ontology/smartcontract/service/native/ont"
	"github.com/ontio/ontology/smartcontract/service/native/utils"
)

var OntIDVersion = byte(0)

const PROMISE_POS = 200000

func registerCandidate(ontSdk *sdk.OntologySdk, user *sdk.Account, peerPubkey string, initPos uint32) bool {
	params := &governance.RegisterCandidateParam{
		PeerPubkey: peerPubkey,
		Address:    user.Address,
		InitPos:    initPos,
		Caller:     []byte("did:ont:" + user.Address.ToBase58()),
		KeyNo:      1,
	}
	method := "registerCandidate"
	contractAddress := utils.GovernanceContractAddress
	txHash, err := ontSdk.Native.InvokeNativeContract(config.DefConfig.GasPrice, config.DefConfig.GasLimit, user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		log4.Error("invokeNativeContract error :", err)
		return false
	}
	log4.Info("registerCandidate txHash is :", txHash.ToHexString())
	//common.WaitForBlock(ontSdk)

	//let node can be authorized
	//changeMaxAuthorization(ontSdk, user, peerPubkey, PROMISE_POS)
	return true
}

func registerCandidate2Sign(ontSdk *sdk.OntologySdk, ontid *sdk.Account, user *sdk.Account, peerPubkey string, initPos uint32) bool {
	params := &governance.RegisterCandidateParam{
		PeerPubkey: peerPubkey,
		Address:    user.Address,
		InitPos:    initPos,
		Caller:     []byte("did:ont:" + ontid.Address.ToBase58()),
		KeyNo:      1,
	}
	method := "registerCandidate"
	contractAddress := utils.GovernanceContractAddress
	tx, err := ontSdk.Native.NewNativeInvokeTransaction(config.DefConfig.GasPrice, config.DefConfig.GasLimit, OntIDVersion, contractAddress, method, []interface{}{params})
	if err != nil {
		log4.Error("NewNativeInvokeTransaction error")
		return false
	}
	err = ontSdk.SignToTransaction(tx, user)
	if err != nil {
		log4.Error("SignToTransaction error")
		return false
	}
	err = ontSdk.SignToTransaction(tx, ontid)
	if err != nil {
		log4.Error("SignToTransaction error")
		return false
	}
	txHash, err := ontSdk.SendTransaction(tx)
	if err != nil {
		log4.Error("SendRawTransaction error", err)
		return false
	}
	log4.Info("registerCandidate2Sign txHash is :", txHash.ToHexString())
	return true
}

func unRegisterCandidate(ontSdk *sdk.OntologySdk, user *sdk.Account, peerPubkey string) bool {
	params := &governance.UnRegisterCandidateParam{
		PeerPubkey: peerPubkey,
		Address:    user.Address,
	}
	method := "unRegisterCandidate"
	contractAddress := utils.GovernanceContractAddress
	txHash, err := ontSdk.Native.InvokeNativeContract(config.DefConfig.GasPrice, config.DefConfig.GasLimit, user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		log4.Error("invokeNativeContract error :", err)
		return false
	}
	log4.Info("unRegisterCandidate txHash is :", txHash.ToHexString())
	return true
}

func approveCandidate(ontSdk *sdk.OntologySdk, user *sdk.Account, peerPubkey string) bool {
	params := &governance.ApproveCandidateParam{
		PeerPubkey: peerPubkey,
	}
	contractAddress := utils.GovernanceContractAddress
	method := "approveCandidate"
	txHash, err := ontSdk.Native.InvokeNativeContract(config.DefConfig.GasPrice, config.DefConfig.GasLimit, user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		log4.Error("invokeNativeContract error :", err)
		return false
	}
	log4.Info("approveCandidate txHash is :", txHash.ToHexString())
	return true
}

func approveCandidateMultiSign(ontSdk *sdk.OntologySdk, pubKeys []keypair.PublicKey, users []*sdk.Account, peerPubkey string) bool {
	params := &governance.ApproveCandidateParam{
		PeerPubkey: peerPubkey,
	}
	contractAddress := utils.GovernanceContractAddress
	method := "approveCandidate"
	txHash, err := common.InvokeNativeContractWithMultiSign(ontSdk, config.DefConfig.GasPrice, config.DefConfig.GasLimit, pubKeys, users, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		log4.Error("invokeNativeContract error :", err)
		return false
	}
	log4.Info("approveCandidateMultiSign txHash is :", txHash.ToHexString())
	return true
}

func rejectCandidate(ontSdk *sdk.OntologySdk, user *sdk.Account, peerPubkey string) bool {
	params := &governance.RejectCandidateParam{
		PeerPubkey: peerPubkey,
	}
	contractAddress := utils.GovernanceContractAddress
	method := "rejectCandidate"
	txHash, err := ontSdk.Native.InvokeNativeContract(config.DefConfig.GasPrice, config.DefConfig.GasLimit, user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		log4.Error("invokeNativeContract error :", err)
		return false
	}
	log4.Info("rejectCandidate txHash is :", txHash.ToHexString())
	return true
}

func rejectCandidateMultiSign(ontSdk *sdk.OntologySdk, pubKeys []keypair.PublicKey, users []*sdk.Account, peerPubkey string) bool {
	params := &governance.RejectCandidateParam{
		PeerPubkey: peerPubkey,
	}
	contractAddress := utils.GovernanceContractAddress
	method := "rejectCandidate"
	txHash, err := common.InvokeNativeContractWithMultiSign(ontSdk, config.DefConfig.GasPrice, config.DefConfig.GasLimit, pubKeys, users, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		log4.Error("invokeNativeContract error :", err)
		return false
	}
	log4.Info("rejectCandidateMultiSign txHash is :", txHash.ToHexString())
	return true
}

func changeMaxAuthorization(ontSdk *sdk.OntologySdk, user *sdk.Account, peerPubkey string, maxAuthorize uint32) bool {
	params := &governance.ChangeMaxAuthorizationParam{
		Address:      user.Address,
		PeerPubkey:   peerPubkey,
		MaxAuthorize: maxAuthorize,
	}
	contractAddress := utils.GovernanceContractAddress
	method := "changeMaxAuthorization"
	txHash, err := ontSdk.Native.InvokeNativeContract(config.DefConfig.GasPrice, config.DefConfig.GasLimit, user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		log4.Error("invokeNativeContract error :", err)
		return false
	}
	log4.Info("changeMaxAuthorization txHash is :", txHash.ToHexString())
	return true
}

func setPeerCost(ontSdk *sdk.OntologySdk, user *sdk.Account, peerPubkey string, peerCost uint32) bool {
	params := &governance.SetPeerCostParam{
		Address:    user.Address,
		PeerPubkey: peerPubkey,
		PeerCost:   peerCost,
	}
	contractAddress := utils.GovernanceContractAddress
	method := "setPeerCost"
	txHash, err := ontSdk.Native.InvokeNativeContract(config.DefConfig.GasPrice, config.DefConfig.GasLimit, user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		log4.Error("invokeNativeContract error :", err)
		return false
	}
	log4.Info("setPeerCost txHash is :", txHash.ToHexString())
	return true
}

func addInitPos(ontSdk *sdk.OntologySdk, user *sdk.Account, peerPubkey string, pos uint32) bool {
	params := &governance.ChangeInitPosParam{
		Address:    user.Address,
		PeerPubkey: peerPubkey,
		Pos:        pos,
	}
	contractAddress := utils.GovernanceContractAddress
	method := "addInitPos"
	txHash, err := ontSdk.Native.InvokeNativeContract(config.DefConfig.GasPrice, config.DefConfig.GasLimit, user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		log4.Error("invokeNativeContract error :", err)
		return false
	}
	log4.Info("addInitPos txHash is :", txHash.ToHexString())
	return true
}

func reduceInitPos(ontSdk *sdk.OntologySdk, user *sdk.Account, peerPubkey string, pos uint32) bool {
	params := &governance.ChangeInitPosParam{
		Address:    user.Address,
		PeerPubkey: peerPubkey,
		Pos:        pos,
	}
	contractAddress := utils.GovernanceContractAddress
	method := "reduceInitPos"
	txHash, err := ontSdk.Native.InvokeNativeContract(config.DefConfig.GasPrice, config.DefConfig.GasLimit, user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		log4.Error("invokeNativeContract error :", err)
		return false
	}
	log4.Info("reduceInitPos txHash is :", txHash.ToHexString())
	return true
}

func authorizeForPeer(ontSdk *sdk.OntologySdk, user *sdk.Account, peerPubkeyList []string, posList []uint32) bool {
	params := &governance.AuthorizeForPeerParam{
		Address:        user.Address,
		PeerPubkeyList: peerPubkeyList,
		PosList:        posList,
	}
	contractAddress := utils.GovernanceContractAddress
	method := "authorizeForPeer"
	_, err := ontSdk.Native.InvokeNativeContract(config.DefConfig.GasPrice, config.DefConfig.GasLimit, user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		log4.Error("invokeNativeContract error :", err)
		return false
	}
	//log4.Info("authorizeForPeer txHash is :", txHash.ToHexString())
	return true
}

func unAuthorizeForPeer(ontSdk *sdk.OntologySdk, user *sdk.Account, peerPubkeyList []string, posList []uint32) bool {
	params := &governance.AuthorizeForPeerParam{
		Address:        user.Address,
		PeerPubkeyList: peerPubkeyList,
		PosList:        posList,
	}
	contractAddress := utils.GovernanceContractAddress
	method := "unAuthorizeForPeer"
	txHash, err := ontSdk.Native.InvokeNativeContract(config.DefConfig.GasPrice, config.DefConfig.GasLimit, user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		log4.Error("invokeNativeContract error :", err)
		return false
	}
	log4.Info("unAuthorizeForPeer txHash is :", txHash.ToHexString())
	return true
}

func withdraw(ontSdk *sdk.OntologySdk, user *sdk.Account, peerPubkeyList []string, withdrawList []uint32) bool {
	params := &governance.WithdrawParam{
		Address:        user.Address,
		PeerPubkeyList: peerPubkeyList,
		WithdrawList:   withdrawList,
	}
	contractAddress := utils.GovernanceContractAddress
	method := "withdraw"
	txHash, err := ontSdk.Native.InvokeNativeContract(config.DefConfig.GasPrice, config.DefConfig.GasLimit, user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		log4.Error("invokeNativeContract error :", err)
		return false
	}
	log4.Info("withdraw txHash is :", txHash.ToHexString())
	return true
}

func withdrawOng(ontSdk *sdk.OntologySdk, user *sdk.Account) bool {
	params := &governance.WithdrawOngParam{
		Address: user.Address,
	}
	contractAddress := utils.GovernanceContractAddress
	method := "withdrawOng"
	txHash, err := ontSdk.Native.InvokeNativeContract(config.DefConfig.GasPrice, config.DefConfig.GasLimit, user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		log4.Error("invokeNativeContract error :", err)
		return false
	}
	log4.Info("withdrawOng txHash is :", txHash.ToHexString())
	return true
}

type commitDposParam struct {
}

func commitDpos(ontSdk *sdk.OntologySdk, user *sdk.Account) bool {
	params := &commitDposParam{}
	contractAddress := utils.GovernanceContractAddress
	method := "commitDpos"
	txHash, err := ontSdk.Native.InvokeNativeContract(config.DefConfig.GasPrice, config.DefConfig.GasLimit, user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		log4.Error("invokeNativeContract error :", err)
		return false
	}
	log4.Info("commitDpos txHash is :", txHash.ToHexString())
	return true
}

func commitDposMultiSign(ontSdk *sdk.OntologySdk, pubKeys []keypair.PublicKey, users []*sdk.Account) bool {
	contractAddress := utils.GovernanceContractAddress
	method := "commitDpos"
	txHash, err := common.InvokeNativeContractWithMultiSign(ontSdk, config.DefConfig.GasPrice, config.DefConfig.GasLimit, pubKeys, users, OntIDVersion,
		contractAddress, method, []interface{}{})
	if err != nil {
		log4.Error("invokeNativeContract error :", err)
		return false
	}
	log4.Info("commitDposMultiSign txHash is :", txHash.ToHexString())
	return true
}

func quitNode(ontSdk *sdk.OntologySdk, user *sdk.Account, peerPubkey string) bool {
	params := &governance.QuitNodeParam{
		PeerPubkey: peerPubkey,
		Address:    user.Address,
	}
	contractAddress := utils.GovernanceContractAddress
	method := "quitNode"
	txHash, err := ontSdk.Native.InvokeNativeContract(config.DefConfig.GasPrice, config.DefConfig.GasLimit, user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		log4.Error("invokeNativeContract error :", err)
		return false
	}
	log4.Info("quitNode txHash is :", txHash.ToHexString())
	return true
}

func blackNode(ontSdk *sdk.OntologySdk, user *sdk.Account, peerPubkeyList []string) bool {
	params := &governance.BlackNodeParam{
		PeerPubkeyList: peerPubkeyList,
	}
	contractAddress := utils.GovernanceContractAddress
	method := "blackNode"
	txHash, err := ontSdk.Native.InvokeNativeContract(config.DefConfig.GasPrice, config.DefConfig.GasLimit, user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		log4.Error("invokeNativeContract error :", err)
		return false
	}
	log4.Info("blackNode txHash is :", txHash.ToHexString())
	return true
}

func blackNodeMultiSign(ontSdk *sdk.OntologySdk, pubKeys []keypair.PublicKey, users []*sdk.Account, peerPubkeyList []string) bool {
	params := &governance.BlackNodeParam{
		PeerPubkeyList: peerPubkeyList,
	}
	contractAddress := utils.GovernanceContractAddress
	method := "blackNode"
	txHash, err := common.InvokeNativeContractWithMultiSign(ontSdk, config.DefConfig.GasPrice, config.DefConfig.GasLimit, pubKeys, users, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		log4.Error("invokeNativeContract error :", err)
		return false
	}
	log4.Info("blackNodeMultiSign txHash is :", txHash.ToHexString())
	return true
}

func whiteNode(ontSdk *sdk.OntologySdk, user *sdk.Account, peerPubkey string) bool {
	params := &governance.WhiteNodeParam{
		PeerPubkey: peerPubkey,
	}
	contractAddress := utils.GovernanceContractAddress
	method := "whiteNode"
	txHash, err := ontSdk.Native.InvokeNativeContract(config.DefConfig.GasPrice, config.DefConfig.GasLimit, user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		log4.Error("invokeNativeContract error :", err)
		return false
	}
	log4.Info("whiteNode txHash is :", txHash.ToHexString())
	return true
}

func whiteNodeMultiSign(ontSdk *sdk.OntologySdk, pubKeys []keypair.PublicKey, users []*sdk.Account, peerPubkey string) bool {
	params := &governance.WhiteNodeParam{
		PeerPubkey: peerPubkey,
	}
	contractAddress := utils.GovernanceContractAddress
	method := "whiteNode"
	txHash, err := common.InvokeNativeContractWithMultiSign(ontSdk, config.DefConfig.GasPrice, config.DefConfig.GasLimit, pubKeys, users, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		log4.Error("invokeNativeContract error :", err)
		return false
	}
	log4.Info("whiteNodeMultiSign txHash is :", txHash.ToHexString())
	return true
}

func updateConfig(ontSdk *sdk.OntologySdk, user *sdk.Account, conf *governance.Configuration) bool {
	contractAddress := utils.GovernanceContractAddress
	method := "updateConfig"
	txHash, err := ontSdk.Native.InvokeNativeContract(config.DefConfig.GasPrice, config.DefConfig.GasLimit, user, OntIDVersion,
		contractAddress, method, []interface{}{conf})
	if err != nil {
		log4.Error("invokeNativeContract error :", err)
		return false
	}
	log4.Info("updateConfig txHash is :", txHash.ToHexString())
	return true
}

func updateConfigMultiSign(ontSdk *sdk.OntologySdk, pubKeys []keypair.PublicKey, users []*sdk.Account, conf *governance.Configuration) bool {
	contractAddress := utils.GovernanceContractAddress
	method := "updateConfig"
	txHash, err := common.InvokeNativeContractWithMultiSign(ontSdk, config.DefConfig.GasPrice, config.DefConfig.GasLimit, pubKeys, users, OntIDVersion,
		contractAddress, method, []interface{}{conf})
	if err != nil {
		log4.Error("invokeNativeContract error :", err)
		return false
	}
	log4.Info("updateConfigMultiSign txHash is :", txHash.ToHexString())
	return true
}

func updateGlobalParam(ontSdk *sdk.OntologySdk, user *sdk.Account, globalParam *governance.GlobalParam) bool {
	contractAddress := utils.GovernanceContractAddress
	method := "updateGlobalParam"
	txHash, err := ontSdk.Native.InvokeNativeContract(config.DefConfig.GasPrice, config.DefConfig.GasLimit, user, OntIDVersion,
		contractAddress, method, []interface{}{globalParam})
	if err != nil {
		log4.Error("invokeNativeContract error :", err)
		return false
	}
	log4.Info("updateGlobalParam txHash is :", txHash.ToHexString())
	return true
}

func updateGlobalParamMultiSign(ontSdk *sdk.OntologySdk, pubKeys []keypair.PublicKey, users []*sdk.Account, globalParam *governance.GlobalParam) bool {
	contractAddress := utils.GovernanceContractAddress
	method := "updateGlobalParam"
	txHash, err := common.InvokeNativeContractWithMultiSign(ontSdk, config.DefConfig.GasPrice, config.DefConfig.GasLimit, pubKeys, users, OntIDVersion,
		contractAddress, method, []interface{}{globalParam})
	if err != nil {
		log4.Error("invokeNativeContract error :", err)
		return false
	}
	log4.Info("updateGlobalParamMultiSign txHash is :", txHash.ToHexString())
	return true
}

func updateGlobalParam2(ontSdk *sdk.OntologySdk, user *sdk.Account, globalParam2 *governance.GlobalParam2) bool {
	contractAddress := utils.GovernanceContractAddress
	method := "updateGlobalParam2"
	txHash, err := ontSdk.Native.InvokeNativeContract(config.DefConfig.GasPrice, config.DefConfig.GasLimit, user, OntIDVersion,
		contractAddress, method, []interface{}{globalParam2})
	if err != nil {
		log4.Error("invokeNativeContract error :", err)
		return false
	}
	log4.Info("updateGlobalParam2 txHash is :", txHash.ToHexString())
	return true
}

func updateGlobalParam2MultiSign(ontSdk *sdk.OntologySdk, pubKeys []keypair.PublicKey, users []*sdk.Account, globalParam2 *governance.GlobalParam2) bool {
	contractAddress := utils.GovernanceContractAddress
	method := "updateGlobalParam2"
	txHash, err := common.InvokeNativeContractWithMultiSign(ontSdk, config.DefConfig.GasPrice, config.DefConfig.GasLimit, pubKeys, users, OntIDVersion,
		contractAddress, method, []interface{}{globalParam2})
	if err != nil {
		log4.Error("invokeNativeContract error :", err)
		return false
	}
	log4.Info("updateGlobalParam2MultiSign txHash is :", txHash.ToHexString())
	return true
}

func updateSplitCurve(ontSdk *sdk.OntologySdk, user *sdk.Account, splitCurve *governance.SplitCurve) bool {
	contractAddress := utils.GovernanceContractAddress
	method := "updateSplitCurve"
	txHash, err := ontSdk.Native.InvokeNativeContract(config.DefConfig.GasPrice, config.DefConfig.GasLimit, user, OntIDVersion,
		contractAddress, method, []interface{}{splitCurve})
	if err != nil {
		log4.Error("invokeNativeContract error :", err)
		return false
	}
	log4.Info("updateSplitCurve txHash is :", txHash.ToHexString())
	return true
}

func updateSplitCurveMultiSign(ontSdk *sdk.OntologySdk, pubKeys []keypair.PublicKey, users []*sdk.Account, splitCurve *governance.SplitCurve) bool {
	contractAddress := utils.GovernanceContractAddress
	method := "updateSplitCurve"
	txHash, err := common.InvokeNativeContractWithMultiSign(ontSdk, config.DefConfig.GasPrice, config.DefConfig.GasLimit, pubKeys, users, OntIDVersion,
		contractAddress, method, []interface{}{splitCurve})
	if err != nil {
		log4.Error("invokeNativeContract error :", err)
		return false
	}
	log4.Info("updateSplitCurveMultiSign txHash is :", txHash.ToHexString())
	return true
}

func setPromisePos(ontSdk *sdk.OntologySdk, user *sdk.Account, promisePos *governance.PromisePos) bool {
	contractAddress := utils.GovernanceContractAddress
	method := "setPromisePos"
	txHash, err := ontSdk.Native.InvokeNativeContract(config.DefConfig.GasPrice, config.DefConfig.GasLimit, user, OntIDVersion,
		contractAddress, method, []interface{}{promisePos})
	if err != nil {
		log4.Error("invokeNativeContract error :", err)
		return false
	}
	log4.Info("setPromisePos txHash is :", txHash.ToHexString())
	return true
}

func setPromisePosMultiSign(ontSdk *sdk.OntologySdk, pubKeys []keypair.PublicKey, users []*sdk.Account, promisePos *governance.PromisePos) bool {
	contractAddress := utils.GovernanceContractAddress
	method := "setPromisePos"
	txHash, err := common.InvokeNativeContractWithMultiSign(ontSdk, config.DefConfig.GasPrice, config.DefConfig.GasLimit, pubKeys, users, OntIDVersion,
		contractAddress, method, []interface{}{promisePos})
	if err != nil {
		log4.Error("invokeNativeContract error :", err)
		return false
	}
	log4.Info("setPromisePosMultiSign txHash is :", txHash.ToHexString())
	return true
}

func transferPenalty(ontSdk *sdk.OntologySdk, user *sdk.Account, peerPubkey string, address ontcommon.Address) bool {
	params := &governance.TransferPenaltyParam{
		PeerPubkey: peerPubkey,
		Address:    address,
	}
	contractAddress := utils.GovernanceContractAddress
	method := "transferPenalty"
	txHash, err := ontSdk.Native.InvokeNativeContract(config.DefConfig.GasPrice, config.DefConfig.GasLimit, user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		log4.Error("invokeNativeContract error :", err)
		return false
	}
	log4.Info("transferPenalty txHash is :", txHash.ToHexString())
	return true
}

func withdrawFee(ontSdk *sdk.OntologySdk, user *sdk.Account) bool {
	params := &governance.WithdrawFeeParam{
		Address: user.Address,
	}
	contractAddress := utils.GovernanceContractAddress
	method := "withdrawFee"
	txHash, err := ontSdk.Native.InvokeNativeContract(config.DefConfig.GasPrice, config.DefConfig.GasLimit, user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		log4.Error("invokeNativeContract error :", err)
		return false
	}
	log4.Info("withdrawFee txHash is :", txHash.ToHexString())
	return true
}

func transferPenaltyMultiSign(ontSdk *sdk.OntologySdk, pubKeys []keypair.PublicKey, users []*sdk.Account, peerPubkey string, address ontcommon.Address) bool {
	params := &governance.TransferPenaltyParam{
		PeerPubkey: peerPubkey,
		Address:    address,
	}
	contractAddress := utils.GovernanceContractAddress
	method := "transferPenalty"
	txHash, err := common.InvokeNativeContractWithMultiSign(ontSdk, config.DefConfig.GasPrice, config.DefConfig.GasLimit, pubKeys, users, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		log4.Error("invokeNativeContract error :", err)
		return false
	}
	log4.Info("transferPenaltyMultiSign txHash is :", txHash.ToHexString())
	return true
}

func multiTransfer(ontSdk *sdk.OntologySdk, contract ontcommon.Address, from []*sdk.Account, to []string, amount []uint64) bool {
	var sts []ont.State
	if len(from) != len(to) || len(from) != len(amount) {
		log4.Error("input length error")
		return false
	}
	for i := 0; i < len(from); i++ {
		address, err := ontcommon.AddressFromBase58(to[i])
		if err != nil {
			log4.Error("common.AddressFromBase58 failed %v", err)
			return false
		}
		sts = append(sts, ont.State{
			From:  from[i].Address,
			To:    address,
			Value: amount[i],
		})
	}
	transfers := ont.Transfers{
		States: sts,
	}
	contractAddress := contract
	method := "transfer"
	tx, err := ontSdk.Native.NewNativeInvokeTransaction(config.DefConfig.GasPrice, config.DefConfig.GasLimit, OntIDVersion, contractAddress, method, []interface{}{transfers})
	if err != nil {
		return false
	}
	for _, singer := range from {
		err = ontSdk.SignToTransaction(tx, singer)
		if err != nil {
			return false
		}
	}
	txHash, err := ontSdk.SendTransaction(tx)
	if err != nil {
		log4.Error("invokeNativeContract error :", err)
		return false
		return false
	}
	log4.Info("multiTransfer txHash is :", txHash.ToHexString())
	return true
}

func transferOntMultiSign(ontSdk *sdk.OntologySdk, pubKeys []keypair.PublicKey, users []*sdk.Account, address ontcommon.Address, amount uint64) bool {
	var sts []ont.State
	from, err := types.AddressFromMultiPubKeys(pubKeys, int((5*len(pubKeys)+6)/7))
	if err != nil {
		log4.Error("types.AddressFromMultiPubKeys error", err)
	}
	sts = append(sts, ont.State{
		From:  from,
		To:    address,
		Value: amount,
	})
	transfers := ont.Transfers{
		States: sts,
	}
	contractAddress := utils.OntContractAddress
	method := "transfer"
	txHash, err := common.InvokeNativeContractWithMultiSign(ontSdk, config.DefConfig.GasPrice, config.DefConfig.GasLimit, pubKeys, users, OntIDVersion,
		contractAddress, method, []interface{}{transfers})
	if err != nil {
		log4.Error("invokeNativeContract error :", err)
		return false
	}
	log4.Info("transferOntMultiSign txHash is :", txHash.ToHexString())
	return true
}

func transferOntMultiSignToMultiSign(ontSdk *sdk.OntologySdk, pubKeys []keypair.PublicKey, users []*sdk.Account, address ontcommon.Address, amount uint64) bool {
	var sts []ont.State
	from, err := types.AddressFromMultiPubKeys(pubKeys, int((5*len(pubKeys)+6)/7))
	if err != nil {
		log4.Error("types.AddressFromMultiPubKeys error", err)
	}
	sts = append(sts, ont.State{
		From:  from,
		To:    address,
		Value: amount,
	})
	transfers := ont.Transfers{
		States: sts,
	}
	contractAddress := utils.OntContractAddress
	method := "transfer"
	txHash, err := common.InvokeNativeContractWithMultiSign(ontSdk, config.DefConfig.GasPrice, config.DefConfig.GasLimit, pubKeys, users, OntIDVersion,
		contractAddress, method, []interface{}{transfers})
	if err != nil {
		log4.Error("invokeNativeContract error :", err)
		return false
	}
	log4.Info("transferOntMultiSignToMultiSign txHash is :", txHash.ToHexString())
	return true
}

func transferOngMultiSign(ontSdk *sdk.OntologySdk, pubKeys []keypair.PublicKey, users []*sdk.Account, address ontcommon.Address, amount uint64) bool {
	var sts []ont.State
	from, err := types.AddressFromMultiPubKeys(pubKeys, int((5*len(pubKeys)+6)/7))
	if err != nil {
		log4.Error("types.AddressFromMultiPubKeys error", err)
	}
	sts = append(sts, ont.State{
		From:  from,
		To:    address,
		Value: amount,
	})
	transfers := ont.Transfers{
		States: sts,
	}
	contractAddress := utils.OngContractAddress
	method := "transfer"
	txHash, err := common.InvokeNativeContractWithMultiSign(ontSdk, config.DefConfig.GasPrice, config.DefConfig.GasLimit, pubKeys, users, OntIDVersion,
		contractAddress, method, []interface{}{transfers})
	if err != nil {
		log4.Error("invokeNativeContract error :", err)
		return false
	}
	log4.Info("transferOngMultiSign txHash is :", txHash.ToHexString())
	return true
}

func transferOngMultiSignToMultiSign(ontSdk *sdk.OntologySdk, pubKeys []keypair.PublicKey, users []*sdk.Account, address ontcommon.Address, amount uint64) bool {
	var sts []ont.State
	from, err := types.AddressFromMultiPubKeys(pubKeys, int((5*len(pubKeys)+6)/7))
	if err != nil {
		log4.Error("types.AddressFromMultiPubKeys error", err)
	}
	sts = append(sts, ont.State{
		From:  from,
		To:    address,
		Value: amount,
	})
	transfers := ont.Transfers{
		States: sts,
	}
	contractAddress := utils.OngContractAddress
	method := "transfer"
	txHash, err := common.InvokeNativeContractWithMultiSign(ontSdk, config.DefConfig.GasPrice, config.DefConfig.GasLimit, pubKeys, users, OntIDVersion,
		contractAddress, method, []interface{}{transfers})
	if err != nil {
		log4.Error("invokeNativeContract error :", err)
		return false
	}
	log4.Info("transferOngMultiSignToMultiSign txHash is :", txHash.ToHexString())
	return true
}

func transferFromOngMultiSign(ontSdk *sdk.OntologySdk, pubKeys []keypair.PublicKey, users []*sdk.Account, address ontcommon.Address, amount uint64) bool {
	from, err := types.AddressFromMultiPubKeys(pubKeys, int((5*len(pubKeys)+6)/7))
	if err != nil {
		log4.Error("types.AddressFromMultiPubKeys error", err)
	}
	params := &ont.TransferFrom{
		Sender: from,
		From:   utils.OntContractAddress,
		To:     address,
		Value:  amount,
	}
	contractAddress := utils.OngContractAddress
	method := "transferFrom"
	txHash, err := common.InvokeNativeContractWithMultiSign(ontSdk, config.DefConfig.GasPrice, config.DefConfig.GasLimit, pubKeys, users, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		log4.Error("invokeNativeContract error :", err)
		return false
	}
	log4.Info("transferFromOngMultiSign txHash is :", txHash.ToHexString())
	return true
}

func transferFromOngMultiSignToMultiSign(ontSdk *sdk.OntologySdk, pubKeys []keypair.PublicKey, users []*sdk.Account, address ontcommon.Address, amount uint64) bool {
	from, err := types.AddressFromMultiPubKeys(pubKeys, int((5*len(pubKeys)+6)/7))
	if err != nil {
		log4.Error("types.AddressFromMultiPubKeys error", err)
	}
	params := &ont.TransferFrom{
		Sender: from,
		From:   utils.OntContractAddress,
		To:     address,
		Value:  amount,
	}
	contractAddress := utils.OngContractAddress
	method := "transferFrom"
	txHash, err := common.InvokeNativeContractWithMultiSign(ontSdk, config.DefConfig.GasPrice, config.DefConfig.GasLimit, pubKeys, users, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		log4.Error("invokeNativeContract error :", err)
		return false
	}
	log4.Info("transferFromOngMultiSignToMultiSign txHash is :", txHash.ToHexString())
	return true
}

func assignFuncsToRole(ontSdk *sdk.OntologySdk, user *sdk.Account, contract ontcommon.Address, role string, function string) bool {
	params := &auth.FuncsToRoleParam{
		ContractAddr: contract,
		AdminOntID:   []byte("did:ont:" + user.Address.ToBase58()),
		Role:         []byte(role),
		FuncNames:    []string{function},
		KeyNo:        1,
	}
	method := "assignFuncsToRole"
	contractAddress := utils.AuthContractAddress
	txHash, err := ontSdk.Native.InvokeNativeContract(config.DefConfig.GasPrice, config.DefConfig.GasLimit, user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		log4.Error("invokeNativeContract error :", err)
		return false
	}
	log4.Info("assignFuncsToRole txHash is :", txHash.ToHexString())
	return true
}

func assignOntIDsToRole(ontSdk *sdk.OntologySdk, user *sdk.Account, contract ontcommon.Address, role string, ontids []string) bool {
	params := &auth.OntIDsToRoleParam{
		ContractAddr: contract,
		AdminOntID:   []byte("did:ont:" + user.Address.ToBase58()),
		Role:         []byte(role),
		Persons:      [][]byte{},
		KeyNo:        1,
	}
	for _, ontid := range ontids {
		params.Persons = append(params.Persons, []byte(ontid))
	}
	contractAddress := utils.AuthContractAddress
	method := "assignOntIDsToRole"
	txHash, err := ontSdk.Native.InvokeNativeContract(config.DefConfig.GasPrice, config.DefConfig.GasLimit, user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		log4.Error("invokeNativeContract error :", err)
		return false
	}
	log4.Info("assignOntIDsToRole txHash is :", txHash.ToHexString())
	return true
}

type RegIDWithPublicKeyParam struct {
	OntID  []byte
	Pubkey []byte
}

func regIdWithPublicKey(ontSdk *sdk.OntologySdk, user *sdk.Account) bool {
	params := RegIDWithPublicKeyParam{
		OntID:  []byte("did:ont:" + user.Address.ToBase58()),
		Pubkey: keypair.SerializePublicKey(user.PublicKey),
	}
	method := "regIDWithPublicKey"
	contractAddress := utils.OntIDContractAddress
	txHash, err := ontSdk.Native.InvokeNativeContract(config.DefConfig.GasPrice, config.DefConfig.GasLimit, user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		log4.Error("invokeNativeContract error :", err)
		return false
	}
	log4.Info("RegIDWithPublicKeyParam txHash is :", txHash.ToHexString())
	return true
}

func getVbftConfig(ontSdk *sdk.OntologySdk) (*governance.Configuration, error) {
	contractAddress := utils.GovernanceContractAddress
	config := new(governance.Configuration)
	key := []byte(governance.VBFT_CONFIG)
	value, err := ontSdk.GetStorage(contractAddress.ToHexString(), key)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "getStorage error")
	}
	if err := config.Deserialize(bytes.NewBuffer(value)); err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "deserialize, deserialize config error!")
	}
	return config, nil
}

func getGlobalParam(ontSdk *sdk.OntologySdk) (*governance.GlobalParam, error) {
	contractAddress := utils.GovernanceContractAddress
	globalParam := new(governance.GlobalParam)
	key := []byte(governance.GLOBAL_PARAM)
	value, err := ontSdk.GetStorage(contractAddress.ToHexString(), key)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "getStorage error")
	}
	if err := globalParam.Deserialize(bytes.NewBuffer(value)); err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "deserialize, deserialize globalParam error!")
	}
	return globalParam, nil
}

func getGlobalParam2(ontSdk *sdk.OntologySdk) (*governance.GlobalParam2, error) {
	contractAddress := utils.GovernanceContractAddress
	globalParam2 := new(governance.GlobalParam2)
	key := []byte(governance.GLOBAL_PARAM2)
	value, err := ontSdk.GetStorage(contractAddress.ToHexString(), key)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "getStorage error")
	}
	if len(value) != 0 {
		if err := globalParam2.Deserialize(bytes.NewBuffer(value)); err != nil {
			return nil, errors.NewDetailErr(err, errors.ErrNoCode, "deserialize, deserialize globalParam2 error!")
		}
	}
	return globalParam2, nil
}

func getSplitCurve(ontSdk *sdk.OntologySdk) (*governance.SplitCurve, error) {
	contractAddress := utils.GovernanceContractAddress
	splitCurve := new(governance.SplitCurve)
	key := []byte(governance.SPLIT_CURVE)
	value, err := ontSdk.GetStorage(contractAddress.ToHexString(), key)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "getStorage error")
	}
	if err := splitCurve.Deserialize(bytes.NewBuffer(value)); err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "deserialize, deserialize splitCurve error!")
	}
	return splitCurve, nil
}

func getGovernanceView(ontSdk *sdk.OntologySdk) (*governance.GovernanceView, error) {
	contractAddress := utils.GovernanceContractAddress
	governanceView := new(governance.GovernanceView)
	key := []byte(governance.GOVERNANCE_VIEW)
	value, err := ontSdk.GetStorage(contractAddress.ToHexString(), key)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "getStorage error")
	}
	if err := governanceView.Deserialize(bytes.NewBuffer(value)); err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "deserialize, deserialize governanceView error!")
	}
	return governanceView, nil
}

func getView(ontSdk *sdk.OntologySdk) (uint32, error) {
	governanceView, err := getGovernanceView(ontSdk)
	if err != nil {
		return 0, errors.NewDetailErr(err, errors.ErrNoCode, "getGovernanceView error")
	}
	return governanceView.View, nil
}

func getPeerPoolMap(ontSdk *sdk.OntologySdk) (*governance.PeerPoolMap, error) {
	contractAddress := utils.GovernanceContractAddress
	view, err := getView(ontSdk)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "getView error")
	}
	peerPoolMap := &governance.PeerPoolMap{
		PeerPoolMap: make(map[string]*governance.PeerPoolItem),
	}
	viewBytes, err := governance.GetUint32Bytes(view)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "GetUint32Bytes, get viewBytes error!")
	}
	key := common.ConcatKey([]byte(governance.PEER_POOL), viewBytes)
	value, err := ontSdk.GetStorage(contractAddress.ToHexString(), key)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "getStorage error")
	}
	if err := peerPoolMap.Deserialize(bytes.NewBuffer(value)); err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "deserialize, deserialize peerPoolMap error!")
	}
	return peerPoolMap, nil
}

func getAuthorizeInfo(ontSdk *sdk.OntologySdk, peerPubkey string, address ontcommon.Address) (*governance.AuthorizeInfo, error) {
	contractAddress := utils.GovernanceContractAddress
	peerPubkeyPrefix, err := hex.DecodeString(peerPubkey)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "hex.DecodeString, peerPubkey format error!")
	}
	authorizeInfo := new(governance.AuthorizeInfo)
	key := common.ConcatKey([]byte(governance.AUTHORIZE_INFO_POOL), peerPubkeyPrefix, address[:])
	value, err := ontSdk.GetStorage(contractAddress.ToHexString(), key)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "getStorage error")
	}
	if err := authorizeInfo.Deserialize(bytes.NewBuffer(value)); err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "deserialize, deserialize authorizeInfo error!")
	}
	return authorizeInfo, nil
}

func inBlackList(ontSdk *sdk.OntologySdk, peerPubkey string) (bool, error) {
	contractAddress := utils.GovernanceContractAddress
	peerPubkeyPrefix, err := hex.DecodeString(peerPubkey)
	if err != nil {
		return false, errors.NewDetailErr(err, errors.ErrNoCode, "hex.DecodeString, peerPubkey format error!")
	}
	key := common.ConcatKey([]byte(governance.BLACK_LIST), peerPubkeyPrefix)
	value, err := ontSdk.GetStorage(contractAddress.ToHexString(), key)
	if err != nil {
		return false, errors.NewDetailErr(err, errors.ErrNoCode, "getStorage error")
	}
	if len(value) != 0 {
		return true, nil
	}
	return false, nil
}

func getTotalStake(ontSdk *sdk.OntologySdk, address ontcommon.Address) (*governance.TotalStake, error) {
	contractAddress := utils.GovernanceContractAddress
	totalStake := new(governance.TotalStake)
	key := common.ConcatKey([]byte(governance.TOTAL_STAKE), address[:])
	value, err := ontSdk.GetStorage(contractAddress.ToHexString(), key)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "getStorage error")
	}
	if err := totalStake.Deserialize(bytes.NewBuffer(value)); err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "deserialize, deserialize totalStake error!")
	}
	return totalStake, nil
}

func getPenaltyStake(ontSdk *sdk.OntologySdk, peerPubkey string) (*governance.PenaltyStake, error) {
	contractAddress := utils.GovernanceContractAddress
	peerPubkeyPrefix, err := hex.DecodeString(peerPubkey)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "hex.DecodeString, peerPubkey format error!")
	}
	penaltyStake := new(governance.PenaltyStake)
	key := common.ConcatKey([]byte(governance.PENALTY_STAKE), peerPubkeyPrefix)
	value, err := ontSdk.GetStorage(contractAddress.ToHexString(), key)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "getStorage error")
	}
	if err := penaltyStake.Deserialize(bytes.NewBuffer(value)); err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "deserialize, deserialize penaltyStake error!")
	}
	return penaltyStake, nil
}

func getAttributes(ontSdk *sdk.OntologySdk, peerPubkey string) (*governance.PeerAttributes, error) {
	contractAddress := utils.GovernanceContractAddress
	peerPubkeyPrefix, err := hex.DecodeString(peerPubkey)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "hex.DecodeString, peerPubkey format error!")
	}
	peerAttributes := new(governance.PeerAttributes)
	key := common.ConcatKey([]byte(governance.PEER_ATTRIBUTES), peerPubkeyPrefix)
	value, err := ontSdk.GetStorage(contractAddress.ToHexString(), key)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "getStorage error")
	}
	if len(value) != 0 {
		if err := peerAttributes.Deserialize(bytes.NewBuffer(value)); err != nil {
			return nil, errors.NewDetailErr(err, errors.ErrNoCode, "deserialize, deserialize peerAttributes error!")
		}
	}
	return peerAttributes, nil
}

func getSplitFeeAddress(ontSdk *sdk.OntologySdk, address ontcommon.Address) (*governance.SplitFeeAddress, error) {
	contractAddress := utils.GovernanceContractAddress
	splitFeeAddress := new(governance.SplitFeeAddress)
	key := common.ConcatKey([]byte(governance.SPLIT_FEE_ADDRESS), address[:])
	value, err := ontSdk.GetStorage(contractAddress.ToHexString(), key)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "getStorage error")
	}
	if err := splitFeeAddress.Deserialize(bytes.NewBuffer(value)); err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "deserialize, deserialize splitFeeAddress error!")
	}
	return splitFeeAddress, nil
}

func getSplitFee(ontSdk *sdk.OntologySdk) (uint64, error) {
	contractAddress := utils.GovernanceContractAddress
	key := common.ConcatKey([]byte(governance.SPLIT_FEE))
	value, err := ontSdk.GetStorage(contractAddress.ToHexString(), key)
	if err != nil {
		return 0, errors.NewDetailErr(err, errors.ErrNoCode, "getStorage error")
	}
	splitFee, err := serialization.ReadUint64(bytes.NewBuffer(value))
	if err != nil {
		return 0, errors.NewDetailErr(err, errors.ErrNoCode, "serialization.ReadUint64, deserialize splitFee error!")
	}
	return splitFee, nil
}

func getPromisePos(ontSdk *sdk.OntologySdk, peerPubkey string) (*governance.PromisePos, error) {
	contractAddress := utils.GovernanceContractAddress
	peerPubkeyPrefix, err := hex.DecodeString(peerPubkey)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "hex.DecodeString, peerPubkey format error!")
	}
	key := common.ConcatKey([]byte(governance.PROMISE_POS), peerPubkeyPrefix)
	value, err := ontSdk.GetStorage(contractAddress.ToHexString(), key)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "getStorage error")
	}
	promisePos := new(governance.PromisePos)
	if err := promisePos.Deserialize(bytes.NewBuffer(value)); err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "deserialize, deserialize promisePos error!")
	}
	return promisePos, nil
}
