package semantic

import "github.com/strict-lang/sdk/pkg/compiler/pass"

func Run(context *pass.Context) error {
	return pass.RunWithId(NameResolutionPassId, context)
}
