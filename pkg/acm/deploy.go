package acm

import (
	"github.com/karmab/tasty/pkg/operator"
)

func (a *ACM) RunDeployACM() error {
	o := operator.NewOperator()
	err := o.InstallOperator(true, false, a.Namespace, a.DefaultChannel, []string{a.Name})
	if err != nil {
		return err
	}

	return nil
}
