package native

import (
	"github.com/ontio/ontology-tool/testcase/smartcontract/native/side_chain_governance"
)

func TestNative() {
	side_chain_governance.TestGovernanceMethods()
}
