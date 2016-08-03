package main

import (
	"github.com/tealeg/xlsx"
	"log"
	"strings"
	"fmt"
	"path/filepath"
)

type XlsxData struct {
	file_path 	string
	fh 		*xlsx.File
}

func NewXlsxData(file_path string) (*XlsxData, error) {

	data_xlsx, err := xlsx.OpenFile(file_path)

	if err != nil {
		return nil, err
	} else {
		return &XlsxData{file_path: file_path, fh: data_xlsx}, nil
	}
}

func (xlsx *XlsxData) getSheets() ([]string) {

	var result []string

	sheets := xlsx.fh.Sheets
	for _, sheet:= range sheets {
		result = append(result, sheet.Name)
	}

	return result
}

func (xlsx *XlsxData) getTitles(sheet_index int) ([]string){

	var result []string

	for i, sheet := range xlsx.fh.Sheets {
		if i == sheet_index {
			log.Printf("i=%d", i)
			for index, row := range sheet.Rows {
				if index == 0 {
					for _, col := range row.Cells {
						title, _ := col.String()
						result = append(result, title)
					}
					return result
				}
			}
		}
	}

	return result
}

func (xlsx *XlsxData) replace(sheet_index, name_index int, template, save_path string) {

	template_doc, err := ReadDocxFile(template)

	if err != nil {
		log.Fatalf("open template file with err: %v", err)
	}
	defer template_doc.Close()

	titles := make(map[int] string)

	var max_col int

	for index, sheet := range xlsx.fh.Sheets {
		if index == sheet_index {

			for r, row := range sheet.Rows {
				if r== 0 {
					var replace_strs []string

					for i, col := range row.Cells {
						titles[i], _ = col.String()
						titles[i] = "{{" + strings.TrimSpace(titles[i]) + "}}"
						replace_strs = append(replace_strs, titles[i])
						max_col = i
						fmt.Printf("title: %d:\t%s\n", i, titles[i])
					}
					template_doc.set_reapace_map(replace_strs)
				} else {

					file_name := ""

					for i, col := range row.Cells {
						if i > max_col {
							break
						}

						value, _ := col.String()

						if i == name_index {
							file_name = value
						}

						fmt.Printf("replace %s:\t%s\n", titles[i], value)
						template_doc.replace(titles[i], value)
					}

					new_doc := template_doc.Editable()

					new_path := fmt.Sprintf("%s/%s.docx", filepath.ToSlash(save_path), file_name)
					log.Printf("write new file %s", filepath.FromSlash(new_path))

					new_doc.WriteToFile(new_path)
				}
			}
		}
	}
}