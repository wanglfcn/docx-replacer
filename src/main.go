package main

import (
	"log"
	"github.com/andlabs/ui"
)

func main() {

	err := ui.Main(main_win)


	if err != nil {
		log.Fatalf("fatal error %v", err)
	}
}

func main_win(){
	window := ui.NewWindow("小鹅专用", 500, 250, false)


	outer_box 	:= ui.NewVerticalBox();
	replacer_btn 	:= ui.NewButton("生成合同")
	picker_btn	:= ui.NewButton("提取数据")
	email_sender_btn:= ui.NewButton("发送邮件")

	replacer 	:= replacer_win(window, outer_box)
	picker 		:= picker_win(window, outer_box)
	email_sender 	:= email_sender_win(window, outer_box)

	outer_box.Append(replacer_btn, false)
	outer_box.Append(picker_btn, false)
	outer_box.Append(email_sender_btn, false)

	replacer_btn.OnClicked(func(*ui.Button) {
		window.SetChild(replacer)
	})

	picker_btn.OnClicked(func(*ui.Button) {
		window.SetChild(picker)
	})

	email_sender_btn.OnClicked(func(*ui.Button) {
		window.SetChild(email_sender)
	})

	window.SetChild(outer_box)
	window.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})
	window.Show()
}
