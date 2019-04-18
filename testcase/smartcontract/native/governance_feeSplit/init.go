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
	"github.com/ontio/ontology-tool/testframework"
)

func TestGovernanceMethods() {
	testframework.TFramework.RegTestCase("RegIdWithPublicKey", RegIdWithPublicKey)
	testframework.TFramework.RegTestCase("AssignFuncsToRole", AssignFuncsToRole)
	testframework.TFramework.RegTestCase("AssignFuncsToRoleAny", AssignFuncsToRoleAny)
	testframework.TFramework.RegTestCase("AssignOntIDsToRole", AssignOntIDsToRole)
	testframework.TFramework.RegTestCase("AssignOntIDsToRoleAny", AssignOntIDsToRoleAny)
	testframework.TFramework.RegTestCase("Vrf", Vrf)
	testframework.TFramework.RegTestCase("TransferOntMultiSign", TransferOntMultiSign)
	testframework.TFramework.RegTestCase("TransferFromOngMultiSign", TransferFromOngMultiSign)
	testframework.TFramework.RegTestCase("GetAddressMultiSign", GetAddressMultiSign)
	testframework.TFramework.RegTestCase("GetVbftInfo", GetVbftInfo)
}
