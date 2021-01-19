package identity

import (
	"github.com/casbin/casbin/v2"
)

type AccessControl interface {
	Enforce(identityID, domain, model, field, act string) (bool, error)
	AddPolicy(identityID, domain, model, field, act string) error
}

type accessControl struct {
	enforcer *casbin.Enforcer
}

func NewAccessControl(enforcer *casbin.Enforcer) AccessControl {
	return &accessControl{
		enforcer: enforcer,
	}
}

func (ac *accessControl) Enforce(identityID, domain, model, field, act string) (bool, error) {
	ok, err := ac.enforcer.Enforce(identityID, domain, model, field, act)
	if err != nil {
		return ok, err
	}

	return ok, nil
}

func (ac *accessControl) AddPolicy(identityID, domain, model, field, act string) error {
	_, err := ac.enforcer.AddPolicy(identityID, domain, model, field, act)
	if err != nil {
		return err
	}

	return nil
}
