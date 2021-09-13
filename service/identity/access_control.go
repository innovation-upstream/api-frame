package identity

import (
	"github.com/casbin/casbin/v2"
)

type (
	AccessControl interface {
		Enforce(identityID, domain, model, field, act string) (bool, error)
		AddPolicy(identityID, domain, model, field, act string) error
		RemovePolicy(identityID, domain, model, field, act string) error
		HasPolicy(identityID, domain, model, field, act string) bool
	}

	accessControl struct {
		enforcer *casbin.Enforcer
		root     string
	}

	AccessControlOption func(ac *accessControl)

	AccessControlFactory func(
		enforcer *casbin.Enforcer,
		opts ...AccessControlOption,
	) AccessControl
)

var NewAccessControl = AccessControlFactory(
	func(
		enforcer *casbin.Enforcer,
		opts ...AccessControlOption,
	) AccessControl {
		ac := &accessControl{
			enforcer: enforcer,
		}

		for _, o := range opts {
			o(ac)
		}

		return ac
	},
)

func WithRoot(root string) AccessControlOption {
	return func(ac *accessControl) {
		ac.root = root
	}
}

func (ac *accessControl) Enforce(
	identityID,
	domain,
	model,
	field,
	act string,
) (bool, error) {
	if ac.root != "" {
		return true, nil
	}

	ok, err := ac.enforcer.Enforce(identityID, domain, model, field, act)
	if err != nil {
		return ok, err
	}

	return ok, nil
}

func (ac *accessControl) AddPolicy(
	identityID,
	domain,
	model,
	field,
	act string,
) error {
	_, err := ac.enforcer.AddPolicy(identityID, domain, model, field, act, "")
	if err != nil {
		return err
	}

	return nil
}

func (ac *accessControl) RemovePolicy(
	identityID,
	domain,
	model,
	field,
	act string,
) error {
	_, err := ac.enforcer.RemovePolicy(identityID, domain, model, field, act, "")
	if err != nil {
		return err
	}

	return nil
}

func (ac *accessControl) HasPolicy(
	identityID,
	domain,
	model,
	field,
	act string,
) bool {
	return ac.enforcer.HasPolicy(identityID, domain, model, field, act, "")
}
