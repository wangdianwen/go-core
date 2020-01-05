package scripts

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/wangdianwen/go-core/middlewares/logger"
	task "github.com/wangdianwen/go-core/middlewares/task/models"
	"gopkg.in/gomail.v2"
)

func EmailTask() error {
	Logger, bErr := logger.NewLogger("middleware/email/EmailTask.log")
	if bErr != nil {
		return bErr.Error
	}

	bErr = Logger.Info("[sending email]start task...")

	list, err := task.UnDoTasks(task.TaskTypeEmail, 30)
	if err != nil {
		errStr := fmt.Sprintf("[sending email] create task error: %s", err.Error())
		Logger.Error(errStr)
		return errors.New(errStr)
	}
	Logger.Info("[sending email] get total list number %d", len(list))
	if len(list) == 0 {
		Logger.Info("[sending email] found no data, exit")
		return nil
	}

	host := beego.AppConfig.String("SMTPHost")
	port, _ := beego.AppConfig.Int("SMTPPort")
	sender := beego.AppConfig.String("SenderName")
	pass := beego.AppConfig.String("SenderPass")

	d := gomail.NewDialer(host, port, sender, pass)
	s, err := d.Dial()
	if err != nil {
		errStr := fmt.Sprintf("[sending email] dial server error: %s", err.Error())
		Logger.Info(errStr)
		return errors.New(errStr)
	}
	m := gomail.NewMessage()
	for _, r := range list {
		var dat map[string]interface{}
		err = json.Unmarshal([]byte(r.Data), &dat)
		if err != nil {
			errStr := fmt.Sprintf("[sending email][taskid %d] parse json error: %s", r.Id, err.Error())
			_, _ = task.TaskFail(r.Id, errStr)
			Logger.Error(errStr)
			return errors.New(errStr)
		}
		toAddress := dat["toAddress"].([]interface{})
		subject := dat["subject"].(string)
		content := dat["content"].(string)

		addresses := make([]string, len(toAddress))
		for i := range toAddress {
			for k, v := range toAddress[i].(map[string]interface{}) {
				addresses[i] = m.FormatAddress(k, v.(string))
			}
		}
		m.SetHeader("From", sender)
		m.SetHeader("To", addresses...)
		m.SetHeader("Subject", subject)
		// logo pics
		m.Embed("assert/images/bizex_logo.png")
		m.Embed("assert/images/brunton_footer.jpg")
		m.SetBody("text/html", content)
		if err := gomail.Send(s, m); err != nil {
			m.Reset()
			errStr := fmt.Sprintf("[sending email][taskid %d] Could not send email to %s, reason: %s", r.Id, toAddress, err.Error())
			_, _ = task.TaskFail(r.Id, errStr)
			Logger.Error(errStr)
			return errors.New(errStr)
		}
		m.Reset()
		_, _ = task.TaskSuccess(r.Id)
		Logger.Info("[sending email][taskid %d] send email success", r.Id)
	}
	Logger.Info("[sending email] all jobs done!")
	return nil
}
