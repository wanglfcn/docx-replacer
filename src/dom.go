package main

import (
	"strings"
)

type Dom struct {
	node_list 	[]*Node
	beginTag	string
	endTag		string
}

type Node struct {
	prefix 	string
	content string
	tail 	bool
}

func Parse(content, begin_tag, end_tag string)(*Dom, error) {

	var node_list []*Node

	pos := 0

	for {
		if pos = strings.Index(content, end_tag); pos > 0 {
			pos += len(end_tag)
			node := generate_node(content[:pos], begin_tag, end_tag, false)
			node_list = append(node_list, node)
			content = content[pos:]
		} else {
			node := generate_node(content, "", "", true)
			node_list = append(node_list, node)
			break
		}
	}

	return &Dom{
		node_list: 	node_list,
		beginTag:  	begin_tag,
		endTag: 	end_tag,
	}, nil

}

func generate_node(value, start_tag, end_tag string, tail bool) (*Node) {
	prefix_content := ""
	content := ""

	if start_pos := strings.Index(value, start_tag); start_pos >= 0 && len(start_tag) > 0 {
		prefix_content = value[: start_pos]
		content = value[start_pos + len(start_tag) : len(value) - len(end_tag)]
	} else {
		prefix_content = value
	}

	return &Node{
		prefix: 	prefix_content,
		content: 	content,
		tail: 		tail,
	}
}

func (node *Node)set_value(value string) {
	node.content = value
}

func (dom *Dom)toString() (content string) {
	content = ""

	for _, node := range dom.node_list {
		if node.tail == false {
			content += node.prefix + dom.beginTag + node.content + dom.endTag
		} else {
			content += node.prefix
		}
	}

	return
}
