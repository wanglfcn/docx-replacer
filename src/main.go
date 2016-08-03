package main

import (
	"github.com/nguyenthenguyen/docx"
	"github.com/tealeg/xlsx"
	"log"
	"path/filepath"
	"flag"
	"fmt"
	"strings"
	"github.com/andlabs/ui"
)


func main() {

	//template_file := flag.String("temp", "", "template file path")
	//data_file := flag.String("data", "", "data file path")

	flag.Parse()

	err := ui.Main(mainwindow)
	if err != nil {
		log.Fatalf("fatal error %v", err)
	}
	//write_file(*template_file, *data_file)

}


func write_file(template_file, data_file string) {
	template_doc, err := docx.ReadDocxFile(template_file)

	if err != nil {
		log.Fatalf("open template file with err: %v", err)
	}
	defer template_doc.Close()

	data_xlsx, err := xlsx.OpenFile(data_file)
	if err != nil {
		log.Fatalf("open template file with err: %v", err)
	}

	titles := make(map[int] string)
	for _, sheet := range data_xlsx.Sheets {
		log.Printf("use sheed %s", sheet.Name)
		for index, row := range sheet.Rows {
			if index == 0 {
				for i, col := range row.Cells {
					titles[i], _ = col.String()
					titles[i] = strings.TrimSpace(titles[i])
					fmt.Printf("title: %d:\t%s\n", i, titles[i])
				}
			} else {
				new_doc := template_doc.Editable()

				file_name := ""

				for i, col := range row.Cells {
					value, _ := col.String()

					if i == 0 {
						file_name = value
					}

					fmt.Printf("replace %s:\t%s\n", titles[i], value)
					new_doc.Replace(titles[i], value, -1)
				}

				new_path := fmt.Sprintf("%s/%s.docx", filepath.Dir(data_file), file_name)
				log.Printf("write new file %s", new_path)

				new_doc.WriteToFile(new_path)
			}
		}
		break
	}
}
