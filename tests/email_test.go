package tests

import (
	"../middlewares/email/models"
	"testing"
)

func TestSysAlert(t *testing.T) {
	service := Prepare()
	defer Finish(service)
	bErr := models.SysAlert("POS", "this is a test message!<br> sending automatically!")
	t.Log(bErr)
}
