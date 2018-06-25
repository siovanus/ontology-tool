package native

import (
	"github.com/ontio/ontology-test/testcase/smartcontract/native/governance_feeSplit"
)

func TestNative() {
	//testframework.TFramework.RegTestCase("TestOntTransfer", TestOntTransfer)
	//testframework.TFramework.RegTestCase("TestWithdrawONG", TestWithdrawONG)
	//testframework.TFramework.RegTestCase("TestGlobalParam", TestGlobalParam)
	//testframework.TFramework.RegTestCase("TestAuth", auth.TestAuthContract)
	//ontid.TestNativeOntID()
	governance_feeSplit.TestGovernanceMethods()
	//governance_feeSplit.TestGovernanceContract()
	//governance_feeSplit.TestGovernanceContractError()
}
