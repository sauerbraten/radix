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
//		// Current key is r.Key()
//		// The full key would be prefix + r.Key()
//		for _, child := range r.Children() {
//			iter(child, prefix + r.Key()
//		}
//	}
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
	value interface{}
}

// Key returns the (partial) key under which r is stored.
func (r *Radix) Key() string {
	if r != nil {
		return r.key
	}
	return ""
}

// Children returns the children of r or nil if there are none.
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
// it just inserted.
func (r *Radix) Insert(key string, value interface{}) *Radix {
	// look up the child starting with the same letter as key
	// if there is no child with the same starting letter, insert a new one

	newR, ok := r.children[key[0]]

	for ok {
		key = key[len(r.key):]
		r = newR
		if len(key) < len(r.key) {
			break
		} else {
			newR, ok = r.children[key[len(r.key)]]
		}
	}
	
	// r is the deepest we can go now

	if key == r.key {
		r.value = value
		return r
	}

	// commonPrefix is now the longest common substring of key and child.key [e.g. only "ab" from "abab" is contained in "abba"]
	cP, prefixEnd := longestCommonPrefix(key, r.key)

	if len(cP) < len(key) {
		if len(cP) < len(r.key) {
			// newOldR := &Radix{r.children, cP, r.value}
			newOldR := &Radix{r.children, cP, r.value}
			// newR: empty children, uncommmon part, value
			newR := &Radix{make(map[byte]*Radix), r.key[prefixEnd:], value}
			// r: [newOldR, newR], cP, nil
			r.children = make(map[byte]*Radix)
			r.children[newOldR.key[0]] = newOldR
			r.children[newR.key[0]] = newR
			r.key = cP
			r.value = nil
			return newR
		} else { //len(cP) == len(r.key)
			// newR: empty children, uncommmon part, value
			newR := &Radix{make(map[byte]*Radix), key[prefixEnd:], value}
			// r: r.children + newR
			r.children[newR.key[0]] = newR
			return newR
		}
	} else { // len(cP) == len(key)
		// newOldR := &Radix{r.children, uncommmon part, r.value}
		newOldR := &Radix{r.children, r.key[prefixEnd:], r.value}
		// r: [newOldR], key, value
		r.children = make(map[byte]*Radix)
		r.children[newOldR.key[0]] = newOldR
		r.key = key
		r.value = value
		return r
	}

	// should never be reached
	return nil
}

// Find returns the node associated with key. All childeren of this node share the same prefix
func (r *Radix) Find(key string) *Radix {
	// look up the child starting with the same letter as key
	// if there is no child with the same starting letter, return false
	child, ok := r.children[key[0]]
	if !ok {
		return nil
	}

	posInKey := 0

	for child.key != key[posInKey:] {
		// commonPrefix is now the longest common substring of key and child.key [e.g. only "ab" from "abab" is contained in "abba"]
		commonPrefix, prefixEnd := longestCommonPrefix(key[posInKey:], child.key)
		posInKey = posInKey + prefixEnd

		// if child.key is not completely contained in key, abort [e.g. trying to find "ab" in "abc"]
		if child.key != commonPrefix {
			return nil
		}

		// if there is no child starting with the leftover key, abort
		child, ok = child.children[key[posInKey]]
		if !ok {
			return nil
		}
	}

	return child
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
				child.value = subchild.value
				child.children = subchild.children
			}
		default:
			// if there are >= 2 subchilds, we can only set the value to nil, thus delete any value set to key
			child.value = nil
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

// Do calls function f on each node in the tree. f's parameter will be r.Value. The behavior of Do is              
// undefined if f changes r.                                                       
func (r *Radix) Do(f func(interface{})) {
	if r != nil {
		f(r.value)
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
		if r.value != nil {
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
