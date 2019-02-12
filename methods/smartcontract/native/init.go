package native

import (
	"github.com/ontio/ontology-tool/methods/smartcontract/native/governance"
)

func RegisterNative() {
	governance.RegisterGovernance()
}
