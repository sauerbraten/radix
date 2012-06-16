package radixtrie

import (
	"fmt"
)

type RadixTrie struct {
	// children maps the first letter of each child to the child itself, e.g. "a" -> "ab", "x" -> "xyz", "y" -> "yza", ...
	children map[byte]*RadixTrie
	chars    string
	value    interface{}
}

func getLongestCommonPrefix(key, bar string) (string, int) {
	x := 0
	for key[x] == bar[x] {
		x = x + 1
		if x == len(key) || x == len(bar) {
			break
		}
	}
	return key[:x], x // == bar[:x]
}

func (node *RadixTrie) Insert(key string, value interface{}) {
	// look up the child starting with the same letter as key
	// if there is no child with the same starting letter, insert a new one
	child, ok := node.children[key[0]]
	if !ok {
		newChild := &RadixTrie{make(map[byte]*RadixTrie), key, value}
		node.children[key[0]] = newChild
		return
	}

	// if key == child.chars, don't have to create a new child, but only have to set the current child's .isEnd to true
	if key == child.chars {
		child.value = value
		return
	}

	// commonPrefix is now the longest common substring of key and child.chars [e.g. only "ab" from "abab" is contained in "abba"]
	commonPrefix, prefixEnd := getLongestCommonPrefix(key, child.chars)

	// insert chars not shared if commonPrefix == child.chars [e.g. child is "ab", key is "abc". we only want to insert "c" below "ab"]
	if commonPrefix == child.chars {
		child.Insert(key[prefixEnd:], value)
		return
	}

	// if current child is "abc" and key is "abx", we need to create a new child "ab" with two sub children "c" and "x"

	// create new child node to replace current child
	newChild := &RadixTrie{make(map[byte]*RadixTrie), commonPrefix, nil}

	// replace child of current node with new child: map first letter of common prefix to new child
	node.children[commonPrefix[0]] = newChild

	// shorten old chars to the non-shared part
	child.chars = child.chars[prefixEnd:]

	// map old child's new first letter to old child as a child of the new child
	newChild.children[child.chars[0]] = child

	// if there are chars left of key, insert them into our new child
	if key != newChild.chars {
		// insert chars left of key into new child [insert "abba" into "abab" -> "ab" with "ab" as child. now go into node "ab" and create child node "ba"]
		newChild.Insert(key[prefixEnd:], value)

	// else, set new.Child.isEnd = true and return
	} else {
		newChild.value = value
		return
	}
}

// return true if key is contained in the tree
func (node *RadixTrie) Find(key string) interface{} {
	// look up the child starting with the same letter as key
	// if there is no child with the same starting letter, return false
	child, ok := node.children[key[0]]
	if !ok {
		return nil
	}

	// check if the end of our string is found and return .isEnd
	if key == child.chars {
		return child.value
	}

	// commonPrefix is now the longest common substring of key and child.chars [e.g. only "ab" from "abab" is contained in "abba"]
	commonPrefix, prefixEnd := getLongestCommonPrefix(key, child.chars)

	// if child.chars is not completely contained in key, abort [e.g. trying to find "ab" in "abc"]
	if child.chars != commonPrefix {
		return nil
	}

	// find the chars left of key in child
	return child.Find(key[prefixEnd:])
}

func (node *RadixTrie) Delete(key string) {
	// look up the child starting with the same letter as key
	// if there is no child with the same starting letter, return
	child, ok := node.children[key[0]]
	if !ok {
		return
	}

	// set child.isEnd = false if key == child.chars and, if child has no children on its own, remove it completely
	if key == child.chars {
		// remove child from current node if empty (when child has no children on its own)	
		if len(child.children) == 0 {
			delete(node.children, key[0])
		} else if len(child.children) == 1 {
			// iterate over map to get the single key-value pair
			for _, subchild := range child.children {
				child.chars = child.chars + subchild.chars
				child.value = subchild.value
				child.children = subchild.children
			}
		} else {
			child.value = nil
		}
		return
	}

	// key != child.chars

	// commonPrefix is now the longest common substring of key and child.chars [e.g. only "ab" from "abab" is contained in "abba"]
	commonPrefix, prefixEnd := getLongestCommonPrefix(key, child.chars)

	// if child.chars is not completely contained in key, abort [e.g. trying to delete "ab" from "abc"]
	if child.chars != commonPrefix {
		return
	}

	// else: cut off common prefix and delete left string in child
	child.Delete(key[prefixEnd:])
}

func (node *RadixTrie) Print(level int) {
	x := level
	for x > 0 {
		fmt.Print("\t")
		x--
	}
	fmt.Printf("'%v'  end: %v\n", node.chars, node.value)
	if len(node.children) != 0 {
		// iterate over children, print each
		for _, child := range node.children {
			child.Print(level + 1)
		}
	}
}

func New() *RadixTrie {
	return &RadixTrie{make(map[byte]*RadixTrie), "", nil}
}
