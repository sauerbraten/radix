package radixtrie

import (
	"fmt"
)

type RadixTrie struct {
	// children maps the first letter of each child to the child itself, e.g. "a" -> "ab", "x" -> "xyz", "y" -> "yza", ...
	children	map[byte]*RadixTrie
	chars		string
	isEnd		bool
}

func getLongestCommonPrefix(foo, bar string) (string, int) {
	x := 0
	for foo[x] == bar[x] {
		x = x + 1
		if x == len(foo) || x == len(bar) {
			break
		}
	}
	return foo[:x], x // == bar[:x]
}

func (node *RadixTrie) Insert(foo string) {
	// look up the child starting with the same letter as foo
	// if there is no child with the same starting letter, insert a new one
	child, ok := node.children[foo[0]]	
	if !ok {
		newChild := &RadixTrie{make(map[byte]*RadixTrie), foo, true}
		node.children[foo[0]] = newChild
		return
	}
	
	// if foo == child.chars, don't have to create a new child, but only have to set the current child's .isEnd to true
	if foo == child.chars {
		child.isEnd = true
		return
	}

	// commonPrefix is now the longest common substring of foo and child.chars [e.g. only "ab" from "abab" is contained in "abba"]
	commonPrefix, prefixEnd := getLongestCommonPrefix(foo, child.chars)

	// insert chars not shared if commonPrefix == child.chars [e.g. child is "ab", foo is "abc". we only want to insert "c" below "ab"]
	if commonPrefix == child.chars {
		child.Insert(foo[prefixEnd:])
		return
	}

	// if current child is "abc" and foo is "abx", we need to create a new child "ab" with two sub children "c" and "x"
	
	// create new child node to replace current child
	newChild := &RadixTrie{make(map[byte]*RadixTrie), commonPrefix, false}
	
	// replace child of current node with new child: map first letter of common prefix to new child
	node.children[commonPrefix[0]] = newChild
	
	// shorten old chars to the non-shared part
	child.chars = child.chars[prefixEnd:]

	// map old child's new first letter to old child as a child of the new child
	newChild.children[child.chars[0]] = child
	
	// insert chars left of foo into new child [insert "abba" into "abab" -> "ab" with "ab" as child. now go into node "ab" and create child node "ba"]
	newChild.Insert(foo[prefixEnd:])
}


// return true if foo is contained in the tree
func (node *RadixTrie) Find(foo string) bool {
	// look up the child starting with the same letter as foo
	// if there is no child with the same starting letter, return false
	child, ok := node.children[foo[0]]	
	if !ok {
		return false
	}

	// check if the end of our string is found and return .isEnd
	if foo == child.chars {
		return child.isEnd
	}
	
	// commonPrefix is now the longest common substring of foo and child.chars [e.g. only "ab" from "abab" is contained in "abba"]
	commonPrefix, prefixEnd := getLongestCommonPrefix(foo, child.chars)
	
	// if child.chars is not completely contained in foo, abort [e.g. trying to find "ab" in "abc"]
	if  child.chars != commonPrefix {
		return false
	}

	// find the chars left of foo in child
	return child.Find(foo[prefixEnd:])
}

func (node *RadixTrie) Delete(foo string) {
	// look up the child starting with the same letter as foo
	// if there is no child with the same starting letter, return
	child, ok := node.children[foo[0]]	
	if !ok {
		return
	}
	
	// set child.isEnd = false if foo == child.chars and, if child has no children on its own, remove it completely
	if foo == child.chars {
		child.isEnd = false
		// remove child from current node if empty (when child has no children on its own)	
		if len(child.children) == 0 {
			delete(node.children, foo[0])
		}
		return
	}

	// foo != child.chars

	// commonPrefix is now the longest common substring of foo and child.chars [e.g. only "ab" from "abab" is contained in "abba"]
	commonPrefix, prefixEnd := getLongestCommonPrefix(foo, child.chars)

	// if child.chars is not completely contained in foo, abort [e.g. trying to delete "ab" from "abc"]
	if  child.chars != commonPrefix {
		return
	}
	
	// else: cut off common prefix and delete left string in child
	child.Delete(foo[prefixEnd:])
}

func (node *RadixTrie) Print(level int) {
	x := level
	for x > 0 {
		fmt.Print("\t")
		x--
	}
	fmt.Printf("'%v'  end: %v\n", node.chars, node.isEnd)
	if len(node.children) != 0 {
		// iterate over children, print each
		for _, child := range node.children {
			child.Print(level + 1)
		}
	}
}

func New() *RadixTrie {
	return &RadixTrie{make(map[byte]*RadixTrie), "", false}
}
