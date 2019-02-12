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
	"github.com/ontio/ontology-tool/core"
)

func RegisterGovernance() {
	core.OntTool.RegMethod("RegIdWithPublicKey", RegIdWithPublicKey)
	core.OntTool.RegMethod("AssignFuncsToRole", AssignFuncsToRole)
	core.OntTool.RegMethod("AssignFuncsToRoleAny", AssignFuncsToRoleAny)
	core.OntTool.RegMethod("AssignOntIDsToRole", AssignOntIDsToRole)
	core.OntTool.RegMethod("AssignOntIDsToRoleAny", AssignOntIDsToRoleAny)
	core.OntTool.RegMethod("RegisterCandidate", RegisterCandidate)
	core.OntTool.RegMethod("RegisterCandidate2Sign", RegisterCandidate2Sign)
	core.OntTool.RegMethod("UnRegisterCandidate", UnRegisterCandidate)
	core.OntTool.RegMethod("ApproveCandidate", ApproveCandidate)
	core.OntTool.RegMethod("RejectCandidate", RejectCandidate)
	core.OntTool.RegMethod("ChangeMaxAuthorization", ChangeMaxAuthorization)
	core.OntTool.RegMethod("SetPeerCost", SetPeerCost)
	core.OntTool.RegMethod("AddInitPos", AddInitPos)
	core.OntTool.RegMethod("ReduceInitPos", ReduceInitPos)
	core.OntTool.RegMethod("AuthorizeForPeer", AuthorizeForPeer)
	core.OntTool.RegMethod("UnAuthorizeForPeer", UnAuthorizeForPeer)
	core.OntTool.RegMethod("Withdraw", Withdraw)
	core.OntTool.RegMethod("QuitNode", QuitNode)
	core.OntTool.RegMethod("BlackNode", BlackNode)
	core.OntTool.RegMethod("WhiteNode", WhiteNode)
	core.OntTool.RegMethod("CommitDpos", CommitDpos)
	core.OntTool.RegMethod("UpdateConfig", UpdateConfig)
	core.OntTool.RegMethod("UpdateGlobalParam", UpdateGlobalParam)
	core.OntTool.RegMethod("UpdateGlobalParam2", UpdateGlobalParam2)
	core.OntTool.RegMethod("UpdateSplitCurve", UpdateSplitCurve)
	core.OntTool.RegMethod("TransferPenalty", TransferPenalty)
	core.OntTool.RegMethod("SetPromisePos", SetPromisePos)
	core.OntTool.RegMethod("GetVbftConfig", GetVbftConfig)
	core.OntTool.RegMethod("GetGlobalParam", GetGlobalParam)
	core.OntTool.RegMethod("GetGlobalParam2", GetGlobalParam2)
	core.OntTool.RegMethod("GetSplitCurve", GetSplitCurve)
	core.OntTool.RegMethod("GetGovernanceView", GetGovernanceView)
	core.OntTool.RegMethod("GetPeerPoolItem", GetPeerPoolItem)
	core.OntTool.RegMethod("GetPeerPoolMap", GetPeerPoolMap)
	core.OntTool.RegMethod("GetAuthorizeInfo", GetAuthorizeInfo)
	core.OntTool.RegMethod("GetTotalStake", GetTotalStake)
	core.OntTool.RegMethod("GetPenaltyStake", GetPenaltyStake)
	core.OntTool.RegMethod("GetAttributes", GetAttributes)
	core.OntTool.RegMethod("GetSplitFee", GetSplitFee)
	core.OntTool.RegMethod("GetSplitFeeAddress", GetSplitFeeAddress)
	core.OntTool.RegMethod("GetPromisePos", GetPromisePos)
	core.OntTool.RegMethod("InBlackList", InBlackList)
	core.OntTool.RegMethod("WithdrawOng", WithdrawOng)
	core.OntTool.RegMethod("Vrf", Vrf)
	core.OntTool.RegMethod("MultiTransferOnt", MultiTransferOnt)
	core.OntTool.RegMethod("MultiTransferOng", MultiTransferOng)
	core.OntTool.RegMethod("TransferOntMultiSign", TransferOntMultiSign)
	core.OntTool.RegMethod("TransferOngMultiSign", TransferOngMultiSign)
	core.OntTool.RegMethod("TransferFromOngMultiSign", TransferFromOngMultiSign)
	core.OntTool.RegMethod("TransferOntMultiSignAddress", TransferOntMultiSignAddress)
	core.OntTool.RegMethod("TransferOngMultiSignAddress", TransferOngMultiSignAddress)
	core.OntTool.RegMethod("TransferFromOngMultiSignAddress", TransferFromOngMultiSignAddress)
	core.OntTool.RegMethod("GetAddressMultiSign", GetAddressMultiSign)
	core.OntTool.RegMethod("TransferOntMultiSignToMultiSign", TransferOntMultiSignToMultiSign)
	core.OntTool.RegMethod("TransferOngMultiSignToMultiSign", TransferOngMultiSignToMultiSign)
	core.OntTool.RegMethod("TransferFromOngMultiSignToMultiSign", TransferFromOngMultiSignToMultiSign)
	core.OntTool.RegMethod("GetVbftInfo", GetVbftInfo)
}
