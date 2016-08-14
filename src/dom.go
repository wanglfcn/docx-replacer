package main

import (
	"regexp"

	"github.com/pkg/errors"
)

type Dom struct {
	node_list []*Node
}

type Node struct {
	prefix  string
	content string
	l_tag   string
	r_tag   string
	tail    bool
}

func Parse(content, begin_tag, end_tag string) (*Dom, error) {

	var node_list []*Node

	pos := 0

	re_begin_tag := regexp.MustCompile(begin_tag)
	re_end_tag := regexp.MustCompile(end_tag)

	begin_pos := re_begin_tag.FindAllStringIndex(content, -1)
	end_pos := re_end_tag.FindAllStringIndex(content, -1)

	if len(begin_pos) != len(end_pos) {
		return nil, errors.New("left tag not eq right tag")
	}

	for index, b_pos := range begin_pos {
		e_pos := end_pos[index]
		node := &Node{
			prefix:  content[pos:b_pos[0]],
			l_tag:   content[b_pos[0]:b_pos[1]],
			content: content[b_pos[1]:e_pos[0]],
			r_tag:   content[e_pos[0]:e_pos[1]],
			tail:    false,
		}

		node_list = append(node_list, node)

		pos = e_pos[1]
	}

	node := &Node{
		prefix:  content[pos:],
		l_tag:   "",
		content: "",
		r_tag:   "",
		tail:    true,
	}
	node_list = append(node_list, node)

	return &Dom{
		node_list: node_list,
	}, nil

}

func (node *Node) set_value(value string) {
	node.content = value
}

func (dom *Dom) toString() (content string) {
	content = ""

	for _, node := range dom.node_list {
		if node.tail == false {
			content += node.prefix + node.l_tag + node.content + node.r_tag
		} else {
			content += node.prefix
		}
	}

	return
}
