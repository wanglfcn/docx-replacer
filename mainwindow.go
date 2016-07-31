package main

import "github.com/andlabs/ui"

func mainwindow() {

	var window *ui.Window

	outer_box := ui.NewVerticalBox();

	temp_box := ui.NewHorizontalBox();
	temp_text_field := ui.NewEntry()
	temp_text_field.SetReadOnly(true)
	temp_btn := ui.NewButton("打开")
	temp_box.Append(ui.NewLabel("模板文件:"), false)
	temp_box.Append(temp_text_field, true)
	temp_box.Append(temp_btn, false)

	temp_btn.OnClicked(
		func(*ui.Button) {
			file_name := ui.OpenFile(window)
			temp_text_field.SetText(file_name)
		},
	)

	outer_box.Append(temp_box, false)

	data_box := ui.NewHorizontalBox()
	data_text_field := ui.NewEntry()
	data_text_field.SetReadOnly(true)
	data_btn := ui.NewButton("打开")
	data_box.Append(ui.NewLabel("数据文件:"), false)
	data_box.Append(data_text_field, true)
	data_box.Append(data_btn, false)

	data_btn.OnClicked(
		func(*ui.Button) {
			file_name := ui.OpenFile(window)
			data_text_field.SetText(file_name)
		},
	)

	outer_box.Append(data_box, false)

	export_box := ui.NewHorizontalBox()
	export_text_field := ui.NewEntry()
	export_text_field.SetReadOnly(true)
	export_btn := ui.NewButton("打开")
	export_box.Append(ui.NewLabel("导出位置:"), false)
	export_box.Append(export_text_field, true)
	export_box.Append(export_btn, false)

	export_btn.OnClicked(
		func(*ui.Button) {
			file_name := ui.SaveFile(window)
			export_text_field.SetText(file_name)
		},
	)

	outer_box.Append(export_box, false)

	select_box := ui.NewHorizontalBox()

	sheet_box := ui.NewHorizontalBox()
	sheet_comb := ui.NewCombobox()
	sheet_comb.Append("小鹅")
	sheet_comb.Append("刘鹅")
	sheet_comb.Append("花鹅")
	sheet_box.Append(ui.NewLabel("sheet:"), false)
	sheet_box.Append(sheet_comb, true)

	select_box.Append(sheet_box, true)

	file_name_box := ui.NewHorizontalBox()
	file_name_comb := ui.NewCombobox()
	file_name_comb.Append("小鹅")
	file_name_comb.Append("刘鹅")
	file_name_comb.Append("花鹅")
	file_name_box.Append(ui.NewLabel("file_name:"), false)
	file_name_box.Append(file_name_comb, true)

	select_box.Append(file_name_box, true)

	outer_box.Append(select_box, false)

	window = ui.NewWindow("小鹅专用", 500, 250, false)
	window.SetChild(outer_box)
	window.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})
	window.Show()

}
