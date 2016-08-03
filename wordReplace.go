package main

import (
	"os"
	"github.com/lestrrat/go-libxml2/parser"
	"github.com/lestrrat/go-libxml2/types"
	"log"
	"strings"
	"archive/zip"
	"io"
	"bytes"
	"io/ioutil"
	"errors"
)

type NodeList struct {
	node 	types.Node
	prefix 	string
	suffix 	string
	start 	int
	end 	int
}

type ReplaceDocx struct {
	zipReader 	*zip.ReadCloser
	dom 		types.Document
	nodelist 	[]NodeList
	content_str	string
	content_xml	string
	replace_map 	map[string][]NodeList
}

func ParseDocXml(doc_xml string) (dom types.Document, nodelist []NodeList, content string, err error) {
	doc_parse := parser.New()

	dom, err = doc_parse.ParseString(doc_xml)

	if err != nil {
		log.Printf("parse xml with error %v", err)
		return
	}

	content_length := 0

	dom.Walk(
		func(node types.Node) error {
			if strings.EqualFold(node.NodeName(), "w:t") {
				content += node.TextContent()
				new_length := len(content) - 1
				nodelist = append(nodelist,
					NodeList{
						node:node,
						start:content_length,
						end: new_length,
						prefix: "",
						suffix: "",
					},
				)
				content_length = new_length + 1
			}
			return nil
	})

	return
}

func (r *ReplaceDocx)find_nodes(start, end int) (node NodeList, ok bool) {

	var targetList []NodeList

	for _, node := range r.nodelist {
		if (node.start <= end && node.end >= start) {
			if node.start < start {
				node.prefix = node.node.TextContent()[:start - node.start]
			}

			if node.end > end {
				node.suffix = node.node.TextContent()[end - node.start + 1:]
			}

			node.node.SetNodeValue(node.prefix + node.suffix)
			targetList = append(targetList, node)
		}
	}

	if len(targetList) == 0 {
		ok = false
	} else {
		ok = true
		node = targetList[len(targetList)/2]
	}

	return
}

func (r *ReplaceDocx)set_reapace_map(replaces []string) {

	replace_map := make(map[string][]NodeList)

	for _, replace_str := range replaces {
		start_pos := 0
		var node_list []NodeList
		for {
			find_start := strings.Index(r.content_str[start_pos:], replace_str)
			if find_start >= 0 {
				find_start += start_pos
				find_end := find_start + len(replace_str) - 1
				node, ok := r.find_nodes(find_start, find_end)
				if ok {
					node_list = append(node_list, node)
				}

				start_pos = find_end + 1
			} else {
				break
			}

			if len(node_list) > 0 {
				replace_map[replace_str] = node_list
			}
		}

	}
	r.replace_map = replace_map
}

func(r *ReplaceDocx)replace(orig, target string) {
	if rep_map, ok := r.replace_map[orig]; ok {

		for _, node := range rep_map {
			node.node.SetNodeValue(node.prefix + target + node.suffix)
		}
	}
}

func (r *ReplaceDocx) Editable() *Docx {
	return &Docx{
		files:   r.zipReader.File,
		content: r.dom.String(),
	}
}

func (r *ReplaceDocx) Close() error {
	return r.zipReader.Close()
}

type Docx struct {
	files   []*zip.File
	content string
}

func (d *Docx) WriteToFile(path string) (err error) {
	var target *os.File
	target, err = os.Create(path)
	if err != nil {
		return
	}
	defer target.Close()
	err = d.Write(target)
	return
}

func (d *Docx) Write(ioWriter io.Writer) (err error) {
	w := zip.NewWriter(ioWriter)
	for _, file := range d.files {
		var writer io.Writer
		var readCloser io.ReadCloser

		writer, err = w.Create(file.Name)
		if err != nil {
			return err
		}
		readCloser, err = file.Open()
		if err != nil {
			return err
		}
		if file.Name == "word/document.xml" {
			writer.Write([]byte(d.content))
		} else {
			writer.Write(streamToByte(readCloser))
		}
	}
	w.Close()
	return
}

func ReadDocxFile(path string) (*ReplaceDocx, error) {
	reader, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	content_xml, err := readText(reader.File)
	if err != nil {
		return nil, err
	}

	dom, nodelist, content_str, err := ParseDocXml(content_xml)
	if err != nil {
		return nil, err
	}

	return &ReplaceDocx{zipReader: reader, dom: dom, nodelist: nodelist, content_str: content_str, content_xml: content_xml}, nil
}

func readText(files []*zip.File) (text string, err error) {
	var documentFile *zip.File
	documentFile, err = retrieveWordDoc(files)
	if err != nil {
		return text, err
	}
	var documentReader io.ReadCloser
	documentReader, err = documentFile.Open()
	if err != nil {
		return text, err
	}

	text, err = wordDocToString(documentReader)
	return
}

func wordDocToString(reader io.Reader) (string, error) {
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func retrieveWordDoc(files []*zip.File) (file *zip.File, err error) {
	for _, f := range files {
		if f.Name == "word/document.xml" {
			file = f
		}
	}
	if file == nil {
		err = errors.New("document.xml file not found")
	}
	return
}

func streamToByte(stream io.Reader) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.Bytes()
}
