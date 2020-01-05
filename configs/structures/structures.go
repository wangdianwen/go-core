package structures

import (
	"github.com/wangdianwen/go-core.git/models"
)

type UserDetail struct {
	Info          *models.User        `json:"userinfo"`
	IsSystemAdmin bool                `json:"superAdmin"`
	CInfo         *models.Client      `json:"client"`
	Roles         []*models.Role      `json:"roles"`
	Privileges    []*models.Privilege `json:"privileges"`
	Menus         []*models.Menu      `json:"menus"`
	Products      []*models.CoreApps  `json:"products"`
}
