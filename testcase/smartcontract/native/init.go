package native

import (
	//"github.com/ontio/ontology-test/testcase/smartcontract/native/auth"
	"github.com/ontio/ontology-test/testcase/smartcontract/native/ontid"
	"github.com/ontio/ontology-test/testframework"
)

func TestNative() {
	testframework.TFramework.RegTestCase("TestGlobalParam", TestGlobalParam)
	testframework.TFramework.RegTestCase("TestOntTransfer", TestOntTransfer)
	//testframework.TFramework.RegTestCase("TestAuth", auth.TestAuthContract)
	ontid.TestNativeOntID()
}
