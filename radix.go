// Package radix implements a radix tree.                                                           
//                                                                                                  
// A radix tree is defined in:                                                                      
//    Donald R. Morrison. "PATRICIA -- practical algorithm to retrieve                              
//    information coded in alphanumeric". Journal of the ACM, 15(4):514-534,                        
//    October 1968                                                                                  
//
// Also see http://en.wikipedia.org/wiki/Radix_tree for more information.
//
// Basic use pattern for iterating over a radix tree and retrieving the full
// keys under which nodes are stored:
//
//	func iter(r *Radix, prefix string) {
//		// current key r.Key()
//		// full key = prefix + r.Key()
//		for _, child := range r.Children() {
//			iter(child, prefix + r.Key()
//		}
//      }
//
//	f := r.Find("tester")		// Look for "tester"
//	iter(f, f.Prefix("tester"))	// Get all the keys from "tester" down
package radix

import (
	"strings"
)

// Radix represents a radix tree.
type Radix struct {
	// children maps the first letter of each child to the child itself, e.g. "a" -> "ab", "x" -> "xyz", "y" -> "yza", ...
	children map[byte]*Radix
	key      string

	// The contents of the radix node.
	Value interface{}
}

// Key returns the (partial) key under which this node is stored.
func (r *Radix) Key() string {
	if r != nil {
		return r.key
	}
	return ""
}

// Children returns the children of this node r or nil if there are none.
func (r *Radix) Children() map[byte]*Radix {
	if r != nil {
		return r.children
	}
	return nil
}

func longestCommonPrefix(key, bar string) (string, int) {
	if key == "" || bar == "" {
		return "", 0
	}
	x := 0
	for key[x] == bar[x] {
		x = x + 1
		if x == len(key) || x == len(bar) {
			break
		}
	}
	return key[:x], x // == bar[:x]
}

// Insert inserts the value into the tree with the specified key. It returns the radix node
// it just inserted
func (r *Radix) Insert(key string, value interface{}) *Radix {
	// look up the child starting with the same letter as key
	// if there is no child with the same starting letter, insert a new one
	child, ok := r.children[key[0]]
	if !ok {
		r.children[key[0]] = &Radix{make(map[byte]*Radix), key, value}
		return r.children[key[0]]
	}

	// if key == child.key, don't have to create a new child, but only have to set the (maybe new) value
	if key == child.key {
		child.Value = value
		return child
	}

	// commonPrefix is now the longest common substring of key and child.key [e.g. only "ab" from "abab" is contained in "abba"]
	commonPrefix, prefixEnd := longestCommonPrefix(key, child.key)

	// insert key not shared if commonPrefix == child.key [e.g. child is "ab", key is "abc". we only want to insert "c" below "ab"]
	if commonPrefix == child.key {
		return child.Insert(key[prefixEnd:], value)
	}

	// if current child is "abc" and key is "abx", we need to create a new child "ab" with two sub children "c" and "x"

	// create new child node to replace current child
	newChild := &Radix{make(map[byte]*Radix), commonPrefix, nil}

	// replace child of current node with new child: map first letter of common prefix to new child
	r.children[commonPrefix[0]] = newChild

	// shorten old key to the non-shared part
	child.key = child.key[prefixEnd:]

	// map old child's new first letter to old child as a child of the new child
	newChild.children[child.key[0]] = child

	// if there are key left of key, insert them into our new child
	if key != newChild.key {
		// insert key left of key into new child [insert "abba" into "abab" -> "ab" with "ab" as child. now go into node "ab" and create child node "ba"]
		newChild.Insert(key[prefixEnd:], value)
		// else, set new.Child.Value to the value to insert and return
	} else {
		newChild.Value = value
	}
	return newChild
}

// Find returns the node associated with key. All childeren of this node share the same prefix
func (r *Radix) Find(key string) *Radix {
	r1 := r.find(key)
	if r1 == nil {
		return nil
	}
	return r1
}

// find returns the radix (sub)tree associated with key, or nil if there is nothing found
func (r *Radix) find(key string) *Radix {
	// look up the child starting with the same letter as key
	// if there is no child with the same starting letter, return false
	child, ok := r.children[key[0]]
	if !ok {
		return nil
	}

	// check if the end of our string is found and return .isEnd
	if key == child.key {
		return child
	}

	// commonPrefix is now the longest common substring of key and child.key [e.g. only "ab" from "abab" is contained in "abba"]
	commonPrefix, prefixEnd := longestCommonPrefix(key, child.key)

	// if child.key is not completely contained in key, abort [e.g. trying to find "ab" in "abc"]
	if child.key != commonPrefix {
		return nil
	}

	// find the key left of key in child
	return child.find(key[prefixEnd:])
}

// Remove removes any value set to key. It returns the removed node or nil if the
// node cannot be found.
func (r *Radix) Remove(key string) *Radix {
	// look up the child starting with the same letter as key
	// if there is no child with the same starting letter, return
	child, ok := r.children[key[0]]
	if !ok {
		return nil
	}

	// if the correct end node is found...
	if key == child.key {
		switch len(child.children) {
		case 0:
			// remove child from current node if child has no children on its own
			delete(r.children, key[0])
		case 1:
			// since len(child.children) == 1, there is only one subchild; we have to use range to get the value, though, since we do not know the key
			for _, subchild := range child.children {
				// essentially moves the subchild up one level to replace the child we want to delete, while keeping the key of child
				child.key = child.key + subchild.key
				child.Value = subchild.Value
				child.children = subchild.children
			}
		default:
			// if there are >= 2 subchilds, we can only set the value to nil, thus delete any value set to key
			child.Value = nil
		}
		return child
	}

	// key != child.keys

	// commonPrefix is now the longest common substring of key and child.key [e.g. only "ab" from "abab" is contained in "abba"]
	commonPrefix, prefixEnd := longestCommonPrefix(key, child.key)

	// if child.key is not completely contained in key, abort [e.g. trying to delete "ab" from "abc"]
	if child.key != commonPrefix {
		return nil
	}

	// else: cut off common prefix and delete left string in child
	return child.Remove(key[prefixEnd:])
}

// Do calls function f on each node in the tree. The behavior of Do is              
// undefined if f changes *r.                                                       
func (r *Radix) Do(f func(interface{})) {
	if r != nil {
		f(r.Value)
		for _, child := range r.children {
			child.Do(f)
		}
	}
}

// Prefix returns the string that is the result of "subtracting" r.Key from s. 
// If s equals "tester" and the key were r is stored is "ster", Prefix returns
// "te".
func (r *Radix) Prefix(s string) string {
	l := strings.LastIndex(s, r.Key())
	if l == -1 {
		return ""
	}
	return s[:l]
}

// Len computes the number of nodes in the radix tree r.
func (r *Radix) Len() int {
	i := 0
	if r != nil {
		if r.Value != nil {
			i++
		}
		for _, child := range r.children {
			i += child.Len()
		}
	}
	return i
}

// New returns an initialized radix tree.
func New() *Radix {
	return &Radix{make(map[byte]*Radix), "", nil}
}
