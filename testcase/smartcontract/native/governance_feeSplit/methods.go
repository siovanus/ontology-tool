package governance_feeSplit

import (
	"github.com/ontio/ontology-crypto/keypair"
	sdk "github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology-tool/testframework"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/core/types"
	"github.com/ontio/ontology/smartcontract/service/native/auth"
	"github.com/ontio/ontology/smartcontract/service/native/ont"
	"github.com/ontio/ontology/smartcontract/service/native/utils"
)

var OntIDVersion = byte(0)

func multiTransfer(ctx *testframework.TestFrameworkContext, contract common.Address, from []*sdk.Account, to []string, amount []uint64) bool {
	var sts []ont.State
	if len(from) != len(to) || len(from) != len(amount) {
		ctx.LogError("input length error")
		return false
	}
	for i := 0; i < len(from); i++ {
		address, err := common.AddressFromBase58(to[i])
		if err != nil {
			ctx.LogError("common.AddressFromBase58 failed %v", err)
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
	tx, err := ctx.Ont.Native.NewNativeInvokeTransaction(ctx.GetChainID(), ctx.GetGasPrice(), ctx.GetGasLimit(), OntIDVersion, contractAddress, method, []interface{}{transfers})
	if err != nil {
		return false
	}
	for _, singer := range from {
		err = ctx.Ont.SignToTransaction(tx, singer)
		if err != nil {
			return false
		}
	}
	txHash, err := ctx.Ont.SendTransaction(tx)
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
		return false
		return false
	}
	ctx.LogInfo("multiTransfer txHash is :", txHash.ToHexString())
	return true
}

func transferOngMultiSign(ctx *testframework.TestFrameworkContext, pubKeys []keypair.PublicKey, users []*sdk.Account, address common.Address, amount uint64) bool {
	var sts []ont.State
	from, err := types.AddressFromMultiPubKeys(pubKeys, int((5*len(pubKeys)+6)/7))
	if err != nil {
		ctx.LogError("types.AddressFromMultiPubKeys error", err)
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
	txHash, err := invokeNativeContractWithMultiSign(ctx, ctx.GetChainID(), ctx.GetGasPrice(), ctx.GetGasLimit(), pubKeys, users, OntIDVersion,
		contractAddress, method, []interface{}{transfers})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
		return false
	}
	ctx.LogInfo("transferOngMultiSign txHash is :", txHash.ToHexString())
	return true
}

func transferOngMultiSignToMultiSign(ctx *testframework.TestFrameworkContext, pubKeys []keypair.PublicKey, users []*sdk.Account, address common.Address, amount uint64) bool {
	var sts []ont.State
	from, err := types.AddressFromMultiPubKeys(pubKeys, int((5*len(pubKeys)+6)/7))
	if err != nil {
		ctx.LogError("types.AddressFromMultiPubKeys error", err)
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
	txHash, err := invokeNativeContractWithMultiSign(ctx, ctx.GetChainID(), ctx.GetGasPrice(), ctx.GetGasLimit(), pubKeys, users, OntIDVersion,
		contractAddress, method, []interface{}{transfers})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
		return false
	}
	ctx.LogInfo("transferOngMultiSignToMultiSign txHash is :", txHash.ToHexString())
	return true
}

func assignFuncsToRole(ctx *testframework.TestFrameworkContext, user *sdk.Account, contract common.Address, role string, function string) bool {
	params := &auth.FuncsToRoleParam{
		ContractAddr: contract,
		AdminOntID:   []byte("did:ont:" + user.Address.ToBase58()),
		Role:         []byte(role),
		FuncNames:    []string{function},
		KeyNo:        1,
	}
	method := "assignFuncsToRole"
	contractAddress := utils.AuthContractAddress
	txHash, err := ctx.Ont.Native.InvokeNativeContract(ctx.GetChainID(), ctx.GetGasPrice(), ctx.GetGasLimit(), user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
		return false
	}
	ctx.LogInfo("assignFuncsToRole txHash is :", txHash.ToHexString())
	return true
}

func assignOntIDsToRole(ctx *testframework.TestFrameworkContext, user *sdk.Account, contract common.Address, role string, ontids []string) bool {
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
	txHash, err := ctx.Ont.Native.InvokeNativeContract(ctx.GetChainID(), ctx.GetGasPrice(), ctx.GetGasLimit(), user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
		return false
	}
	ctx.LogInfo("assignOntIDsToRole txHash is :", txHash.ToHexString())
	return true
}

type RegIDWithPublicKeyParam struct {
	OntID  []byte
	Pubkey []byte
}

func regIdWithPublicKey(ctx *testframework.TestFrameworkContext, user *sdk.Account) bool {
	params := RegIDWithPublicKeyParam{
		OntID:  []byte("did:ont:" + user.Address.ToBase58()),
		Pubkey: keypair.SerializePublicKey(user.PublicKey),
	}
	method := "regIDWithPublicKey"
	contractAddress := utils.OntIDContractAddress
	txHash, err := ctx.Ont.Native.InvokeNativeContract(ctx.GetChainID(), ctx.GetGasPrice(), ctx.GetGasLimit(), user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
		return false
	}
	ctx.LogInfo("RegIDWithPublicKeyParam txHash is :", txHash.ToHexString())
	return true
}
