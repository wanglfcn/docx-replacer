package main

import "github.com/andlabs/ui"

func picker_win(window *ui.Window, main *ui.Box)(*ui.Box) {

	outer_box := ui.NewVerticalBox()
	label := ui.NewLabel("收集数据")
	return_btn := ui.NewButton("返回")

	return_btn.OnClicked(func(*ui.Button) {
		window.SetChild(main)
	})

	outer_box.Append(label, false)
	outer_box.Append(return_btn, false)

	return outer_box

}
