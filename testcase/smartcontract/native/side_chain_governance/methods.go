package side_chain_governance

import (
	"bytes"

	"encoding/hex"
	"fmt"
	"github.com/ontio/ontology-crypto/keypair"
	sdk "github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology-tool/testframework"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/errors"
	"github.com/ontio/ontology/smartcontract/service/native/governance"
	"github.com/ontio/ontology/smartcontract/service/native/side_chain"
	"github.com/ontio/ontology/smartcontract/service/native/utils"
)

var OntIDVersion = byte(0)

func registerSideChain(ctx *testframework.TestFrameworkContext, user *sdk.Account, ratio uint32,
	deposit uint64, ongPool uint64, genesisBlockHeader []byte) bool {
	params := &side_chain.RegisterSideChainParam{
		Address:            user.Address,
		Ratio:              ratio,
		Deposit:            deposit,
		OngPool:            ongPool,
		GenesisBlockHeader: genesisBlockHeader,
		Caller:             []byte("did:ont:" + user.Address.ToBase58()),
		KeyNo:              1,
	}
	method := "registerSideChain"
	contractAddress := utils.SideChainGovernanceContractAddress
	txHash, err := ctx.Ont.Native.InvokeNativeContract(ctx.GetGasPrice(), ctx.GetGasLimit(), user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
		return false
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
		return false
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
		return false
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
		return false
	}
	ctx.LogInfo("registerNodeToSideChain txHash is :", txHash.ToHexString())
	waitForBlock(ctx)
	return true
}

func ongLock(ctx *testframework.TestFrameworkContext, user *sdk.Account, sideChainID uint32, ongxAmount uint64) bool {
	params := &side_chain.OngLockParam{
		SideChainID: sideChainID,
		Address:     user.Address,
		OngxAmount:  ongxAmount,
	}
	method := "ongLock"
	contractAddress := utils.SideChainGovernanceContractAddress
	txHash, err := ctx.Ont.Native.InvokeNativeContract(ctx.GetGasPrice(), ctx.GetGasLimit(), user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
		return false
	}
	ctx.LogInfo("ongLock txHash is :", txHash.ToHexString())
	waitForBlock(ctx)
	return true
}

func ongUnlock(ctx *testframework.TestFrameworkContext, user *sdk.Account, sideChainID uint32, txHash common.Uint256, rpc string) bool {
	//get block hash and mpt proof of side chain
	sideChainSdk := sdk.NewOntologySdk()
	sideChainSdk.NewRpcClient().SetAddress(rpc)
	key := hex.EncodeToString(append(utils.OngContractAddress[:], txHash.ToArray()...))
	height, err := sideChainSdk.ClientMgr.GetBlockHeightByTxHash(txHash.ToHexString())
	if err != nil {
		ctx.LogError("sideChainSdk.ClientMgr.GetBlockHeightByTxHash error :", err)
		return false
	}
	blockHash, err := sideChainSdk.GetBlockHash(height)
	if err != nil {
		ctx.LogError("sideChainSdk.GetBlockHash error :", err)
		return false
	}
	mptProof, err := sideChainSdk.ClientMgr.GetMPTProof(key, blockHash.ToHexString())
	if err != nil {
		ctx.LogError("sideChainSdk.ClientMgr.GetMPTProof error :", err)
		return false
	}
	var proof [][]byte
	for _, v := range mptProof.MPTProof {
		proof = append(proof, v)
	}
	params := &side_chain.OngUnlockParam{
		SideChainID: sideChainID,
		BlockHash:   blockHash,
		TxHash:      txHash,
		Proof:       proof,
	}
	method := "ongUnlock"
	contractAddress := utils.SideChainGovernanceContractAddress
	txHash2, err := ctx.Ont.Native.InvokeNativeContract(ctx.GetGasPrice(), ctx.GetGasLimit(), user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
		return false
	}
	ctx.LogInfo("ongUnlock txHash is :", txHash2.ToHexString())
	waitForBlock(ctx)
	return true
}

func getSideChain(ctx *testframework.TestFrameworkContext, sideChainID uint32) (*side_chain.SideChain, error) {
	sideChainIDBytes, err := governance.GetUint32Bytes(sideChainID)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "getUint32Bytes error")
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
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "getUint32Bytes error")
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
