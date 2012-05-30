package main

import (
	"fmt"
)

type Node struct {
	children	[]*Node
	chars		string
	isEnd		bool
}

func (node *Node) insert(foo string) {

	//fmt.Printf("\ntrying to insert '%v' into '%v'...\n", foo, node.chars)
	//fmt.Printf("children: %v\n", len(node.children))
	// iterate over children, for each check if the first char of chars is the same as in foo
	n := -1
	for i, child := range node.children {
		//fmt.Printf("child: %v\n", child.chars)
		if len(child.chars) > 0 && len(foo) > 0 && child.chars[0] == foo[0] {
			// remember child
			n = i
			break
		}
	}
	
	//fmt.Printf("n is %v...\n", n)
	
	// if n still -1, then there is no child starting with the same letter
	if n == -1 {
		//fmt.Printf("creating new node '%v' after '%v'...\n", foo, node.chars)
		// make new end node
		newChild := &Node{nil, foo, true}
		node.children = append(node.children, newChild)
		return
	}
	
	child := node.children[n]
	
	x := 0
	for foo[x] == child.chars[x] {
		x = x + 1
		if x == len(child.chars) || x == len(foo) {
			break
		}
	}
	
	// foo[:x] is now the substring of foo contained in child.chars [e.g. only "ab" from "abab" is contained in "abba"]
	//fmt.Printf("matching chars: %v\n", foo[:x])
	//fmt.Printf("not matching chars: %v\n", foo[x:])
	
	// don't create new child if shared chars = child.chars
	if foo[:x] == child.chars {
		child.insert(foo[x:])
		return
	}
	// create new child node to replace current child, keep the .isEnd value from current child
	//fmt.Printf("creating new node '%v' after '%v'...\n", foo[:x], node.chars)
	newChild := &Node{nil, foo[:x], child.isEnd}
	
	// shorten old chars to the non-shared part
	//fmt.Printf("shortening old node '%v' to '%v'...\n", child.chars, child.chars[x:])
	child.chars = child.chars[x:]
	
	// replace child of current node with new child
	node.children[n] = newChild
	
	// old child is now a child of the new child
	newChild.children = append(newChild.children, child)
	
	
	// insert non-shared string into new child [go into node "ab" and create child node "ab" next to "ba"]
	newChild.insert(foo[x:])
}



func (node *Node) find(foo string) bool {
	
	// check if the end of our string is found
	// if chars equal, but node is no end node, return false, else return true
	if foo==node.chars {
		if node.isEnd {
			return true
		} else {
			return false
		}
	}
	
	// if we are not at the end yet, check if the current node has children
	// if not, stop
	if len(node.children) == 0 {
		return false
	}
	
	// find matching chars both foo and node.chars begin with
	// x is the index of the first char not equal
	x := 0
	if len(node.chars) > 0 {		
		for foo[x] == node.chars[x] {
			x = x + 1
			if x == len(node.chars) || x == len(foo) {
				break
			}
		}
	}
	
	// foo[:x] is now the substring of foo contained in node.chars [e.g. only "ab" from "abab" is contained in "abba"]
	//fmt.Printf("matching chars: %v\n", foo[:x])
	//fmt.Printf("not matching chars: %v\n", foo[x:])
	
	// iterate over children and continue find in the one starting with the correct letter
	for _, child := range node.children {
		// only search child with the correct first letter
		if child.chars[0] == foo[x] {
			return child.find(foo[x:])
		}
	}
	
	return false
}
	
func (node *Node) print(level int) {
	x := level
	for x > 0 {
		fmt.Print("\t")
		x--
	}
	fmt.Printf("%v %v\n", node.chars, node.isEnd)
	if len(node.children) != 0 {
		// iterate over children, print each
		for _, child := range node.children {
			child.print(level + 1)
		}
	}
}

func main() {
	root := new(Node)
	
	root.insert("a")
	root.insert("aba")

	root.print(0)
	
	fmt.Printf("contains 'a': %v\n", root.find("a"))
	
	fmt.Printf("contains 'ab': %v\n", root.find("ab"))
	
	fmt.Printf("contains 'aba': %v\n", root.find("aba"))
}
