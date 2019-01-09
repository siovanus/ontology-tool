package side_chain_governance

import (
	"bytes"

	"github.com/ontio/ontology-crypto/keypair"
	sdk "github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology-tool/testframework"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/errors"
	"github.com/ontio/ontology/smartcontract/service/native/governance"
	"github.com/ontio/ontology/smartcontract/service/native/side_chain"
	"github.com/ontio/ontology/smartcontract/service/native/utils"
	"fmt"
	"encoding/hex"
)

var OntIDVersion = byte(0)

func registerSideChain(ctx *testframework.TestFrameworkContext, user *sdk.Account, sideChainID uint32, ratio uint32,
	deposit uint64, ongPool uint64) bool {
	params := &side_chain.RegisterSideChainParam{
		SideChainID: sideChainID,
		Address:     user.Address,
		Ratio:       ratio,
		Deposit:     deposit,
		OngPool:     ongPool,
		Caller:      []byte("did:ont:" + user.Address.ToBase58()),
		KeyNo:       1,
	}
	method := "registerSideChain"
	contractAddress := utils.SideChainGovernanceContractAddress
	txHash, err := ctx.Ont.Native.InvokeNativeContract(ctx.GetGasPrice(), ctx.GetGasLimit(), user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
	}
	ctx.LogInfo("registerSideChain txHash is :", txHash.ToHexString())
	waitForBlock(ctx)

	return true
}

func approveSideChainMultiSign(ctx *testframework.TestFrameworkContext, pubKeys []keypair.PublicKey, users []*sdk.Account,
	sideChainID uint32) bool {
	params := &side_chain.SideChainIDParam{
		SideChainID: sideChainID,
	}
	contractAddress := utils.SideChainGovernanceContractAddress
	method := "approveSideChain"
	txHash, err := invokeNativeContractWithMultiSign(ctx, ctx.GetGasPrice(), ctx.GetGasLimit(), pubKeys, users, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
	}
	ctx.LogInfo("approveSideChainMultiSign txHash is :", txHash.ToHexString())
	return true
}

func rejectSideChainMultiSign(ctx *testframework.TestFrameworkContext, pubKeys []keypair.PublicKey, users []*sdk.Account,
	sideChainID uint32) bool {
	params := &side_chain.SideChainIDParam{
		SideChainID: sideChainID,
	}
	contractAddress := utils.SideChainGovernanceContractAddress
	method := "rejectSideChain"
	txHash, err := invokeNativeContractWithMultiSign(ctx, ctx.GetGasPrice(), ctx.GetGasLimit(), pubKeys, users, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
	}
	ctx.LogInfo("rejectSideChainMultiSign txHash is :", txHash.ToHexString())
	return true
}

func registerNodeToSideChain(ctx *testframework.TestFrameworkContext, user *sdk.Account, sideChainID uint32, peerPubkey string) bool {
	params := &side_chain.NodeToSideChainParams{
		PeerPubkey:  peerPubkey,
		Address:     user.Address,
		SideChainID: sideChainID,
	}
	method := "registerNodeToSideChain"
	contractAddress := utils.SideChainGovernanceContractAddress
	txHash, err := ctx.Ont.Native.InvokeNativeContract(ctx.GetGasPrice(), ctx.GetGasLimit(), user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
	}
	ctx.LogInfo("registerNodeToSideChain txHash is :", txHash.ToHexString())
	waitForBlock(ctx)
	return true
}

func getSideChain(ctx *testframework.TestFrameworkContext, sideChainID uint32) (*side_chain.SideChain, error) {
	sideChainIDBytes, err := governance.GetUint32Bytes(sideChainID)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode,"getUint32Bytes error")
	}
	contractAddress := utils.SideChainGovernanceContractAddress
	sideChain := new(side_chain.SideChain)
	key := ConcatKey([]byte(side_chain.SIDE_CHAIN), sideChainIDBytes)
	fmt.Println(hex.EncodeToString(key))
	value, err := ctx.Ont.GetStorage(contractAddress.ToHexString(), key)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "getStorage error")
	}
	if len(value) != 0 {
		if err := sideChain.Deserialize(common.NewZeroCopySource(value)); err != nil {
			return nil, errors.NewDetailErr(err, errors.ErrNoCode, "deserialize, deserialize sideChain error!")
		}
	}
	return sideChain, nil
}

func getSideChainNodeInfo(ctx *testframework.TestFrameworkContext, sideChainID uint32) (*side_chain.SideChainNodeInfo, error) {
	sideChainIDBytes, err := governance.GetUint32Bytes(sideChainID)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode,"getUint32Bytes error")
	}
	contractAddress := utils.SideChainGovernanceContractAddress
	sideChainNodeInfo := new(side_chain.SideChainNodeInfo)
	key := ConcatKey([]byte(side_chain.SIDE_CHAIN_NODE_INFO), sideChainIDBytes)
	value, err := ctx.Ont.GetStorage(contractAddress.ToHexString(), key)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "getStorage error")
	}
	if len(value) != 0 {
		if err := sideChainNodeInfo.Deserialize(bytes.NewBuffer(value)); err != nil {
			return nil, errors.NewDetailErr(err, errors.ErrNoCode, "deserialize, deserialize sideChainNodeInfo error!")
		}
	}
	return sideChainNodeInfo, nil
}

func getSideChainPeerPoolMap(ctx *testframework.TestFrameworkContext) (*governance.PeerPoolMap, error) {
	contractAddress := utils.GovernanceContractAddress
	peerPoolMap := &governance.PeerPoolMap{
		PeerPoolMap: make(map[string]*governance.PeerPoolItem),
	}
	key := []byte(governance.PEER_POOL)
	value, err := ctx.Ont.GetStorage(contractAddress.ToHexString(), key)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "getStorage error")
	}
	if err := peerPoolMap.Deserialize(bytes.NewBuffer(value)); err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "deserialize, deserialize peerPoolMap error!")
	}
	return peerPoolMap, nil
}
