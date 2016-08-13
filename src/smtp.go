package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/smtp"
	"strings"

	"github.com/andlabs/ui"
)

func send(auth smtp.Auth, from, to, subject, content string) (err error) {
	msg := []byte("To: " + to + "\r\nFrom: " + from + "<" + from + ">\r\nSubject: " +
		subject + "\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n" + content)
	err = smtp.SendMail(
		"smtp.126.com:25",
		auth,
		from,
		[]string{to},
		msg,
	)

	return
}

func send_all(data *XlsxData, temp_file, from, password string, sheet_index, to_index, sub_index int, status_label *ui.Label, self bool) (err error) {
	var max_col int
	var to string
	var subject string
	var content string
	var log_msg string

	temp_data, err := ioutil.ReadFile(temp_file)

	if err != nil {
		log_msg = fmt.Sprintf("打开模版文件失败:%s", err.Error())
		status_label.SetText(status_label.Text() + "\n" + log_msg)
		log.Print(log_msg)
		return
	}

	auth := smtp.PlainAuth(
		"",
		from,
		password,
		"smtp.126.com",
	)

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
					}

					if self {
						to = from
					}

					log.Print(content)
					err = send(auth, from, to, subject, content)

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
			status_label.SetText(status_label.Text() + "\n发送完成")
		}
	}

	return
}
