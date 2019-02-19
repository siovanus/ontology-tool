package md5check

import (
	"github.com/ontio/ontology-tool/core"
)

func RegisterMd5Check() {
	core.OntTool.RegMethod("Md5Check", Md5Check)
}
