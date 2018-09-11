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

package main

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"math/rand"
	"net/http"
	"os"
	"time"

	"bytes"
	"github.com/ontio/ontology-crypto/keypair"
	"github.com/ontio/ontology/account"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/common/config"
	"github.com/ontio/ontology/common/constants"
	"github.com/ontio/ontology/common/log"
	"github.com/ontio/ontology/core/genesis"
	"github.com/ontio/ontology/core/ledger"
	"github.com/ontio/ontology/core/payload"
	"github.com/ontio/ontology/core/signature"
	"github.com/ontio/ontology/core/types"
	"github.com/ontio/ontology/core/validation"
	"github.com/ontio/ontology/events"
	common2 "github.com/ontio/ontology/http/base/common"
	"github.com/ontio/ontology/smartcontract/service/native/auth"
	"github.com/ontio/ontology/smartcontract/service/native/governance"
	"github.com/ontio/ontology/smartcontract/service/native/ont"
	utils2 "github.com/ontio/ontology/smartcontract/service/native/utils"
)

const (
	GAS_PRICE    = 0
	GAS_LIMIT    = 100000
	PEER_PUBKEY1 = "02de965de4cf08b85f03ab2da0fcf935200facdf663f01c0e635610b53599cf7d2"
	PEER_PUBKEY2 = "037ed3642e1c16fa7a60589abd76c68e0c049e17d80c8629fc09255448b7add26f"
	INIT_POS1    = 300000000
	INIT_POS2    = 100000000
)

func main() {
	go func() {
		fmt.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	datadir := "testdata"
	os.RemoveAll(datadir)
	log.Trace("Node version: ", config.Version)

	acct := account.NewAccount("")
	buf := keypair.SerializePublicKey(acct.PublicKey)
	config.DefConfig.Genesis.ConsensusType = "solo"
	config.DefConfig.Genesis.SOLO.Bookkeepers = []string{hex.EncodeToString(buf)}

	configuration := &config.VBFTConfig{
		N:                    7,
		C:                    2,
		K:                    7,
		L:                    112,
		BlockMsgDelay:        10000,
		HashMsgDelay:         10000,
		PeerHandshakeTimeout: 10,
		MaxBlockChangeView:   100000,
		AdminOntID:           "did:ont:AMAx993nE6NEqZjwBssUfopxnnvTdob9ij",
		MinInitStake:         100000,
		VrfValue:             "1c9810aa9822e511d5804a9c4db9dd08497c31087b0daafa34d768a3253441fa20515e2f30f81741102af0ca3cefc4818fef16adb825fbaa8cad78647f3afb590e",
		VrfProof:             "c57741f934042cb8d8b087b44b161db56fc3ffd4ffb675d36cd09f83935be853d8729f3f5298d12d6fd28d45dde515a4b9d7f67682d182ba5118abf451ff1988",
	}
	configuration.Peers = append(configuration.Peers, &config.VBFTPeerStakeInfo{Index: 1, PeerPubkey: "0253ccfd439b29eca0fe90ca7c6eaa1f98572a054aa2d1d56e72ad96c466107a85", Address: "AQTnNYR5owE587TaaEsJKUgXZ3Txij741a", InitPos: 200000000})
	configuration.Peers = append(configuration.Peers, &config.VBFTPeerStakeInfo{Index: 2, PeerPubkey: "035eb654bad6c6409894b9b42289a43614874c7984bde6b03aaf6fc1d0486d9d45", Address: "AQTnNYR5owE587TaaEsJKUgXZ3Txij741a", InitPos: 200000000})
	configuration.Peers = append(configuration.Peers, &config.VBFTPeerStakeInfo{Index: 3, PeerPubkey: "0281d198c0dd3737a9c39191bc2d1af7d65a44261a8a64d6ef74d63f27cfb5ed92", Address: "AQTnNYR5owE587TaaEsJKUgXZ3Txij741a", InitPos: 200000000})
	configuration.Peers = append(configuration.Peers, &config.VBFTPeerStakeInfo{Index: 4, PeerPubkey: "023967bba3060bf8ade06d9bad45d02853f6c623e4d4f52d767eb56df4d364a99f", Address: "AQTnNYR5owE587TaaEsJKUgXZ3Txij741a", InitPos: 200000000})
	configuration.Peers = append(configuration.Peers, &config.VBFTPeerStakeInfo{Index: 5, PeerPubkey: "038bfc50b0e3f0e5df6d451069065cbfa7ab5d382a5839cce82e0c963edb026e94", Address: "AQTnNYR5owE587TaaEsJKUgXZ3Txij741a", InitPos: 200000000})
	configuration.Peers = append(configuration.Peers, &config.VBFTPeerStakeInfo{Index: 6, PeerPubkey: "03f1095289e7fddb882f1cb3e158acc1c30d9de606af21c97ba851821e8b6ea535", Address: "AQTnNYR5owE587TaaEsJKUgXZ3Txij741a", InitPos: 200000000})
	configuration.Peers = append(configuration.Peers, &config.VBFTPeerStakeInfo{Index: 7, PeerPubkey: "0215865baab70607f4a2413a7a9ba95ab2c3c0202d5b7731c6824eef48e899fc90", Address: "AQTnNYR5owE587TaaEsJKUgXZ3Txij741a", InitPos: 200000000})

	config.DefConfig.Genesis.VBFT = configuration

	bookkeepers := []keypair.PublicKey{acct.PublicKey}
	//Init event hub
	events.Init()

	//initLog()

	log.Info("1. Loading the Ledger")
	var err error
	ledger.DefLedger, err = ledger.NewLedger(datadir)
	if err != nil {
		log.Fatalf("NewLedger error %s", err)
		os.Exit(1)
	}
	genblock, err := genesis.BuildGenesisBlock(bookkeepers, config.DefConfig.Genesis)
	if err != nil {
		log.Error(err)
		return
	}
	err = ledger.DefLedger.Init(bookkeepers, genblock)
	if err != nil {
		log.Fatalf("DefLedger.Init error %s", err)
		os.Exit(1)
	}

	//blockBuf = bytes.NewBuffer(nil)
	BlockGen(acct)
}

func BlockGen(issuer *account.Account) {
	// gen N accounts
	// initial each account with B ONT
	N := 10000
	B := 10000
	H := 10
	accounts := GenAccounts(N)

	// init accounts setup
	tsTx := make([]*types.Transaction, N)
	for i := 0; i < len(tsTx); i++ {
		ont := uint64(B)
		mutable := NewTransferTransaction(utils2.OntContractAddress, issuer.Address, accounts[i].Address, ont, GAS_PRICE, GAS_LIMIT)
		if err := signTransaction(issuer, mutable); err != nil {
			fmt.Println("signTransaction error:", err)
			os.Exit(1)
		}
		tx, _ := mutable.IntoImmutable()
		validation.VerifyTransaction(tx)
		tsTx[i] = tx
	}

	//mutable := NewOngTransferFromOngTransaction(utils2.OntContractAddress, issuer.Address, issuer.Address, 10000000000000, GAS_PRICE, GAS_LIMIT)
	//if err := signTransaction(issuer, mutable); err != nil {
	//	fmt.Println("signTransaction error:", err)
	//	os.Exit(1)
	//}
	//tx, _ := mutable.IntoImmutable()
	//validation.VerifyTransaction(tx)
	//tsTx = append(tsTx, tx)
	persistBlock(issuer, tsTx)

	//init governance contract setup
	//register ontid
	tsTx = make([]*types.Transaction, 0)
	mutable := NewRegisterOntIDTransaction(issuer, GAS_PRICE, GAS_LIMIT)
	if err := signTransaction(issuer, mutable); err != nil {
		fmt.Println("signTransaction error:", err)
		os.Exit(1)
	}
	tx, _ := mutable.IntoImmutable()
	validation.VerifyTransaction(tx)
	tsTx = append(tsTx, tx)
	admin, err := getAdminAccount()
	if err != nil {
		fmt.Println("getAdminAccount error:", err)
		os.Exit(1)
	}
	mutable = NewRegisterOntIDTransaction(admin, GAS_PRICE, GAS_LIMIT)
	if err := signTransaction(admin, mutable); err != nil {
		fmt.Println("signTransaction error:", err)
		os.Exit(1)
	}
	tx, _ = mutable.IntoImmutable()
	validation.VerifyTransaction(tx)
	tsTx = append(tsTx, tx)
	persistBlock(issuer, tsTx)

	//assign func to role
	tsTx = make([]*types.Transaction, 0)
	mutable = NewAssignFuncToRoleTransaction(admin, GAS_PRICE, GAS_LIMIT)
	if err := signTransaction(admin, mutable); err != nil {
		fmt.Println("signTransaction error:", err)
		os.Exit(1)
	}
	tx, _ = mutable.IntoImmutable()
	validation.VerifyTransaction(tx)
	tsTx = append(tsTx, tx)
	persistBlock(issuer, tsTx)

	//assign ontid to role
	tsTx = make([]*types.Transaction, 0)
	mutable = NewAssignOntIDsToRoleTransaction(admin, issuer, GAS_PRICE, GAS_LIMIT)
	if err := signTransaction(admin, mutable); err != nil {
		fmt.Println("signTransaction error:", err)
		os.Exit(1)
	}
	tx, _ = mutable.IntoImmutable()
	validation.VerifyTransaction(tx)
	tsTx = append(tsTx, tx)
	persistBlock(issuer, tsTx)

	//register two candidate
	tsTx = make([]*types.Transaction, 0)
	mutable = NewRegisterCandidateTransaction(issuer, PEER_PUBKEY1, INIT_POS1, GAS_PRICE, GAS_LIMIT)
	if err := signTransaction(issuer, mutable); err != nil {
		fmt.Println("signTransaction error:", err)
		os.Exit(1)
	}
	tx, _ = mutable.IntoImmutable()
	validation.VerifyTransaction(tx)
	tsTx = append(tsTx, tx)
	mutable = NewRegisterCandidateTransaction(issuer, PEER_PUBKEY2, INIT_POS2, GAS_PRICE, GAS_LIMIT)
	if err := signTransaction(issuer, mutable); err != nil {
		fmt.Println("signTransaction error:", err)
		os.Exit(1)
	}
	tx, _ = mutable.IntoImmutable()
	validation.VerifyTransaction(tx)
	tsTx = append(tsTx, tx)
	persistBlock(issuer, tsTx)

	//approve candidate and change peer maxAuthorize
	tsTx = make([]*types.Transaction, 0)
	mutable = NewApproveCandidateTransaction(issuer, PEER_PUBKEY1, GAS_PRICE, GAS_LIMIT)
	if err := signTransaction(issuer, mutable); err != nil {
		fmt.Println("signTransaction error:", err)
		os.Exit(1)
	}
	tx, _ = mutable.IntoImmutable()
	validation.VerifyTransaction(tx)
	tsTx = append(tsTx, tx)
	mutable = NewApproveCandidateTransaction(issuer, PEER_PUBKEY2, GAS_PRICE, GAS_LIMIT)
	if err := signTransaction(issuer, mutable); err != nil {
		fmt.Println("signTransaction error:", err)
		os.Exit(1)
	}
	tx, _ = mutable.IntoImmutable()
	validation.VerifyTransaction(tx)
	tsTx = append(tsTx, tx)
	mutable = NewChangeMaxAuthorizationTransaction(issuer, PEER_PUBKEY1, INIT_POS1, GAS_PRICE, GAS_LIMIT)
	if err := signTransaction(issuer, mutable); err != nil {
		fmt.Println("signTransaction error:", err)
		os.Exit(1)
	}
	tx, _ = mutable.IntoImmutable()
	validation.VerifyTransaction(tx)
	tsTx = append(tsTx, tx)
	mutable = NewChangeMaxAuthorizationTransaction(issuer, PEER_PUBKEY2, INIT_POS2, GAS_PRICE, GAS_LIMIT)
	if err := signTransaction(issuer, mutable); err != nil {
		fmt.Println("signTransaction error:", err)
		os.Exit(1)
	}
	tx, _ = mutable.IntoImmutable()
	validation.VerifyTransaction(tx)
	tsTx = append(tsTx, tx)
	persistBlock(issuer, tsTx)

	//fuzzy test
	for h := 0; h < H; h++ {
		tsTx := make([]*types.Transaction, 0)
		for i := 0; i < N; i++ {
			authorizeAmount1 := rand.Int() % (B / 20)
			authorizeAmount2 := rand.Int() % (B / 20)
			unAuthorizeAmount1 := rand.Int() % (B / 20)
			unAuthorizeAmount2 := rand.Int() % (B / 20)
			addInitPosAmount1 := rand.Int() % (2500)
			addInitPosAmount2 := rand.Int() % (2500)
			reduceInitPosAmount1 := rand.Int() % (2500)
			reduceInitPosAmount2 := rand.Int() % (2500)

			//authorize fuzzy
			mutable := NewAuthorizeForPeerTransaction(accounts[i], []string{PEER_PUBKEY1},
				[]uint32{uint32(authorizeAmount1)}, GAS_PRICE, GAS_LIMIT)
			if err := signTransaction(accounts[i], mutable); err != nil {
				fmt.Println("signTransaction error:", err)
				os.Exit(1)
			}
			tx, _ := mutable.IntoImmutable()
			validation.VerifyTransaction(tx)
			tsTx = append(tsTx, tx)
			mutable = NewAuthorizeForPeerTransaction(accounts[i], []string{PEER_PUBKEY2},
				[]uint32{uint32(authorizeAmount2)}, GAS_PRICE, GAS_LIMIT)
			if err := signTransaction(accounts[i], mutable); err != nil {
				fmt.Println("signTransaction error:", err)
				os.Exit(1)
			}
			tx, _ = mutable.IntoImmutable()
			validation.VerifyTransaction(tx)
			tsTx = append(tsTx, tx)
			mutable = NewUnAuthorizeForPeerTransaction(accounts[i], []string{PEER_PUBKEY1},
				[]uint32{uint32(unAuthorizeAmount1)}, GAS_PRICE, GAS_LIMIT)
			if err := signTransaction(accounts[i], mutable); err != nil {
				fmt.Println("signTransaction error:", err)
				os.Exit(1)
			}
			tx, _ = mutable.IntoImmutable()
			validation.VerifyTransaction(tx)
			tsTx = append(tsTx, tx)
			mutable = NewUnAuthorizeForPeerTransaction(accounts[i], []string{PEER_PUBKEY2},
				[]uint32{uint32(unAuthorizeAmount2)}, GAS_PRICE, GAS_LIMIT)
			if err := signTransaction(accounts[i], mutable); err != nil {
				fmt.Println("signTransaction error:", err)
				os.Exit(1)
			}
			tx, _ = mutable.IntoImmutable()
			validation.VerifyTransaction(tx)
			tsTx = append(tsTx, tx)

			//peer fuzzy
			mutable = NewAddInitPosTransaction(issuer, PEER_PUBKEY1, uint32(addInitPosAmount1), GAS_PRICE, GAS_LIMIT)
			if err := signTransaction(issuer, mutable); err != nil {
				fmt.Println("signTransaction error:", err)
				os.Exit(1)
			}
			tx, _ = mutable.IntoImmutable()
			validation.VerifyTransaction(tx)
			tsTx = append(tsTx, tx)
			mutable = NewAddInitPosTransaction(issuer, PEER_PUBKEY2, uint32(addInitPosAmount2), GAS_PRICE, GAS_LIMIT)
			if err := signTransaction(issuer, mutable); err != nil {
				fmt.Println("signTransaction error:", err)
				os.Exit(1)
			}
			tx, _ = mutable.IntoImmutable()
			validation.VerifyTransaction(tx)
			tsTx = append(tsTx, tx)
			mutable = NewReduceInitPosTransaction(issuer, PEER_PUBKEY1, uint32(reduceInitPosAmount1), GAS_PRICE, GAS_LIMIT)
			if err := signTransaction(issuer, mutable); err != nil {
				fmt.Println("signTransaction error:", err)
				os.Exit(1)
			}
			tx, _ = mutable.IntoImmutable()
			validation.VerifyTransaction(tx)
			tsTx = append(tsTx, tx)
			mutable = NewReduceInitPosTransaction(issuer, PEER_PUBKEY2, uint32(reduceInitPosAmount2), GAS_PRICE, GAS_LIMIT)
			if err := signTransaction(issuer, mutable); err != nil {
				fmt.Println("signTransaction error:", err)
				os.Exit(1)
			}
			tx, _ = mutable.IntoImmutable()
			validation.VerifyTransaction(tx)
			tsTx = append(tsTx, tx)
		}
		persistBlock(issuer, tsTx)

		//commitDpos
		tsTx = make([]*types.Transaction, 0)
		mutable = NewCommitDposTransaction(issuer, GAS_PRICE, GAS_LIMIT)
		if err := signTransaction(issuer, mutable); err != nil {
			fmt.Println("signTransaction error:", err)
			os.Exit(1)
		}
		tx, _ = mutable.IntoImmutable()
		validation.VerifyTransaction(tx)
		tsTx = append(tsTx, tx)
		persistBlock(issuer, tsTx)
		fmt.Println("current block ", ledger.DefLedger.GetCurrentBlockHeight())
	}

	// check result
	for i := 0; i < N; i++ {
		state := getState(accounts[i].Address)
		authorizeInfo1, err := getAuthorizeInfo(PEER_PUBKEY1, accounts[i].Address)
		if err != nil {
			fmt.Println("getAuthorizeInfo error:", err)
		}
		authorizeInfo2, err := getAuthorizeInfo(PEER_PUBKEY2, accounts[i].Address)
		if err != nil {
			fmt.Println("getAuthorizeInfo error:", err)
		}
		if state["ont"]+authorizeInfo1.ConsensusPos+authorizeInfo1.CandidatePos+authorizeInfo1.NewPos+authorizeInfo1.WithdrawConsensusPos+
			authorizeInfo1.WithdrawCandidatePos+authorizeInfo1.WithdrawUnfreezePos+authorizeInfo2.ConsensusPos+
			authorizeInfo2.CandidatePos+authorizeInfo2.NewPos+authorizeInfo2.WithdrawConsensusPos+authorizeInfo2.WithdrawCandidatePos+
			authorizeInfo2.WithdrawUnfreezePos != uint64(B) {
			fmt.Println("balance of account is error")
		}
	}
	state := getState(issuer.Address)
	authorizeInfo1, err := getAuthorizeInfo(PEER_PUBKEY1, issuer.Address)
	if err != nil {
		fmt.Println("getAuthorizeInfo error:", err)
	}
	authorizeInfo2, err := getAuthorizeInfo(PEER_PUBKEY2, issuer.Address)
	if err != nil {
		fmt.Println("getAuthorizeInfo error:", err)
	}
	peerPoolItem1, err := getPeerPoolItem(PEER_PUBKEY1)
	if err != nil {
		fmt.Println("getPeerPoolItem error:", err)
	}
	peerPoolItem2, err := getPeerPoolItem(PEER_PUBKEY2)
	if err != nil {
		fmt.Println("getPeerPoolItem error:", err)
	}
	fmt.Println(state["ont"])
	if state["ont"]+authorizeInfo1.WithdrawConsensusPos+authorizeInfo1.WithdrawCandidatePos+authorizeInfo1.WithdrawUnfreezePos+
		authorizeInfo2.WithdrawConsensusPos+authorizeInfo2.WithdrawCandidatePos+authorizeInfo2.WithdrawUnfreezePos+peerPoolItem1.InitPos+
		peerPoolItem2.InitPos != 900000000 {
		fmt.Println("balance of issuer is error")
	}
}

func makeBlock(acc *account.Account, txs []*types.Transaction) (*types.Block, error) {
	nextBookkeeper, err := types.AddressFromBookkeepers([]keypair.PublicKey{acc.PublicKey})
	if err != nil {
		return nil, fmt.Errorf("GetBookkeeperAddress error:%s", err)
	}
	prevHash := ledger.DefLedger.GetCurrentBlockHash()
	height := ledger.DefLedger.GetCurrentBlockHeight()

	nonce := uint64(height)
	txHash := []common.Uint256{}
	for _, t := range txs {
		txHash = append(txHash, t.Hash())
	}

	txRoot := common.ComputeMerkleRoot(txHash)
	if err != nil {
		return nil, fmt.Errorf("ComputeRoot error:%s", err)
	}

	blockRoot := ledger.DefLedger.GetBlockRootWithNewTxRoot(txRoot)
	header := &types.Header{
		Version:          0,
		PrevBlockHash:    prevHash,
		TransactionsRoot: txRoot,
		BlockRoot:        blockRoot,
		Timestamp:        constants.GENESIS_BLOCK_TIMESTAMP + height + 1,
		Height:           height + 1,
		ConsensusData:    nonce,
		NextBookkeeper:   nextBookkeeper,
	}
	block := &types.Block{
		Header:       header,
		Transactions: txs,
	}

	blockHash := block.Hash()

	sig, err := signature.Sign(acc, blockHash[:])
	if err != nil {
		return nil, fmt.Errorf("[Signature],Sign error:%s.", err)
	}

	block.Header.Bookkeepers = []keypair.PublicKey{acc.PublicKey}
	block.Header.SigData = [][]byte{sig}
	return block, nil
}

func signTransaction(signer *account.Account, tx *types.MutableTransaction) error {
	hash := tx.Hash()
	sign, _ := signature.Sign(signer, hash[:])
	tx.Sigs = append(tx.Sigs, types.Sig{
		PubKeys: []keypair.PublicKey{signer.PublicKey},
		M:       1,
		SigData: [][]byte{sign},
	})
	return nil
}

func GenAccounts(num int) []*account.Account {
	var accounts []*account.Account
	for i := 0; i < num; i++ {
		acc := account.NewAccount("")
		accounts = append(accounts, acc)
	}
	return accounts
}

func NewOngTransferFromOngTransaction(from, to, sender common.Address, value, gasPrice, gasLimit uint64) *types.MutableTransaction {
	sts := &ont.TransferFrom{
		From:   from,
		To:     to,
		Sender: sender,
		Value:  value,
	}

	invokeCode, _ := common2.BuildNativeInvokeCode(utils2.OngContractAddress, 0, "transferFrom", []interface{}{sts})
	invokePayload := &payload.InvokeCode{
		Code: invokeCode,
	}
	tx := &types.MutableTransaction{
		GasPrice: gasPrice,
		GasLimit: gasLimit,
		TxType:   types.Invoke,
		Payer:    sender,
		Nonce:    uint32(time.Now().Unix()),
		Payload:  invokePayload,
		Sigs:     nil,
	}

	return tx
}

func NewTransferTransaction(asset common.Address, from, to common.Address, value, gasPrice, gasLimit uint64) *types.MutableTransaction {
	var sts []*ont.State
	sts = append(sts, &ont.State{
		From:  from,
		To:    to,
		Value: value,
	})
	invokeCode, _ := common2.BuildNativeInvokeCode(asset, 0, "transfer", []interface{}{sts})
	invokePayload := &payload.InvokeCode{
		Code: invokeCode,
	}
	tx := &types.MutableTransaction{
		GasPrice: gasPrice,
		GasLimit: gasLimit,
		TxType:   types.Invoke,
		Payer:    from,
		Nonce:    uint32(time.Now().Unix()),
		Payload:  invokePayload,
		Sigs:     nil,
	}

	return tx
}

type RegIDWithPublicKeyParam struct {
	OntID  []byte
	Pubkey []byte
}

func NewRegisterOntIDTransaction(user *account.Account, gasPrice, gasLimit uint64) *types.MutableTransaction {
	params := RegIDWithPublicKeyParam{
		OntID:  []byte("did:ont:" + user.Address.ToBase58()),
		Pubkey: keypair.SerializePublicKey(user.PublicKey),
	}
	method := "regIDWithPublicKey"
	contractAddress := utils2.OntIDContractAddress
	invokeCode, _ := common2.BuildNativeInvokeCode(contractAddress, 0, method, []interface{}{params})
	invokePayload := &payload.InvokeCode{
		Code: invokeCode,
	}
	tx := &types.MutableTransaction{
		GasPrice: gasPrice,
		GasLimit: gasLimit,
		TxType:   types.Invoke,
		Payer:    user.Address,
		Nonce:    uint32(time.Now().Unix()),
		Payload:  invokePayload,
		Sigs:     nil,
	}

	return tx
}

func NewAssignFuncToRoleTransaction(admin *account.Account, gasPrice, gasLimit uint64) *types.MutableTransaction {
	params := &auth.FuncsToRoleParam{
		ContractAddr: utils2.GovernanceContractAddress,
		AdminOntID:   []byte("did:ont:" + admin.Address.ToBase58()),
		Role:         []byte("TrionesCandidatePeerOwner"),
		FuncNames:    []string{"registerCandidate"},
		KeyNo:        1,
	}
	method := "assignFuncsToRole"
	contractAddress := utils2.AuthContractAddress
	invokeCode, _ := common2.BuildNativeInvokeCode(contractAddress, 0, method, []interface{}{params})
	invokePayload := &payload.InvokeCode{
		Code: invokeCode,
	}
	tx := &types.MutableTransaction{
		GasPrice: gasPrice,
		GasLimit: gasLimit,
		TxType:   types.Invoke,
		Payer:    admin.Address,
		Nonce:    uint32(time.Now().Unix()),
		Payload:  invokePayload,
		Sigs:     nil,
	}

	return tx
}

func NewAssignOntIDsToRoleTransaction(admin *account.Account, user *account.Account, gasPrice, gasLimit uint64) *types.MutableTransaction {
	params := &auth.OntIDsToRoleParam{
		ContractAddr: utils2.GovernanceContractAddress,
		AdminOntID:   []byte("did:ont:" + admin.Address.ToBase58()),
		Role:         []byte("TrionesCandidatePeerOwner"),
		Persons:      [][]byte{[]byte("did:ont:" + user.Address.ToBase58())},
		KeyNo:        1,
	}
	contractAddress := utils2.AuthContractAddress
	method := "assignOntIDsToRole"
	invokeCode, _ := common2.BuildNativeInvokeCode(contractAddress, 0, method, []interface{}{params})
	invokePayload := &payload.InvokeCode{
		Code: invokeCode,
	}
	tx := &types.MutableTransaction{
		GasPrice: gasPrice,
		GasLimit: gasLimit,
		TxType:   types.Invoke,
		Payer:    admin.Address,
		Nonce:    uint32(time.Now().Unix()),
		Payload:  invokePayload,
		Sigs:     nil,
	}

	return tx
}

func NewRegisterCandidateTransaction(user *account.Account, peerPubkey string, initPos uint32, gasPrice, gasLimit uint64) *types.MutableTransaction {
	params := &governance.RegisterCandidateParam{
		PeerPubkey: peerPubkey,
		Address:    user.Address,
		InitPos:    initPos,
		Caller:     []byte("did:ont:" + user.Address.ToBase58()),
		KeyNo:      1,
	}
	method := "registerCandidate"
	contractAddress := utils2.GovernanceContractAddress
	invokeCode, _ := common2.BuildNativeInvokeCode(contractAddress, 0, method, []interface{}{params})
	invokePayload := &payload.InvokeCode{
		Code: invokeCode,
	}
	tx := &types.MutableTransaction{
		GasPrice: gasPrice,
		GasLimit: gasLimit,
		TxType:   types.Invoke,
		Payer:    user.Address,
		Nonce:    uint32(time.Now().Unix()),
		Payload:  invokePayload,
		Sigs:     nil,
	}

	return tx
}

func NewApproveCandidateTransaction(user *account.Account, peerPubkey string, gasPrice, gasLimit uint64) *types.MutableTransaction {
	params := &governance.ApproveCandidateParam{
		PeerPubkey: peerPubkey,
	}
	contractAddress := utils2.GovernanceContractAddress
	method := "approveCandidate"
	invokeCode, _ := common2.BuildNativeInvokeCode(contractAddress, 0, method, []interface{}{params})
	invokePayload := &payload.InvokeCode{
		Code: invokeCode,
	}
	tx := &types.MutableTransaction{
		GasPrice: gasPrice,
		GasLimit: gasLimit,
		TxType:   types.Invoke,
		Payer:    user.Address,
		Nonce:    uint32(time.Now().Unix()),
		Payload:  invokePayload,
		Sigs:     nil,
	}

	return tx
}

func NewChangeMaxAuthorizationTransaction(user *account.Account, peerPubkey string, maxAuthorize uint32, gasPrice, gasLimit uint64) *types.MutableTransaction {
	params := &governance.ChangeMaxAuthorizationParam{
		Address:      user.Address,
		PeerPubkey:   peerPubkey,
		MaxAuthorize: maxAuthorize,
	}
	contractAddress := utils2.GovernanceContractAddress
	method := "changeMaxAuthorization"
	invokeCode, _ := common2.BuildNativeInvokeCode(contractAddress, 0, method, []interface{}{params})
	invokePayload := &payload.InvokeCode{
		Code: invokeCode,
	}
	tx := &types.MutableTransaction{
		GasPrice: gasPrice,
		GasLimit: gasLimit,
		TxType:   types.Invoke,
		Payer:    user.Address,
		Nonce:    uint32(time.Now().Unix()),
		Payload:  invokePayload,
		Sigs:     nil,
	}

	return tx
}

func NewAuthorizeForPeerTransaction(user *account.Account, peerPubkeyList []string, posList []uint32, gasPrice, gasLimit uint64) *types.MutableTransaction {
	params := &governance.AuthorizeForPeerParam{
		Address:        user.Address,
		PeerPubkeyList: peerPubkeyList,
		PosList:        posList,
	}
	contractAddress := utils2.GovernanceContractAddress
	method := "authorizeForPeer"
	invokeCode, _ := common2.BuildNativeInvokeCode(contractAddress, 0, method, []interface{}{params})
	invokePayload := &payload.InvokeCode{
		Code: invokeCode,
	}
	tx := &types.MutableTransaction{
		GasPrice: gasPrice,
		GasLimit: gasLimit,
		TxType:   types.Invoke,
		Payer:    user.Address,
		Nonce:    uint32(time.Now().Unix()),
		Payload:  invokePayload,
		Sigs:     nil,
	}

	return tx
}

func NewUnAuthorizeForPeerTransaction(user *account.Account, peerPubkeyList []string, posList []uint32, gasPrice, gasLimit uint64) *types.MutableTransaction {
	params := &governance.AuthorizeForPeerParam{
		Address:        user.Address,
		PeerPubkeyList: peerPubkeyList,
		PosList:        posList,
	}
	contractAddress := utils2.GovernanceContractAddress
	method := "unAuthorizeForPeer"
	invokeCode, _ := common2.BuildNativeInvokeCode(contractAddress, 0, method, []interface{}{params})
	invokePayload := &payload.InvokeCode{
		Code: invokeCode,
	}
	tx := &types.MutableTransaction{
		GasPrice: gasPrice,
		GasLimit: gasLimit,
		TxType:   types.Invoke,
		Payer:    user.Address,
		Nonce:    uint32(time.Now().Unix()),
		Payload:  invokePayload,
		Sigs:     nil,
	}

	return tx
}

func NewAddInitPosTransaction(user *account.Account, peerPubkey string, pos uint32, gasPrice, gasLimit uint64) *types.MutableTransaction {
	params := &governance.ChangeInitPosParam{
		Address:    user.Address,
		PeerPubkey: peerPubkey,
		Pos:        pos,
	}
	contractAddress := utils2.GovernanceContractAddress
	method := "addInitPos"
	invokeCode, _ := common2.BuildNativeInvokeCode(contractAddress, 0, method, []interface{}{params})
	invokePayload := &payload.InvokeCode{
		Code: invokeCode,
	}
	tx := &types.MutableTransaction{
		GasPrice: gasPrice,
		GasLimit: gasLimit,
		TxType:   types.Invoke,
		Payer:    user.Address,
		Nonce:    uint32(time.Now().Unix()),
		Payload:  invokePayload,
		Sigs:     nil,
	}

	return tx
}

func NewReduceInitPosTransaction(user *account.Account, peerPubkey string, pos uint32, gasPrice, gasLimit uint64) *types.MutableTransaction {
	params := &governance.ChangeInitPosParam{
		Address:    user.Address,
		PeerPubkey: peerPubkey,
		Pos:        pos,
	}
	contractAddress := utils2.GovernanceContractAddress
	method := "reduceInitPos"
	invokeCode, _ := common2.BuildNativeInvokeCode(contractAddress, 0, method, []interface{}{params})
	invokePayload := &payload.InvokeCode{
		Code: invokeCode,
	}
	tx := &types.MutableTransaction{
		GasPrice: gasPrice,
		GasLimit: gasLimit,
		TxType:   types.Invoke,
		Payer:    user.Address,
		Nonce:    uint32(time.Now().Unix()),
		Payload:  invokePayload,
		Sigs:     nil,
	}

	return tx
}

func NewCommitDposTransaction(user *account.Account, gasPrice, gasLimit uint64) *types.MutableTransaction {
	params := &governance.ApproveCandidateParam{}
	contractAddress := utils2.GovernanceContractAddress
	method := "commitDpos"
	invokeCode, _ := common2.BuildNativeInvokeCode(contractAddress, 0, method, []interface{}{params})
	invokePayload := &payload.InvokeCode{
		Code: invokeCode,
	}
	tx := &types.MutableTransaction{
		GasPrice: gasPrice,
		GasLimit: gasLimit,
		TxType:   types.Invoke,
		Payer:    user.Address,
		Nonce:    uint32(time.Now().Unix()),
		Payload:  invokePayload,
		Sigs:     nil,
	}

	return tx
}

func getAdminAccount() (*account.Account, error) {
	wallet, err := account.Open("wallet.dat")
	if err != nil {
		return nil, err
	}
	user, err := wallet.GetDefaultAccount([]byte("passwordtest"))
	if err != nil {
		return nil, err
	}
	return user, nil
}

func getState(addr common.Address) map[string]uint64 {
	ont := new(big.Int)
	ong := new(big.Int)
	appove := new(big.Int)

	ontBalance, _ := ledger.DefLedger.GetStorageItem(utils2.OntContractAddress, addr[:])
	if ontBalance != nil {
		ont = common.BigIntFromNeoBytes(ontBalance)
	}
	ongBalance, _ := ledger.DefLedger.GetStorageItem(utils2.OngContractAddress, addr[:])
	if ongBalance != nil {
		ong = common.BigIntFromNeoBytes(ongBalance)
	}

	appoveKey := append(utils2.OntContractAddress[:], addr[:]...)
	ongappove, _ := ledger.DefLedger.GetStorageItem(utils2.OngContractAddress, appoveKey[:])
	if ongappove != nil {
		appove = common.BigIntFromNeoBytes(ongappove)
	}

	rsp := make(map[string]uint64)
	rsp["ont"] = ont.Uint64()
	rsp["ong"] = ong.Uint64()
	rsp["ongAppove"] = appove.Uint64()

	return rsp
}

func getAuthorizeInfo(peerPubkey string, addr common.Address) (*governance.AuthorizeInfo, error) {
	authorizeInfo := new(governance.AuthorizeInfo)
	peerPubkeyPrefix, err := hex.DecodeString(peerPubkey)
	if err != nil {
		return nil, err
	}
	key := concatKey([]byte(governance.AUTHORIZE_INFO_POOL), peerPubkeyPrefix, addr[:])
	authorizeInfoBytes, _ := ledger.DefLedger.GetStorageItem(utils2.GovernanceContractAddress, key)
	if authorizeInfoBytes != nil {
		if err := authorizeInfo.Deserialize(bytes.NewBuffer(authorizeInfoBytes)); err != nil {
			return nil, err
		}
	}

	return authorizeInfo, nil
}

func getPeerPoolItem(peerPubkey string) (*governance.PeerPoolItem, error) {
	governanceView := new(governance.GovernanceView)
	key := concatKey([]byte(governance.GOVERNANCE_VIEW))
	governanceViewBytes, _ := ledger.DefLedger.GetStorageItem(utils2.GovernanceContractAddress, key)
	if governanceViewBytes != nil {
		if err := governanceView.Deserialize(bytes.NewBuffer(governanceViewBytes)); err != nil {
			return nil, err
		}
	}
	view := governanceView.View
	viewBytes, err := governance.GetUint32Bytes(view)
	if err != nil {
		return nil, err
	}
	peerPoolMap := new(governance.PeerPoolMap)
	key = concatKey([]byte(governance.PEER_POOL), viewBytes)
	peerPoolMapBytes, _ := ledger.DefLedger.GetStorageItem(utils2.GovernanceContractAddress, key)
	if peerPoolMapBytes != nil {
		if err := peerPoolMap.Deserialize(bytes.NewBuffer(peerPoolMapBytes)); err != nil {
			return nil, err
		}
	}

	return peerPoolMap.PeerPoolMap[peerPubkey], nil
}

func persistBlock(issuer *account.Account, tsTx []*types.Transaction) {
	block, _ := makeBlock(issuer, tsTx)
	err := ledger.DefLedger.AddBlock(block)
	if err != nil {
		fmt.Println("persist block error", err)
	}
}

func concatKey(args ...[]byte) []byte {
	temp := []byte{}
	for _, arg := range args {
		temp = append(temp, arg...)
	}
	return temp
}

func initLog() {
	//init log module
	logLevel := 1
	log.InitLog(logLevel, log.PATH, log.Stdout)
}
