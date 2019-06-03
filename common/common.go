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

package common

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	log4 "github.com/alecthomas/log4go"
	"github.com/ontio/ontology-crypto/keypair"
	sdk "github.com/ontio/ontology-go-sdk"
	scommon "github.com/ontio/ontology/common"
	"github.com/ontio/ontology/common/password"
	"github.com/ontio/ontology/consensus/vbft"
	"github.com/ontio/ontology/consensus/vbft/config"
	"github.com/ontio/ontology/core/types"
)

func GetAccountByPassword(sdk *sdk.OntologySdk, path string) (*sdk.Account, bool) {
	wallet, err := sdk.OpenWallet(path)
	if err != nil {
		log4.Error("open wallet error:", err)
		return nil, false
	}
	pwd, err := password.GetPassword()
	if err != nil {
		log4.Error("getPassword error:", err)
		return nil, false
	}
	user, err := wallet.GetDefaultAccount(pwd)
	if err != nil {
		log4.Error("getDefaultAccount error:", err)
		return nil, false
	}
	return user, true
}

func InvokeNativeContractWithMultiSign(
	sdk *sdk.OntologySdk,
	gasPrice,
	gasLimit uint64,
	pubKeys []keypair.PublicKey,
	singers []*sdk.Account,
	cversion byte,
	contractAddress scommon.Address,
	method string,
	params []interface{},
) (scommon.Uint256, error) {
	tx, err := sdk.Native.NewNativeInvokeTransaction(gasPrice, gasLimit, cversion, contractAddress, method, params)
	if err != nil {
		return scommon.UINT256_EMPTY, err
	}
	for _, singer := range singers {
		err = sdk.MultiSignToTransaction(tx, uint16((5*len(pubKeys)+6)/7), pubKeys, singer)
		if err != nil {
			return scommon.UINT256_EMPTY, err
		}
	}
	return sdk.SendTransaction(tx)
}

func WaitForBlock(sdk *sdk.OntologySdk) bool {
	_, err := sdk.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		log4.Error("WaitForGenerateBlock error:", err)
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

func InitVbftBlock(block *types.Block) (*vbft.Block, error) {
	if block == nil {
		return nil, fmt.Errorf("nil block in initVbftBlock")
	}

	blkInfo := &vconfig.VbftBlockInfo{}
	if err := json.Unmarshal(block.Header.ConsensusPayload, blkInfo); err != nil {
		return nil, fmt.Errorf("unmarshal blockInfo: %s", err)
	}

	return &vbft.Block{
		Block: block,
		Info:  blkInfo,
	}, nil
}

func GetAddressByHexString(hexString string) (scommon.Address, error) {
	contractByte, err := hex.DecodeString(hexString)
	if err != nil {
		return scommon.Address{}, fmt.Errorf("hex.DecodeString failed %v", err)
	}
	contractAddress, err := scommon.AddressParseFromBytes(scommon.ToArrayReverse(contractByte))
	if err != nil {
		return scommon.Address{}, fmt.Errorf("common.AddressParseFromBytes failed %v", err)
	}
	return contractAddress, nil
}
