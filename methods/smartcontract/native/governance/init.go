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

package governance

import (
	"github.com/ontio/ontology-tool/testframework"
)

func RegisterGovernanceMethods() {
	testframework.TFramework.RegMethod("RegIdWithPublicKey", RegIdWithPublicKey)
	testframework.TFramework.RegMethod("AssignFuncsToRole", AssignFuncsToRole)
	testframework.TFramework.RegMethod("AssignFuncsToRoleAny", AssignFuncsToRoleAny)
	testframework.TFramework.RegMethod("AssignOntIDsToRole", AssignOntIDsToRole)
	testframework.TFramework.RegMethod("AssignOntIDsToRoleAny", AssignOntIDsToRoleAny)
	testframework.TFramework.RegMethod("RegisterCandidate", RegisterCandidate)
	testframework.TFramework.RegMethod("RegisterCandidate2Sign", RegisterCandidate2Sign)
	testframework.TFramework.RegMethod("UnRegisterCandidate", UnRegisterCandidate)
	testframework.TFramework.RegMethod("ApproveCandidate", ApproveCandidate)
	testframework.TFramework.RegMethod("RejectCandidate", RejectCandidate)
	testframework.TFramework.RegMethod("ChangeMaxAuthorization", ChangeMaxAuthorization)
	testframework.TFramework.RegMethod("SetPeerCost", SetPeerCost)
	testframework.TFramework.RegMethod("AddInitPos", AddInitPos)
	testframework.TFramework.RegMethod("ReduceInitPos", ReduceInitPos)
	testframework.TFramework.RegMethod("AuthorizeForPeer", AuthorizeForPeer)
	testframework.TFramework.RegMethod("UnAuthorizeForPeer", UnAuthorizeForPeer)
	testframework.TFramework.RegMethod("Withdraw", Withdraw)
	testframework.TFramework.RegMethod("QuitNode", QuitNode)
	testframework.TFramework.RegMethod("BlackNode", BlackNode)
	testframework.TFramework.RegMethod("WhiteNode", WhiteNode)
	testframework.TFramework.RegMethod("CommitDpos", CommitDpos)
	testframework.TFramework.RegMethod("UpdateConfig", UpdateConfig)
	testframework.TFramework.RegMethod("UpdateGlobalParam", UpdateGlobalParam)
	testframework.TFramework.RegMethod("UpdateGlobalParam2", UpdateGlobalParam2)
	testframework.TFramework.RegMethod("UpdateSplitCurve", UpdateSplitCurve)
	testframework.TFramework.RegMethod("TransferPenalty", TransferPenalty)
	testframework.TFramework.RegMethod("SetPromisePos", SetPromisePos)
	testframework.TFramework.RegMethod("GetVbftConfig", GetVbftConfig)
	testframework.TFramework.RegMethod("GetGlobalParam", GetGlobalParam)
	testframework.TFramework.RegMethod("GetGlobalParam2", GetGlobalParam2)
	testframework.TFramework.RegMethod("GetSplitCurve", GetSplitCurve)
	testframework.TFramework.RegMethod("GetGovernanceView", GetGovernanceView)
	testframework.TFramework.RegMethod("GetPeerPoolItem", GetPeerPoolItem)
	testframework.TFramework.RegMethod("GetPeerPoolMap", GetPeerPoolMap)
	testframework.TFramework.RegMethod("GetAuthorizeInfo", GetAuthorizeInfo)
	testframework.TFramework.RegMethod("GetTotalStake", GetTotalStake)
	testframework.TFramework.RegMethod("GetPenaltyStake", GetPenaltyStake)
	testframework.TFramework.RegMethod("GetAttributes", GetAttributes)
	testframework.TFramework.RegMethod("GetSplitFee", GetSplitFee)
	testframework.TFramework.RegMethod("GetSplitFeeAddress", GetSplitFeeAddress)
	testframework.TFramework.RegMethod("GetPromisePos", GetPromisePos)
	testframework.TFramework.RegMethod("InBlackList", InBlackList)
	testframework.TFramework.RegMethod("WithdrawOng", WithdrawOng)
	testframework.TFramework.RegMethod("Vrf", Vrf)
	testframework.TFramework.RegMethod("MultiTransferOnt", MultiTransferOnt)
	testframework.TFramework.RegMethod("MultiTransferOng", MultiTransferOng)
	testframework.TFramework.RegMethod("TransferOntMultiSign", TransferOntMultiSign)
	testframework.TFramework.RegMethod("TransferOngMultiSign", TransferOngMultiSign)
	testframework.TFramework.RegMethod("TransferFromOngMultiSign", TransferFromOngMultiSign)
	testframework.TFramework.RegMethod("TransferOntMultiSignAddress", TransferOntMultiSignAddress)
	testframework.TFramework.RegMethod("TransferOngMultiSignAddress", TransferOngMultiSignAddress)
	testframework.TFramework.RegMethod("TransferFromOngMultiSignAddress", TransferFromOngMultiSignAddress)
	testframework.TFramework.RegMethod("GetAddressMultiSign", GetAddressMultiSign)
	testframework.TFramework.RegMethod("TransferOntMultiSignToMultiSign", TransferOntMultiSignToMultiSign)
	testframework.TFramework.RegMethod("TransferOngMultiSignToMultiSign", TransferOngMultiSignToMultiSign)
	testframework.TFramework.RegMethod("TransferFromOngMultiSignToMultiSign", TransferFromOngMultiSignToMultiSign)
	testframework.TFramework.RegMethod("GetVbftInfo", GetVbftInfo)
}
