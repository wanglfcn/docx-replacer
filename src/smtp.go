package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/andlabs/ui"
	"gopkg.in/gomail.v2"
)

func compose_email(from, to, subject, cc, content string) (msg *gomail.Message) {

	msg = gomail.NewMessage()
	msg.SetHeader("From", from)
	msg.SetHeader("To", to)
	if len(cc) > 3 {
		msg.SetHeader("Cc", cc)
	}
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", content)

	return
}

func send_all(data *XlsxData, temp_file, from, password string, sheet_index, to_index, sub_index, cc_index int, status_label *ui.Label, self bool) (err error) {
	var max_col int
	var to string
	var subject string
	var content string
	var cc string
	var log_msg string

	temp_data, err := ioutil.ReadFile(temp_file)

	if err != nil {
		log_msg = fmt.Sprintf("打开模版文件失败:%s", err.Error())
		status_label.SetText(status_label.Text() + "\n" + log_msg)
		log.Print(log_msg)
		return
	}

	dialog := gomail.NewDialer("smtp.163.com", 465, from, password)

	sender, err := dialog.Dial()

	if err != nil {

		log_msg = fmt.Sprintf("连接服务器错误: %s", err.Error())
		status_label.SetText(status_label.Text() + "\n" + log_msg)
		log.Print(log_msg)
	}
	defer sender.Close()

	col_name := make(map[int]string)

	for s_i, sheet := range data.fh.Sheets {

		if s_i == sheet_index {

			for r_i, row := range sheet.Rows {
				if r_i == 0 {
					for c_i, col := range row.Cells {
						col_name[c_i] = col.Value
						max_col = c_i
					}
				} else {
					content = string(temp_data)
					to = ""
					subject = ""
					content = string(temp_data)

					for c_i, col := range row.Cells {
						if c_i > max_col {
							break
						}

						content = strings.Replace(content, "{{"+col_name[c_i]+"}}", col.Value, -1)
						if c_i == to_index {
							to = col.Value
						}

						if c_i == sub_index {
							subject = col.Value
						}

						if c_i == cc_index {
							cc = col.Value
						}
					}

					if self {
						to = from
						if len(cc) > 3 {
							cc = from
						}
					}

					log.Print(content)
					if len(to) > 3 {

						msg := compose_email(from, to, subject, cc, content)
						err = gomail.Send(sender, msg)

						if err != nil {
							log_msg = fmt.Sprintf("[%d]发送失败:%s,原因: %s", r_i, to, err.Error())
						} else {
							log_msg = fmt.Sprintf("[%d]发送成功:%s", r_i, to)
						}
						status_label.SetText(status_label.Text() + "\n" + log_msg)
						log.Print(log_msg)

						if self {
							return
						}
					}
				}
			}
			status_label.SetText(status_label.Text() + "\n发送完成")
		}
	}

	return
}
