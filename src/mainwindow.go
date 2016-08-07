package main

import (
	"github.com/andlabs/ui"
	"log"
)

func mainwindow() {

	var window *ui.Window
	var excel_data *XlsxData
	var temp_file string
	var save_path string

	outer_box := ui.NewVerticalBox();

	temp_box := ui.NewHorizontalBox();
	temp_text_field := ui.NewEntry()
	temp_text_field.SetReadOnly(true)
	temp_btn := ui.NewButton("打开")
	temp_box.Append(ui.NewLabel("模板文件:"), false)
	temp_box.Append(temp_text_field, true)
	temp_box.Append(temp_btn, false)

	outer_box.Append(temp_box, false)

	data_box := ui.NewHorizontalBox()
	data_text_field := ui.NewEntry()
	data_text_field.SetReadOnly(true)
	data_btn := ui.NewButton("打开")
	data_box.Append(ui.NewLabel("数据文件:"), false)
	data_box.Append(data_text_field, true)
	data_box.Append(data_btn, false)

	outer_box.Append(data_box, false)

	export_box := ui.NewHorizontalBox()
	export_text_field := ui.NewEntry()
	export_text_field.SetReadOnly(false)
	export_btn := ui.NewButton("打开")
	export_box.Append(ui.NewLabel("导出位置:"), false)
	export_box.Append(export_text_field, true)
	export_box.Append(export_btn, false)

	outer_box.Append(export_box, false)

	select_box := ui.NewHorizontalBox()

	sheet_box := ui.NewHorizontalBox()
	sheet_comb := ui.NewCombobox()
	sheet_comb.Append("test")
	sheet_box.Append(ui.NewLabel("sheet:"), false)
	sheet_box.Append(sheet_comb, true)

	select_box.Append(sheet_box, true)

	file_name_box := ui.NewHorizontalBox()
	file_name_comb := ui.NewCombobox()
	file_name_box.Append(ui.NewLabel("file_name:"), false)
	file_name_box.Append(file_name_comb, true)

	select_box.Append(file_name_box, true)

	outer_box.Append(select_box, false)

	exec_btn := ui.NewButton("执行")
	outer_box.Append(exec_btn, false)

	exec_btn.OnClicked(
		func(*ui.Button) {
			excel_data.replace(sheet_comb.Selected(), file_name_comb.Selected(), temp_file , save_path)
		},
	)


	temp_btn.OnClicked(
		func(*ui.Button) {
			temp_file = ui.OpenFile(window)
			temp_text_field.SetText(temp_file)
		},
	)

	export_text_field.OnChanged(
		func(e *ui.Entry) {
			save_path = e.Text()
		},
	)

	export_btn.OnClicked(
		func(*ui.Button) {
			save_path = ui.OpenDir(window)
			export_text_field.SetText(save_path)
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

				file_name_comb.Reset()
				for _, title := range titles {
					file_name_comb.Append(title)
				}
				file_name_comb.SetSelected(0)
			}
		},
	)

	sheet_comb.OnSelected(
		func(comb *ui.Combobox) {
			sel := comb.Selected()
			log.Printf("sel: %d, %v", sel, excel_data)
			titles := excel_data.getTitles(sel)
			log.Printf("titles: %v", titles)

			file_name_comb.Reset()
			for _, title := range titles {
				file_name_comb.Append(title)
			}
			file_name_comb.SetSelected(0)
		},
	)


	window = ui.NewWindow("小鹅专用", 500, 250, false)
	window.SetChild(outer_box)
	window.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})
	window.Show()

}

