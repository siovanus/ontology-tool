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
	"os/exec"
	"time"

	"github.com/ontio/ontology-crypto/keypair"
	sdkcom "github.com/ontio/ontology-go-sdk/common"
	"github.com/ontio/ontology-tool/common"
	"github.com/ontio/ontology-tool/testframework"
	"github.com/ontio/ontology/account"
	scommon "github.com/ontio/ontology/common"
	"github.com/ontio/ontology/common/password"
)

func getDefaultAccount(ctx *testframework.TestFrameworkContext) (*account.Account, bool) {
	user, err := ctx.GetDefaultAccount()
	if err != nil {
		ctx.LogError("GetDefaultAccount error:%s", err)
		return nil, false
	}
	return user, true
}

func getAccountByPassword(ctx *testframework.TestFrameworkContext, path string) (*account.Account, bool) {
	wallet, err := ctx.Ont.OpenWallet(path)
	if err != nil {
		ctx.LogError("open wallet error:%s", err)
		return nil, false
	}
	pwd, err := password.GetPassword()
	if err != nil {
		ctx.LogError("getPassword error:%s", err)
		return nil, false
	}
	user, err := wallet.GetDefaultAccount(pwd)
	if err != nil {
		ctx.LogError("getDefaultAccount error:%s", err)
		return nil, false
	}
	return user, true
}

func getAccount(ctx *testframework.TestFrameworkContext, path string) (*account.Account, bool) {
	wallet, err := ctx.Ont.OpenWallet(path)
	if err != nil {
		ctx.LogError("open wallet error:%s", err)
		return nil, false
	}
	user, err := wallet.GetDefaultAccount([]byte(common.DefConfig.Password))
	if err != nil {
		ctx.LogError("getDefaultAccount error:%s", err)
		return nil, false
	}
	return user, true
}

func getAccount1(ctx *testframework.TestFrameworkContext) (*account.Account, bool) {
	wallet, err := ctx.Ont.OpenWallet("./testcase/smartcontract/native/governance_feeSplit/wallet/wallet1.dat")
	if err != nil {
		ctx.LogError("open wallet error:%s", err)
		return nil, false
	}
	user, err := wallet.GetDefaultAccount([]byte(common.DefConfig.Password))
	if err != nil {
		ctx.LogError("getDefaultAccount error:%s", err)
		return nil, false
	}
	return user, true
}

func getAccount2(ctx *testframework.TestFrameworkContext) (*account.Account, bool) {
	wallet, err := ctx.Ont.OpenWallet("./testcase/smartcontract/native/governance_feeSplit/wallet/wallet2.dat")
	if err != nil {
		ctx.LogError("open wallet error:%s", err)
		return nil, false
	}
	user, err := wallet.GetDefaultAccount([]byte(common.DefConfig.Password))
	if err != nil {
		ctx.LogError("getDefaultAccount error:%s", err)
		return nil, false
	}
	return user, true
}

func getAccount3(ctx *testframework.TestFrameworkContext) (*account.Account, bool) {
	wallet, err := ctx.Ont.OpenWallet("./testcase/smartcontract/native/governance_feeSplit/wallet/wallet3.dat")
	if err != nil {
		ctx.LogError("open wallet error:%s", err)
		return nil, false
	}
	user, err := wallet.GetDefaultAccount([]byte(common.DefConfig.Password))
	if err != nil {
		ctx.LogError("getDefaultAccount error:%s", err)
		return nil, false
	}
	return user, true
}

func invokeNativeContractWithMultiSign(
	ctx *testframework.TestFrameworkContext,
	gasPrice,
	gasLimit uint64,
	pubKeys []keypair.PublicKey,
	singers []*account.Account,
	cversion byte,
	contractAddress scommon.Address,
	method string,
	params []interface{},
) (scommon.Uint256, error) {
	tx, err := ctx.Ont.Rpc.NewNativeInvokeTransaction(gasPrice, gasLimit, cversion, contractAddress, method, params)
	if err != nil {
		return scommon.UINT256_EMPTY, err
	}
	for _, singer := range singers {
		err = sdkcom.MultiSignToTransaction(tx, uint16((5*len(pubKeys)+6)/7), pubKeys, singer)
		if err != nil {
			return scommon.UINT256_EMPTY, err
		}
	}
	return ctx.Ont.Rpc.SendRawTransaction(tx)
}

func waitForBlock(ctx *testframework.TestFrameworkContext) bool {
	_, err := ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("WaitForGenerateBlock error:%s", err)
		return false
	}
	return true
}

func ConcatKey(args ...[]byte) []byte {
	temp := []byte{}
	for _, arg := range args {
		temp = append(temp, arg...)
	}
	return temp
}

func setupTest(ctx *testframework.TestFrameworkContext, user *account.Account) bool {
	cmd := exec.Command("/bin/sh", "./testcase/smartcontract/native/governance_feeSplit/clear.sh")
	err := cmd.Start()
	if err != nil {
		ctx.LogError("run clear.sh error:%s", err)
		return false
	}
	time.Sleep(7 * time.Second)

	user, err = ctx.GetDefaultAccount()
	if err != nil {
		ctx.LogError("Wallet.GetDefaultAccount error:%s", err)
		return false
	}
	user1, ok := getAccount1(ctx)
	if !ok {
		return false
	}
	user2, ok := getAccount2(ctx)
	if !ok {
		return false
	}

	_, err = ctx.Ont.Rpc.Transfer(ctx.GetGasPrice(), ctx.GetGasLimit(), "ONT", user, user1.Address, INIT_ONT)
	if err != nil {
		ctx.LogError("Rpc.Transfer error:%s", err)
		return false
	}
	_, err = ctx.Ont.Rpc.Transfer(ctx.GetGasPrice(), ctx.GetGasLimit(), "ONT", user, user2.Address, INIT_ONT)
	if err != nil {
		ctx.LogError("Rpc.Transfer error:%s", err)
		return false
	}
	waitForBlock(ctx)
	user1Balance, err := ctx.Ont.Rpc.GetBalance(user1.Address)
	if err != nil {
		ctx.LogError("Rpc.GetBalance error:%s", err)
		return false
	}
	if user1Balance.Ont != INIT_ONT {
		ctx.LogError("balance of user1 %v is error", user1Balance.Ont)
		return false
	}
	user2Balance, err := ctx.Ont.Rpc.GetBalance(user2.Address)
	if err != nil {
		ctx.LogError("Rpc.GetBalance error:%s", err)
		return false
	}
	if user2Balance.Ont != INIT_ONT {
		ctx.LogError("balance of user2 %v is error", user2Balance.Ont)
		return false
	}

	ok = regIdWithPublicKey(ctx, user)
	if !ok {
		ctx.LogError("regIdWithPublicKey failed!")
		return false
	}
	ok = regIdWithPublicKey(ctx, user1)
	if !ok {
		ctx.LogError("regIdWithPublicKey failed!")
		return false
	}
	waitForBlock(ctx)

	ok = assignFuncsToRole(ctx, user)
	if !ok {
		ctx.LogError("assignFuncsToRole failed!")
		return false
	}
	waitForBlock(ctx)

	ok = assignOntIDsToRole(ctx, user, []*account.Account{user, user1, user2})
	if !ok {
		ctx.LogError("assignOntIDsToRole failed!")
		return false
	}
	waitForBlock(ctx)

	peerPubkeyList := []string{
		"0253ccfd439b29eca0fe90ca7c6eaa1f98572a054aa2d1d56e72ad96c466107a85",
		"035eb654bad6c6409894b9b42289a43614874c7984bde6b03aaf6fc1d0486d9d45",
		"0281d198c0dd3737a9c39191bc2d1af7d65a44261a8a64d6ef74d63f27cfb5ed92",
		"023967bba3060bf8ade06d9bad45d02853f6c623e4d4f52d767eb56df4d364a99f",
		"038bfc50b0e3f0e5df6d451069065cbfa7ab5d382a5839cce82e0c963edb026e94",
		"03f1095289e7fddb882f1cb3e158acc1c30d9de606af21c97ba851821e8b6ea535",
		"0215865baab70607f4a2413a7a9ba95ab2c3c0202d5b7731c6824eef48e899fc90",
	}
	posList := []uint32{
		20000,
		30000,
		40000,
		50000,
		60000,
		70000,
		80000,
	}
	voteForPeer(ctx, user, peerPubkeyList, posList)
	waitForBlock(ctx)

	registerCandidate(ctx, user, PEER_PUBKEY, 10000)
	waitForBlock(ctx)
	approveCandidate(ctx, user, PEER_PUBKEY)
	waitForBlock(ctx)
	return true
}

func checkBalance(ctx *testframework.TestFrameworkContext, user *account.Account, balance uint64) bool {
	userBalance, err := ctx.Ont.Rpc.GetBalance(user.Address)
	if err != nil {
		ctx.LogError("Rpc.GetBalance error:%s", err)
		return false
	}
	if userBalance.Ont != balance {
		ctx.LogError("balance of user is %v, not %v", userBalance.Ont, balance)
		return false
	}
	return true
}
