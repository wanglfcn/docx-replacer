package main

import (
	"log"
	"strings"

	"github.com/andlabs/ui"
)

func email_sender_win(window *ui.Window, main *ui.Box) *ui.Box {

	var email string
	var password string
	var temp_file string
	var excel_data *XlsxData

	outer_box := ui.NewVerticalBox()

	email_box := ui.NewHorizontalBox()
	email_input := ui.NewEntry()
	email_box.Append(ui.NewLabel("邮箱:"), false)
	email_box.Append(email_input, true)
	outer_box.Append(email_box, false)
	email_input.OnChanged(func(entry *ui.Entry) {
		email = entry.Text()
	})

	passwd_box := ui.NewHorizontalBox()
	passwd_input := ui.NewEntry()
	passwd_box.Append(ui.NewLabel("密码:"), false)
	passwd_box.Append(passwd_input, true)
	outer_box.Append(passwd_box, false)

	passwd_input.OnChanged(func(entry *ui.Entry) {
		input_len := len(entry.Text())
		if input_len < len(password) {
			password = password[:input_len]
		} else if input_len > len(password) {
			password += entry.Text()[len(password):]
			entry.SetText(strings.Repeat("*", input_len))
		}
	})

	temp_box := ui.NewHorizontalBox()
	temp_text_field := ui.NewEntry()
	temp_btn := ui.NewButton("打开")
	temp_text_field.SetReadOnly(true)
	temp_box.Append(ui.NewLabel("模板文件(html):"), false)
	temp_box.Append(temp_text_field, true)
	temp_box.Append(temp_btn, false)

	outer_box.Append(temp_box, false)

	data_box := ui.NewHorizontalBox()
	data_text_field := ui.NewEntry()
	data_btn := ui.NewButton("打开")
	data_text_field.SetReadOnly(true)
	data_box.Append(ui.NewLabel("数据文件(excel):"), false)
	data_box.Append(data_text_field, true)
	data_box.Append(data_btn, false)

	outer_box.Append(data_box, false)

	select_first_box := ui.NewHorizontalBox()

	sheet_box := ui.NewHorizontalBox()
	sheet_comb := ui.NewCombobox()
	sheet_box.Append(ui.NewLabel("sheet:"), false)
	sheet_box.Append(sheet_comb, true)

	select_first_box.Append(sheet_box, true)

	email_pos_box := ui.NewHorizontalBox()
	email_pos_comb := ui.NewCombobox()
	email_pos_box.Append(ui.NewLabel("email:"), false)
	email_pos_box.Append(email_pos_comb, true)

	select_first_box.Append(email_pos_box, true)
	outer_box.Append(select_first_box, false)

	select_sec_box := ui.NewHorizontalBox()

	subject_box := ui.NewHorizontalBox()
	subject_comb := ui.NewCombobox()
	subject_box.Append(ui.NewLabel("主题:"), false)
	subject_box.Append(subject_comb, true)

	select_sec_box.Append(subject_box, true)

	cc_box := ui.NewHorizontalBox()
	cc_comb := ui.NewCombobox()
	cc_box.Append(ui.NewLabel("抄送:"), false)
	cc_box.Append(cc_comb, true)

	select_sec_box.Append(cc_box, true)

	outer_box.Append(select_sec_box, false)

	send_box := ui.NewHorizontalBox()
	send_to_self_btn := ui.NewButton("发送给自己")
	exec_btn := ui.NewButton("执行")
	send_box.Append(send_to_self_btn, false)
	send_box.Append(exec_btn, false)
	outer_box.Append(send_box, false)

	stats_label := ui.NewLabel("")
	outer_box.Append(stats_label, true)

	return_btn := ui.NewButton("返回")
	outer_box.Append(return_btn, false)

	return_btn.OnClicked(func(*ui.Button) {
		window.SetChild(main)
	})

	temp_btn.OnClicked(
		func(*ui.Button) {
			temp_file = ui.OpenFile(window)
			temp_text_field.SetText(temp_file)
		},
	)

	data_btn.OnClicked(
		func(*ui.Button) {
			file_name := ui.OpenFile(window)
			var err error
			excel_data, err = NewXlsxData(file_name)
			if err != nil {
				data_text_field.SetText(err.Error())
			} else {
				sheets := excel_data.getSheets()
				sheet_comb.Reset()
				for _, s_name := range sheets {
					sheet_comb.Append(s_name)
				}
				sheet_comb.SetSelected(0)

				data_text_field.SetText(file_name)

				sel := sheet_comb.Selected()
				log.Printf("sel: %d, %v", sel, excel_data)
				titles := excel_data.getTitles(sel)
				log.Printf("titles: %v", titles)

				email_pos_comb.Reset()
				subject_comb.Reset()
				cc_comb.Reset()
				for _, title := range titles {
					email_pos_comb.Append(title)
					subject_comb.Append(title)
					cc_comb.Append(title)
				}
				email_pos_comb.SetSelected(0)
				subject_comb.SetSelected(0)
			}
		},
	)

	sheet_comb.OnSelected(
		func(comb *ui.Combobox) {
			sel := comb.Selected()
			log.Printf("sel: %d, %v", sel, excel_data)
			titles := excel_data.getTitles(sel)
			log.Printf("titles: %v", titles)

			email_pos_comb.Reset()
			subject_comb.Reset()
			cc_comb.Reset()
			for _, title := range titles {
				email_pos_comb.Append(title)
				subject_comb.Append(title)
				cc_comb.Append(title)
			}
			email_pos_comb.SetSelected(0)
			subject_comb.SetSelected(0)
		},
	)

	send_to_self_btn.OnClicked(func(*ui.Button) {
		go func() {
			send_to_self_btn.Disable()
			exec_btn.Disable()
			stats_label.SetText("开始发送邮件...")
			send_all(excel_data, temp_file, email, password, sheet_comb.Selected(), email_pos_comb.Selected(), subject_comb.Selected(), cc_comb.Selected(), stats_label, true)
			send_to_self_btn.Enable()
			exec_btn.Enable()
		}()
	})

	exec_btn.OnClicked(func(*ui.Button) {
		go func() {
			send_to_self_btn.Disable()
			exec_btn.Disable()
			stats_label.SetText("开始发送邮件...")
			send_all(excel_data, temp_file, email, password, sheet_comb.Selected(), email_pos_comb.Selected(), subject_comb.Selected(), cc_comb.Selected(), stats_label, false)
			send_to_self_btn.Enable()
			exec_btn.Enable()
		}()
	})

	return outer_box

}
