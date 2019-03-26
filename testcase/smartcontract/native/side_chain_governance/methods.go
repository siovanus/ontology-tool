package side_chain_governance

import (
	"bytes"
	"github.com/ontio/ontology-crypto/keypair"
	sdk "github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology-tool/testframework"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/errors"
	"github.com/ontio/ontology/smartcontract/service/native/chain_manager"
	"github.com/ontio/ontology/smartcontract/service/native/governance"
	"github.com/ontio/ontology/smartcontract/service/native/ong"
	"github.com/ontio/ontology/smartcontract/service/native/utils"
)

var OntIDVersion = byte(0)

func registerSideChain(ctx *testframework.TestFrameworkContext, user *sdk.Account, ratio uint32,
	deposit uint64, ongPool uint64, genesisBlockHeader []byte) bool {
	params := &chain_manager.RegisterSideChainParam{
		Address:            user.Address,
		Ratio:              ratio,
		Deposit:            deposit,
		OngPool:            ongPool,
		GenesisBlockHeader: genesisBlockHeader,
		Caller:             []byte("did:ont:" + user.Address.ToBase58()),
		KeyNo:              1,
	}
	method := "registerSideChain"
	contractAddress := utils.ChainManagerContractAddress
	txHash, err := ctx.Ont.Native.InvokeNativeContract(ctx.GetChainID(), ctx.GetGasPrice(), ctx.GetGasLimit(), user, OntIDVersion,
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
	chainID uint64) bool {
	params := &chain_manager.ChainIDParam{
		ChainID: chainID,
	}
	contractAddress := utils.ChainManagerContractAddress
	method := "approveSideChain"
	txHash, err := invokeNativeContractWithMultiSign(ctx, ctx.GetChainID(), ctx.GetGasPrice(), ctx.GetGasLimit(), pubKeys, users, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
		return false
	}
	ctx.LogInfo("approveSideChainMultiSign txHash is :", txHash.ToHexString())
	return true
}

func rejectSideChainMultiSign(ctx *testframework.TestFrameworkContext, pubKeys []keypair.PublicKey, users []*sdk.Account,
	chainID uint64) bool {
	params := &chain_manager.ChainIDParam{
		ChainID: chainID,
	}
	contractAddress := utils.ChainManagerContractAddress
	method := "rejectSideChain"
	txHash, err := invokeNativeContractWithMultiSign(ctx, ctx.GetChainID(), ctx.GetGasPrice(), ctx.GetGasLimit(), pubKeys, users, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
		return false
	}
	ctx.LogInfo("rejectSideChainMultiSign txHash is :", txHash.ToHexString())
	return true
}

func registerNodeToSideChain(ctx *testframework.TestFrameworkContext, user *sdk.Account, chainID uint64, peerPubkey string) bool {
	params := &chain_manager.NodeToSideChainParams{
		PeerPubkey: peerPubkey,
		Address:    user.Address,
		ChainID:    chainID,
	}
	method := "registerNodeToSideChain"
	contractAddress := utils.ChainManagerContractAddress
	txHash, err := ctx.Ont.Native.InvokeNativeContract(ctx.GetChainID(), ctx.GetGasPrice(), ctx.GetGasLimit(), user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
		return false
	}
	ctx.LogInfo("registerNodeToSideChain txHash is :", txHash.ToHexString())
	waitForBlock(ctx)
	return true
}

func ongLock(ctx *testframework.TestFrameworkContext, user *sdk.Account, ongxFee uint64, chainID uint64, ongxAmount uint64) bool {
	params := &ong.OngLockParam{
		OngxFee:    ongxFee,
		ToChainID:  chainID,
		Address:    user.Address,
		OngxAmount: ongxAmount,
	}
	method := "ongLock"
	contractAddress := utils.OngContractAddress
	txHash, err := ctx.Ont.Native.InvokeNativeContract(ctx.GetChainID(), ctx.GetGasPrice(), ctx.GetGasLimit(), user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
		return false
	}
	ctx.LogInfo("ongLock txHash is :", txHash.ToHexString())
	waitForBlock(ctx)
	return true
}

func registerMainChain(ctx *testframework.TestFrameworkContext, pubKeys []keypair.PublicKey, users []*sdk.Account,
	genesisBlockHeader []byte) bool {
	params := &chain_manager.RegisterMainChainParam{
		GenesisHeader: genesisBlockHeader,
	}
	contractAddress := utils.ChainManagerContractAddress
	method := "registerMainChain"
	txHash, err := invokeNativeContractWithMultiSign(ctx, ctx.GetChainID(), ctx.GetGasPrice(), ctx.GetGasLimit(), pubKeys,
		users, OntIDVersion, contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContractWithMultiSign error :", err)
	}
	ctx.LogInfo("registerMainChain txHash is :", txHash.ToHexString())
	return true
}

func getSideChain(ctx *testframework.TestFrameworkContext, sideChainID uint64) (*chain_manager.SideChain, error) {
	sideChainIDBytes, err := utils.GetUint64Bytes(sideChainID)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "getUint32Bytes error")
	}
	contractAddress := utils.ChainManagerContractAddress
	sideChain := new(chain_manager.SideChain)
	key := ConcatKey([]byte(chain_manager.SIDE_CHAIN), sideChainIDBytes)
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

func getSideChainNodeInfo(ctx *testframework.TestFrameworkContext, sideChainID uint64) (*chain_manager.SideChainNodeInfo, error) {
	sideChainIDBytes, err := utils.GetUint64Bytes(sideChainID)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "getUint32Bytes error")
	}
	contractAddress := utils.ChainManagerContractAddress
	sideChainNodeInfo := new(chain_manager.SideChainNodeInfo)
	key := ConcatKey([]byte(chain_manager.SIDE_CHAIN_NODE_INFO), sideChainIDBytes)
	value, err := ctx.Ont.GetStorage(contractAddress.ToHexString(), key)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "getStorage error")
	}
	if len(value) != 0 {
		if err := sideChainNodeInfo.Deserialization(common.NewZeroCopySource(value)); err != nil {
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
