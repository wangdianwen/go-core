package models

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	task "go-core/middlewares/task/models"
	"go-core/utils"
)

var (
	emailHeader = `
<!DOCTYPE html>
<html lang="en">
    <head>
        <meta charset="UTF-8">
        <title>Activate Your BizEX Account</title>
        <style type="text/css">
        .footer {
        margin-top: 100px;
        font-size: 10px;
        }
        .footer p {
        line-height: 5px;
        color: gray;
        }
        .footer a {
        color: gray
        }
        a:link {
        text-decoration: none;
        }
        a:visited {
        text-decoration: none;
        }
        .red {
        color: red
        }
        .blue {
        color: #00aeef
        }
		.header img {
		width: 590px;
		height: 120px;
		}
        </style>
    </head>
    <body>
        <div class="header">
            <img src="cid:bizex_logo.png" alt="logo header" />
            <p class="red"><b>%s</b></p>
        </div>
        <hr>
`
	emailFooter = `
		<br />
        <p>Thanks,</p>
        <p class="blue">BizEX RD.</p>
        <div class="footer">
            <p>Web: <a href="https://www.bizex.co.nz">https://www.bizex.co.nz/</a></p>
            <p>Tel: 09 529 4492</p>
            <p>Address: 109 Great South Rd, Epsom, Auckland</p>
            <!-- <img src="cid:brunton_footer.jpg" alt="logo footer" /> -->
			<p style="font-size:6px;line-height: 15px;">This e-mail is only intended to be read by the named recipient.  It may contain information which is confidential, proprietary or the subject of legal privilege. If you are not the intended recipient please notify the sender immediately and delete this e-mail. Also you may not use any information contained in it.  Legal privilege is not waived because you have read this e-mail. Views expressed in this email are not necessarily those of Brunton (NZ) Ltd.</p>
        </div>
    </body>
</html>
`
	emailActiveAccount = `
        <p>Hello <b>%s</b></p>
        <p>We just received a request to create a Brunton account using your email address: %s.</p>
        <p>To make sure your email address wasn't used fraudulently by someone else, please confirm you want the account to be activated.</p>
        <p><a href="%s">Activate your account here</a></p>
        <p>You need to activate your account within 24 hour from the time this email was sent.</p>
        <p>If you didn't use this email address to create an AT account, you can ignore this email and the account will not be activated.</p>
`
)

func ActiveAccount(url string, email string, username string, app string) (res bool, bErr *utils.BError) {
	emailHeadSubject := "account activation"
	data := make(map[string]interface{})
	data["subject"] = "Please activate your account"
	data["content"] = fmt.Sprintf(emailHeader+emailActiveAccount+emailFooter, emailHeadSubject, username, email, url)
	data["toAddress"] = [1]map[string]string{{email: username}}
	dataJson, err := json.Marshal(data)
	if err != nil {
		return false, &utils.BError{Error: err, Message: err.Error(), Code: utils.ErrorParseJson}
	}
	_, err = task.Insert(task.TaskTypeEmail, string(dataJson), app)
	if err != nil {
		return false, &utils.BError{Error: err, Message: err.Error(), Code: utils.ErrorService}
	}
	return true, nil
}

func SysAlert(app string, msg string) (bErr *utils.BError) {
	DevOps := beego.AppConfig.Strings("DevOps")
	emailHeadSubject := fmt.Sprintf("shit happens! application [%s]", app)
	data := make(map[string]interface{})
	data["subject"] = "SYSTEM ERROR!"
	data["content"] = fmt.Sprintf(emailHeader+"%s"+emailFooter, emailHeadSubject, msg)
	toAddress := make([]map[string]string, len(DevOps))
	for i := range DevOps {
		toAddress[i] = map[string]string{DevOps[i]: ""}
	}
	data["toAddress"] = toAddress
	dataJson, err := json.Marshal(data)
	if err != nil {
		return &utils.BError{Error: err, Message: err.Error(), Code: utils.ErrorParseJson}
	}
	_, err = task.Insert(task.TaskTypeEmail, string(dataJson), app)
	if err != nil {
		return &utils.BError{Error: err, Message: err.Error(), Code: utils.ErrorService}
	}
	return nil
}

func SysInfo(app string, msg string) (bErr *utils.BError) {
	DevOps := beego.AppConfig.Strings("DevOps")
	emailHeadSubject := fmt.Sprintf("important information! application [%s]", app)
	data := make(map[string]interface{})
	data["subject"] = "SYSTEM INFORMATION!"
	data["content"] = fmt.Sprintf(emailHeader+"%s"+emailFooter, emailHeadSubject, msg)
	toAddress := make([]map[string]string, len(DevOps))
	for i := range DevOps {
		toAddress[i] = map[string]string{DevOps[i]: ""}
	}
	data["toAddress"] = toAddress
	dataJson, err := json.Marshal(data)
	if err != nil {
		return &utils.BError{Error: err, Message: err.Error(), Code: utils.ErrorParseJson}
	}
	_, err = task.Insert(task.TaskTypeEmail, string(dataJson), app)
	if err != nil {
		return &utils.BError{Error: err, Message: err.Error(), Code: utils.ErrorService}
	}
	return nil
}
