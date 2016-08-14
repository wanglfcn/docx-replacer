package main

import (
	"fmt"
	"log"

	"regexp"

	"strings"

	"io/ioutil"
	"path/filepath"

	"github.com/andlabs/ui"
	"github.com/tealeg/xlsx"
)

func pick_data_from_docx(temp_file, source_path, target_file string, label *ui.Label) {

	var log_msg string
	var xlsx_file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell

	xlsx_file = xlsx.NewFile()
	sheet, err := xlsx_file.AddSheet("Sheet1")
	row = sheet.AddRow()
	temp_blocks, err := get_blocks(temp_file, row, label)

	block_len := len(temp_blocks)

	if err != nil {
		return
	}

	files, err := ioutil.ReadDir(source_path)
	if err != nil {
		log_msg = fmt.Sprintf("打开目录失败:%s", err.Error())
		log.Print(log_msg)
		label.SetText(label.Text() + "\n" + log_msg)
	}

	if err != nil {

	}

	for _, file := range files {

		if strings.HasSuffix(file.Name(), "docx") {
			doc, err := ReadDocxFile(filepath.Join(source_path, file.Name()))

			if err != nil {
				log_msg = fmt.Sprintf("打开文件[%s]失败:%s", file.Name(), err.Error())
				log.Print(log_msg)
				label.SetText(label.Text() + "\n" + log_msg)

			} else {
				log_msg = fmt.Sprintf("处理文件[%s]", file.Name())
				log.Print(log_msg)
				label.SetText(label.Text() + "\n" + log_msg)

				content := doc.content_str
				doc.Close()
				row = sheet.AddRow()
				for i := 0; i < block_len; i++ {

					pre_pos := strings.Index(content, temp_blocks[i])

					last_pos := -1
					value := ""
					if pre_pos >= 0 {
						pre_pos += len(temp_blocks[i])
						if i < block_len-1 {
							last_pos = strings.Index(content, temp_blocks[i+1])
						}

						if last_pos > 0 {

							value = strings.TrimSpace(content[pre_pos:last_pos])
						} else {
							value = strings.TrimSpace(content[pre_pos:])
						}
					}

					cell = row.AddCell()
					cell.SetString(value)
				}
			}

		}
	}

	if strings.HasSuffix(target_file, ".xlsx") == false {
		target_file += ".xlsx"
	}

	xlsx_file.Save(target_file)
	log_msg = "处理完成"
	log.Print(log_msg)
	label.SetText(label.Text() + "\n" + log_msg)

}

func get_blocks(temp_file string, titles *xlsx.Row, label *ui.Label) (temp_blocks []string, err error) {

	template_doc, err := ReadDocxFile(temp_file)

	var log_msg string
	var template_content string

	if err != nil {
		log_msg = fmt.Sprintf("打开模版文件错误: %s", err.Error())
		label.SetText(label.Text() + "\n" + log_msg)
		log.Print(log_msg)
		return
	}

	template_content = template_doc.content_str
	template_doc.Close()

	find_reg, err := regexp.Compile("{{.*?}}")

	if err != nil {
		log.Printf("compile regex with error: %s", err.Error())
	}

	positons := find_reg.FindAllStringIndex(template_content, -1)

	if positons == nil {
		log_msg = "没有发现需要查找的变量"
	}

	start := 0
	for _, pos := range positons {
		temp_blocks = append(temp_blocks, strings.TrimSpace(template_content[start:pos[0]]))
		cell := titles.AddCell()
		cell.SetString(template_content[pos[0]+2 : pos[1]-2])
		start = pos[1]
	}

	if len(template_content) > start {
		temp_blocks = append(temp_blocks, strings.TrimSpace(template_content[start:]))
	}

	return
}
