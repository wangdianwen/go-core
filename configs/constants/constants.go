package constants

import "github.com/wangdianwen/go-core.git/utils"

const (
	AppName = utils.APPCore

	// user types
	TypeNormalUser          = int8(0)
	TypeFreeTrialUser       = int8(1)
	TypeEmptyUser           = int8(2)
	TypeFreeTrialExpireUser = int8(3)

	// user resources
	SourceAdminAdd = int8(0)
	SourceRegister = int8(1)

	// user email verified
	EmailVerifiedUser    = int8(0)
	EmailNotVerifiedUser = int8(1)

	FreeTrialMaxDays = 30

	VerifyEmailExpireSeconds       = 2 * 86400
	VerifyEmailSendIntervalSeconds = 60

	UndefinedId = -1

	StatusForbidden = int8(1)
	StatusNormal    = int8(0)

	ClientBackEndTypeMySql    = int8(0)
	ClientBackEndTypeAccredo  = int8(1)
	ClientBackEndTypeAdvanced = int8(2)
)
