package providers

import (
	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/contracts/validation"
	"github.com/goravel/framework/facades"

	"github.com/TheTNB/panel/v2/app/rules"
)

type ValidationServiceProvider struct {
}

func (receiver *ValidationServiceProvider) Register(app foundation.Application) {

}

func (receiver *ValidationServiceProvider) Boot(app foundation.Application) {
	if err := facades.Validation().AddRules(receiver.rules()); err != nil {
		facades.Log().Infof("add rules error: %+v", err)
	}
}

func (receiver *ValidationServiceProvider) rules() []validation.Rule {
	return []validation.Rule{
		&rules.Exists{},
		&rules.NotExists{},
		&rules.PathExists{},
		&rules.PathNotExists{},
	}
}
