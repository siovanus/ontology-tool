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
	"github.com/ontio/ontology-tool/ontologytool"
)

func RegisterGovernance() {
	ontologytool.OntTool.RegMethod("RegIdWithPublicKey", RegIdWithPublicKey)
	ontologytool.OntTool.RegMethod("AssignFuncsToRole", AssignFuncsToRole)
	ontologytool.OntTool.RegMethod("AssignFuncsToRoleAny", AssignFuncsToRoleAny)
	ontologytool.OntTool.RegMethod("AssignOntIDsToRole", AssignOntIDsToRole)
	ontologytool.OntTool.RegMethod("AssignOntIDsToRoleAny", AssignOntIDsToRoleAny)
	ontologytool.OntTool.RegMethod("RegisterCandidate", RegisterCandidate)
	ontologytool.OntTool.RegMethod("RegisterCandidate2Sign", RegisterCandidate2Sign)
	ontologytool.OntTool.RegMethod("UnRegisterCandidate", UnRegisterCandidate)
	ontologytool.OntTool.RegMethod("ApproveCandidate", ApproveCandidate)
	ontologytool.OntTool.RegMethod("RejectCandidate", RejectCandidate)
	ontologytool.OntTool.RegMethod("ChangeMaxAuthorization", ChangeMaxAuthorization)
	ontologytool.OntTool.RegMethod("SetPeerCost", SetPeerCost)
	ontologytool.OntTool.RegMethod("AddInitPos", AddInitPos)
	ontologytool.OntTool.RegMethod("ReduceInitPos", ReduceInitPos)
	ontologytool.OntTool.RegMethod("AuthorizeForPeer", AuthorizeForPeer)
	ontologytool.OntTool.RegMethod("UnAuthorizeForPeer", UnAuthorizeForPeer)
	ontologytool.OntTool.RegMethod("Withdraw", Withdraw)
	ontologytool.OntTool.RegMethod("QuitNode", QuitNode)
	ontologytool.OntTool.RegMethod("BlackNode", BlackNode)
	ontologytool.OntTool.RegMethod("WhiteNode", WhiteNode)
	ontologytool.OntTool.RegMethod("CommitDpos", CommitDpos)
	ontologytool.OntTool.RegMethod("UpdateConfig", UpdateConfig)
	ontologytool.OntTool.RegMethod("UpdateGlobalParam", UpdateGlobalParam)
	ontologytool.OntTool.RegMethod("UpdateGlobalParam2", UpdateGlobalParam2)
	ontologytool.OntTool.RegMethod("UpdateSplitCurve", UpdateSplitCurve)
	ontologytool.OntTool.RegMethod("TransferPenalty", TransferPenalty)
	ontologytool.OntTool.RegMethod("SetPromisePos", SetPromisePos)
	ontologytool.OntTool.RegMethod("GetVbftConfig", GetVbftConfig)
	ontologytool.OntTool.RegMethod("GetGlobalParam", GetGlobalParam)
	ontologytool.OntTool.RegMethod("GetGlobalParam2", GetGlobalParam2)
	ontologytool.OntTool.RegMethod("GetSplitCurve", GetSplitCurve)
	ontologytool.OntTool.RegMethod("GetGovernanceView", GetGovernanceView)
	ontologytool.OntTool.RegMethod("GetPeerPoolItem", GetPeerPoolItem)
	ontologytool.OntTool.RegMethod("GetPeerPoolMap", GetPeerPoolMap)
	ontologytool.OntTool.RegMethod("GetAuthorizeInfo", GetAuthorizeInfo)
	ontologytool.OntTool.RegMethod("GetTotalStake", GetTotalStake)
	ontologytool.OntTool.RegMethod("GetPenaltyStake", GetPenaltyStake)
	ontologytool.OntTool.RegMethod("GetAttributes", GetAttributes)
	ontologytool.OntTool.RegMethod("GetSplitFee", GetSplitFee)
	ontologytool.OntTool.RegMethod("GetSplitFeeAddress", GetSplitFeeAddress)
	ontologytool.OntTool.RegMethod("GetPromisePos", GetPromisePos)
	ontologytool.OntTool.RegMethod("InBlackList", InBlackList)
	ontologytool.OntTool.RegMethod("WithdrawOng", WithdrawOng)
	ontologytool.OntTool.RegMethod("Vrf", Vrf)
	ontologytool.OntTool.RegMethod("MultiTransferOnt", MultiTransferOnt)
	ontologytool.OntTool.RegMethod("MultiTransferOng", MultiTransferOng)
	ontologytool.OntTool.RegMethod("TransferOntMultiSign", TransferOntMultiSign)
	ontologytool.OntTool.RegMethod("TransferOngMultiSign", TransferOngMultiSign)
	ontologytool.OntTool.RegMethod("TransferFromOngMultiSign", TransferFromOngMultiSign)
	ontologytool.OntTool.RegMethod("TransferOntMultiSignAddress", TransferOntMultiSignAddress)
	ontologytool.OntTool.RegMethod("TransferOngMultiSignAddress", TransferOngMultiSignAddress)
	ontologytool.OntTool.RegMethod("TransferFromOngMultiSignAddress", TransferFromOngMultiSignAddress)
	ontologytool.OntTool.RegMethod("GetAddressMultiSign", GetAddressMultiSign)
	ontologytool.OntTool.RegMethod("TransferOntMultiSignToMultiSign", TransferOntMultiSignToMultiSign)
	ontologytool.OntTool.RegMethod("TransferOngMultiSignToMultiSign", TransferOngMultiSignToMultiSign)
	ontologytool.OntTool.RegMethod("TransferFromOngMultiSignToMultiSign", TransferFromOngMultiSignToMultiSign)
	ontologytool.OntTool.RegMethod("GetVbftInfo", GetVbftInfo)
}
