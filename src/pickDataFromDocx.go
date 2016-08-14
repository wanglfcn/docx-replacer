package main

import (
	"fmt"
	"log"

	"regexp"

	"strings"

	"io/ioutil"
	"path/filepath"

	"github.com/andlabs/ui"
)

func pick_data_from_docx(temp_file, source_path, target_file string, label *ui.Label) {

	var log_msg string
	temp_blocks, tags, err := get_tags(temp_file, label)

	if err != nil {
		return
	}

	files, err := ioutil.ReadDir(source_path)
	if err != nil {
		log_msg = fmt.Sprintf("打开目录失败:%s", err.Error())
		log.Print(log_msg)
		label.SetText(label.Text() + "\n" + log_msg)
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
				log.Print(content)
				doc.Close()
				for i := 0; i < len(tags); i++ {

					pre_pos := strings.Index(content, temp_blocks[i])

					last_pos := -1
					if pre_pos >= 0 {
						pre_pos += len(temp_blocks[i])
						if i < len(temp_blocks)-1 {
							last_pos = strings.Index(content, temp_blocks[i+1])
						}

						if last_pos > 0 {
							log.Printf("%s => %s", tags[i], strings.TrimSpace(content[pre_pos:last_pos]))
						} else {
							log.Printf("%s => %s", tags[i], strings.TrimSpace(content[pre_pos:]))
						}
					}

				}
			}

		}
	}

}

func get_tags(temp_file string, label *ui.Label) (temp_blocks, tags []string, err error) {

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
		tags = append(tags, template_content[pos[0]+2:pos[1]-2])
		start = pos[1]
	}

	if len(template_content) > start {
		temp_blocks = append(temp_blocks, strings.TrimSpace(template_content[start:]))
	}

	return
}
