package native

import (
	"github.com/ontio/ontology-tool/testcase/smartcontract/native/governance_feeSplit"
	"github.com/ontio/ontology-tool/testcase/smartcontract/native/side_chain_governance"
)

func TestNative() {
	//governance_feeSplit.TestGovernanceContract()
	//governance_feeSplit.TestGovernanceContractError()
	governance_feeSplit.TestGovernanceMethods()
	side_chain_governance.TestGovernanceMethods()
	//governance_feeSplit.TestGovernanceBatch()
}
