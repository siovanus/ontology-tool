package side_chain_governance

import (
	"encoding/hex"
	"github.com/ontio/ontology-crypto/keypair"
	sdk "github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology-tool/testframework"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/errors"
	"github.com/ontio/ontology/smartcontract/service/native/chain_manager"
	"github.com/ontio/ontology/smartcontract/service/native/ong"
	"github.com/ontio/ontology/smartcontract/service/native/utils"
)

var OntIDVersion = byte(0)

func registerSideChain(ctx *testframework.TestFrameworkContext, user *sdk.Account, ratio uint32,
	deposit uint64, ongPool uint64, genesisBlockHeader []byte) bool {
	params := &chain_manager.RegisterSideChainParam{
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

func stakeSideChain(ctx *testframework.TestFrameworkContext, user *sdk.Account, chainID uint64, amount uint64) bool {
	pubKey := keypair.SerializePublicKey(user.PublicKey)
	params := &chain_manager.StakeSideChainParam{
		ChainID: chainID,
		Pubkey:  hex.EncodeToString(pubKey),
		Amount:  amount,
	}
	method := "stakeSideChain"
	contractAddress := utils.ChainManagerContractAddress
	txHash, err := ctx.Ont.Native.InvokeNativeContract(ctx.GetChainID(), ctx.GetGasPrice(), ctx.GetGasLimit(), user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
		return false
	}
	ctx.LogInfo("stakeSideChain txHash is :", txHash.ToHexString())
	waitForBlock(ctx)
	return true
}

func unStakeSideChain(ctx *testframework.TestFrameworkContext, user *sdk.Account, chainID uint64, amount uint64) bool {
	pubKey := keypair.SerializePublicKey(user.PublicKey)
	params := &chain_manager.StakeSideChainParam{
		ChainID: chainID,
		Pubkey:  hex.EncodeToString(pubKey),
		Amount:  amount,
	}
	method := "unStakeSideChain"
	contractAddress := utils.ChainManagerContractAddress
	txHash, err := ctx.Ont.Native.InvokeNativeContract(ctx.GetChainID(), ctx.GetGasPrice(), ctx.GetGasLimit(), user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
		return false
	}
	ctx.LogInfo("unStakeSideChain txHash is :", txHash.ToHexString())
	waitForBlock(ctx)
	return true
}

func inflationMultiSign(ctx *testframework.TestFrameworkContext, pubKeys []keypair.PublicKey, users []*sdk.Account,
	chainID uint64, depositAdd uint64, ongPoolAdd uint64) bool {
	params := &chain_manager.InflationParam{
		ChainID:    chainID,
		DepositAdd: depositAdd,
		OngPoolAdd: ongPoolAdd,
	}
	contractAddress := utils.ChainManagerContractAddress
	method := "inflation"
	txHash, err := invokeNativeContractWithMultiSign(ctx, ctx.GetChainID(), ctx.GetGasPrice(), ctx.GetGasLimit(), pubKeys, users, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
		return false
	}
	ctx.LogInfo("inflationMultiSign txHash is :", txHash.ToHexString())
	return true
}

func approveInflationMultiSign(ctx *testframework.TestFrameworkContext, pubKeys []keypair.PublicKey, users []*sdk.Account,
	chainID uint64) bool {
	params := &chain_manager.ChainIDParam{
		ChainID: chainID,
	}
	contractAddress := utils.ChainManagerContractAddress
	method := "approveInflation"
	txHash, err := invokeNativeContractWithMultiSign(ctx, ctx.GetChainID(), ctx.GetGasPrice(), ctx.GetGasLimit(), pubKeys, users, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
		return false
	}
	ctx.LogInfo("approveInflationMultiSign txHash is :", txHash.ToHexString())
	return true
}

func rejectInflationMultiSign(ctx *testframework.TestFrameworkContext, pubKeys []keypair.PublicKey, users []*sdk.Account,
	chainID uint64) bool {
	params := &chain_manager.ChainIDParam{
		ChainID: chainID,
	}
	contractAddress := utils.ChainManagerContractAddress
	method := "rejectInflation"
	txHash, err := invokeNativeContractWithMultiSign(ctx, ctx.GetChainID(), ctx.GetGasPrice(), ctx.GetGasLimit(), pubKeys, users, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
		return false
	}
	ctx.LogInfo("rejectInflationMultiSign txHash is :", txHash.ToHexString())
	return true
}

func ongLock(ctx *testframework.TestFrameworkContext, user *sdk.Account, ongxFee uint64, chainID uint64, ongxAmount uint64) bool {
	params := &ong.OngLockParam{
		Fee:       ongxFee,
		ToChainID: chainID,
		Address:   user.Address,
		Amount:    ongxAmount,
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
