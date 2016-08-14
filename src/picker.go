package main

import "github.com/andlabs/ui"

func picker_win(window *ui.Window, main *ui.Box) *ui.Box {

	var temp_file string
	var source_path string
	var save_file string

	outer_box := ui.NewVerticalBox()

	temp_box := ui.NewHorizontalBox()
	temp_text_field := ui.NewEntry()
	temp_text_field.SetReadOnly(true)
	temp_btn := ui.NewButton("打开")
	temp_box.Append(ui.NewLabel("模板文件(word):"), false)
	temp_box.Append(temp_text_field, true)
	temp_box.Append(temp_btn, false)

	outer_box.Append(temp_box, false)

	source_box := ui.NewHorizontalBox()
	source_text_field := ui.NewEntry()
	source_text_field.SetReadOnly(true)
	source_btn := ui.NewButton("打开")
	source_box.Append(ui.NewLabel("数据位置:"), false)
	source_box.Append(source_text_field, true)
	source_box.Append(source_btn, false)

	outer_box.Append(source_box, false)

	data_box := ui.NewHorizontalBox()
	data_text_field := ui.NewEntry()
	data_text_field.SetReadOnly(true)
	data_btn := ui.NewButton("打开")
	data_box.Append(ui.NewLabel("保存文件(excel):"), false)
	data_box.Append(data_text_field, true)
	data_box.Append(data_btn, false)

	outer_box.Append(data_box, false)

	exec_btn := ui.NewButton("执行")
	outer_box.Append(exec_btn, false)

	status_label := ui.NewLabel("")
	outer_box.Append(status_label, true)

	temp_btn.OnClicked(func(*ui.Button) {
		temp_file = ui.OpenFile(window)
		temp_text_field.SetText(temp_file)
	})

	source_btn.OnClicked(func(*ui.Button) {
		source_path = ui.OpenDir(window)
		source_text_field.SetText(source_path)
	})

	data_btn.OnClicked(func(*ui.Button) {
		save_file = ui.SaveFile(window)
		data_text_field.SetText(save_file)
	})

	exec_btn.OnClicked(func(*ui.Button) {
		go func() {
			exec_btn.Disable()
			status_label.SetText("开始...")
			pick_data_from_docx(temp_file, source_path, save_file, status_label)
			exec_btn.Enable()
		}()
	})

	return_btn := ui.NewButton("返回")

	return_btn.OnClicked(func(*ui.Button) {
		window.SetChild(main)
	})

	outer_box.Append(return_btn, false)

	return outer_box

}
