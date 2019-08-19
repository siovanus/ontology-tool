package side_chain_governance

import (
	"github.com/ontio/multi-chain/smartcontract/service/native/side_chain_manager"
	"github.com/ontio/ontology-crypto/keypair"
	sdk "github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology-tool/testframework"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/smartcontract/service/native/chain_manager"
	"github.com/ontio/ontology/smartcontract/service/native/utils"
)

var OntIDVersion = byte(0)

func syncGenesisHeader(ctx *testframework.TestFrameworkContext, pubKeys []keypair.PublicKey, users []*sdk.Account,
	genesisBlockHeader []byte) bool {
	params := &chain_manager.RegisterMainChainParam{
		GenesisHeader: genesisBlockHeader,
	}
	contractAddress := utils.HeaderSyncContractAddress
	method := "syncGenesisHeader"
	txHash, err := invokeNativeContractWithMultiSign(ctx, ctx.GetChainID(), ctx.GetGasPrice(), ctx.GetGasLimit(), pubKeys,
		users, OntIDVersion, contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContractWithMultiSign error :", err)
	}
	ctx.LogInfo("syncGenesisHeader txHash is :", txHash.ToHexString())
	return true
}

func registerSideChain(ctx *testframework.TestFrameworkContext, user *sdk.Account, chainid uint64, name string,
	blocksToWait uint64) bool {
	params := &side_chain_manager.RegisterSideChainParam{
		Address:      user.Address.ToBase58(),
		Chainid:      chainid,
		Name:         name,
		BlocksToWait: blocksToWait,
	}
	contractAddress, _ := common.AddressParseFromBytes([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x011})
	method := "registerSideChain"
	txHash, err := ctx.Ont.Native.InvokeNativeContract(ctx.GetChainID(), ctx.GetGasPrice(), ctx.GetGasLimit(), user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
		return false
	}
	ctx.LogInfo("registerSideChain txHash is :", txHash.ToHexString())
	return true
}

func approveRegisterSideChain(ctx *testframework.TestFrameworkContext, pubKeys []keypair.PublicKey, users []*sdk.Account, chainid uint64) bool {
	params := &side_chain_manager.ChainidParam{
		Chainid: chainid,
	}
	contractAddress, _ := common.AddressParseFromBytes([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x011})
	method := "approveRegisterSideChain"
	txHash, err := invokeNativeContractWithMultiSign(ctx, ctx.GetChainID(), ctx.GetGasPrice(), ctx.GetGasLimit(), pubKeys,
		users, OntIDVersion, contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContractWithMultiSign error :", err)
	}
	ctx.LogInfo("approveRegisterSideChain txHash is :", txHash.ToHexString())
	return true
}

func assetMapping(ctx *testframework.TestFrameworkContext, user *sdk.Account, assetName string, assetList []*side_chain_manager.Asset) bool {
	params := &side_chain_manager.AssetMappingParam{
		Address:   user.Address.ToBase58(),
		AssetName: assetName,
		AssetList: assetList,
	}
	contractAddress, _ := common.AddressParseFromBytes([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x011})
	method := "assetMapping"
	txHash, err := ctx.Ont.Native.InvokeNativeContract(ctx.GetChainID(), ctx.GetGasPrice(), ctx.GetGasLimit(), user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContractWithMultiSign error :", err)
	}
	ctx.LogInfo("assetMapping txHash is :", txHash.ToHexString())
	return true
}
