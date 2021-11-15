package types

import (
	"Telegram2VCF/util"
	"context"
	"errors"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/tg"
)

// noSignUp can be embedded to prevent signing up.
type noSignUp struct{}

func (c noSignUp) SignUp(context.Context) (auth.UserInfo, error) {
	return auth.UserInfo{}, errors.New("not implemented")
}

func (c noSignUp) AcceptTermsOfService(_ context.Context, tos tg.HelpTermsOfService) error {
	return &auth.SignUpRequired{TermsOfService: tos}
}

// SimpleAuth implements authentication via terminal.
type SimpleAuth struct {
	noSignUp
	PhoneNumber string
}

func (a SimpleAuth) Phone(context.Context) (string, error) {
	return a.PhoneNumber, nil
}

func (a SimpleAuth) Password(context.Context) (string, error) {
	return util.Prompt("Enter 2FA password: ", true)
}

func (a SimpleAuth) Code(context.Context, *tg.AuthSentCode) (string, error) {
	return util.Prompt("Enter code: ", false)
}
