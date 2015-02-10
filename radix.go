// Package radix implements a radix tree. It uses UTF-8 strings as keys.
//
// A radix tree is defined in:
//    Donald R. Morrison. "PATRICIA -- practical algorithm to retrieve
//    information coded in alphanumeric". Journal of the ACM, 15(4):514-534,
//    October 1968
//
// Also see http://en.wikipedia.org/wiki/Radix_tree for more information.
package radix

import "unicode/utf8"

// Radix represents a radix tree.
type Radix struct {
	parent *Radix
	// children maps the first rune of the key of each child to the child itself, e.g. "a" -> "ab", "x" -> "xyz", "y" -> "yza", ...
	children map[rune]*Radix
	key      string

	// The contents of the radix node.
	value interface{}
}

// Value returns the value stored udner the key ending at the node r.
func (r *Radix) Value() interface{} {
	return r.value
}

// Children returns the children of r or nil if there are none.
func (r *Radix) Children() map[rune]*Radix {
	if r != nil {
		return r.children
	}
	return nil
}

// Key returns the (partial) key under which r is stored.
func (r *Radix) Key() string {
	if r != nil {
		return r.key
	}
	return ""
}

// New returns an initialized radix tree.
func New() *Radix {
	return &Radix{
		parent:   nil,
		children: make(map[rune]*Radix),
		key:      "",
		value:    nil,
	}
}

// helper function, used in Set() and SubTree()
// Returns the longest prefix the two strings have in common by comparing
// runes, as well as the length of that prefix in bytes. Assumes both strings
// have been checked for being valid UTF-8.
func longestCommonPrefix(a, b string) (prefix string, pl int) {
	if a == "" || b == "" {
		return
	}

	for pl < len(a) && pl < len(b) {
		runeA, runeLength := utf8.DecodeRuneInString(a[pl:])
		runeB, _ := utf8.DecodeRuneInString(b[pl:])
		if runeA == runeB {
			pl += runeLength
		} else {
			break
		}
	}

	prefix = a[:pl] // a[:pl] == b[:pl]

	return
}

// Set inserts the value into the tree with the specified key. It returns the
// radix node it just inserted. It returns nil if the key is not valid UTF-8.
func (r *Radix) Set(key string, value interface{}) *Radix {
	if !utf8.ValidString(key) {
		return nil
	}
	// look up the child starting with the same letter as key
	// if there is no child with the same starting letter, insert a new one

	newR, ok := r, true
	if len(key) > 0 {
		firstRune, _ := utf8.DecodeRuneInString(key)
		newR, ok = r.children[firstRune]
	}

	for ok {
		key = key[len(r.key):]
		r = newR
		if len(key) <= len(r.key) {
			break
		} else {
			firstRune, _ := utf8.DecodeRuneInString(key[len(r.key):])
			newR, ok = r.children[firstRune]
		}
	}

	// r now is the deepest we can go

	if key == r.key {
		r.value = value
		return r
	}

	// commonPrefix is now the longest common substring of key and child.key [e.g. only "ab" from "abab" is contained in "abba"]
	commonPrefix, prefixLength := longestCommonPrefix(key, r.key)

	// key: 'abcd', commonPrefix: 'abc'
	if len(commonPrefix) < len(key) {
		// key: 'abcd', commonPrefix: 'abc', r.key: 'abcx'
		if len(commonPrefix) < len(r.key) {
			// newOldR := &Radix{r, r.children, commonPrefix, r.Value}
			newOldR := &Radix{
				parent:   r,
				children: r.children,
				key:      r.key[prefixLength:],
				value:    r.value,
			}

			// newR: r, empty children, uncommmon part, value
			newR := &Radix{
				parent:   r,
				children: make(map[rune]*Radix),
				key:      key[prefixLength:],
				value:    value}

			// r: r.parent, [newOldR, newR], commonPrefix, nil
			newOldRFirstRune, _ := utf8.DecodeRuneInString(newOldR.key)
			newRFirstRune, _ := utf8.DecodeRuneInString(newR.key)
			r.children = map[rune]*Radix{
				newOldRFirstRune: newOldR,
				newRFirstRune:    newR,
			}
			r.key = commonPrefix
			r.value = nil

			// go into newly created r for return statement at the end
			r = newR

		} else { // len(commonPrefix) == len(r.key) → key: 'abcd', commonPrefix: 'abc', r.key: 'abc'
			// newR: r, empty children, uncommmon part, value
			newR := &Radix{
				parent:   r,
				children: map[rune]*Radix{},
				key:      key[prefixLength:],
				value:    value,
			}

			// r: r.parent, r.children + newR, r.key, r.Value
			newRFirstRune, _ := utf8.DecodeRuneInString(newR.key)
			r.children[newRFirstRune] = newR

			// go into newly created r for return statement at the end
			r = newR
		}
		// key: 'abc', commonPrefix: 'abc', r.key: 'abcd'
	} else { // len(commonPrefix) == len(key)
		// newOldR := &Radix{r, r.children, uncommmon part, r.Value}
		newOldR := &Radix{r, r.children, r.key[prefixLength:], r.value}
		// r: r.parent, [newOldR], key, value
		newOldRFirstRune, _ := utf8.DecodeRuneInString(newOldR.key)
		r.children = map[rune]*Radix{
			newOldRFirstRune: newOldR,
		}
		r.key = key
		r.value = value
	}

	return r
}

// SubTree returns the node wich key points to or nil if there is no such key.
func (r *Radix) SubTree(key string) *Radix {
	if len(key) < 1 {
		return nil
	}

	// look up the child starting with the same letter as key
	// if there is no child with the same starting letter, return false
	firstRune, _ := utf8.DecodeRuneInString(key)
	r, ok := r.children[firstRune]
	if !ok {
		return nil
	}

	posInKey := 0

	for r.key != key[posInKey:] {
		// commonPrefix is now the longest common substring of key and child.key [e.g. only "ab" from "abab" is contained in "abba"]
		commonPrefix, prefixLength := longestCommonPrefix(key[posInKey:], r.key)
		posInKey = posInKey + prefixLength

		// if child.key is not completely contained in key, abort [e.g. trying to find "ab" in "abc"]
		if r.key != commonPrefix {
			return nil
		}

		// if there is no child starting with the leftover key, abort
		firstRune, _ := utf8.DecodeRuneInString(key[posInKey:])
		r, ok = r.children[firstRune]
		if !ok {
			return nil
		}
	}

	return r
}

// SubTreeWithPefix returns the node wich key starts with prefix or nil if there is no such node.
func (r *Radix) SubTreeWithPrefix(prefix string) *Radix {
	if len(prefix) < 1 {
		return nil
	}

	// look up the child starting with the same letter as key
	// if there is no child with the same starting letter, return false
	firstRune, _ := utf8.DecodeRuneInString(prefix)
	r, ok := r.children[firstRune]
	if !ok {
		return nil
	}

	posInPrefix := 0

	for posInPrefix != len(prefix) {
		// commonPrefix is now the longest common substring of key and child.key [e.g. only "ab" from "abab" is contained in "abcd"]
		commonPrefix, prefixLength := longestCommonPrefix(prefix[posInPrefix:], r.key)
		posInPrefix += prefixLength

		if posInPrefix > len(prefix)-1 {
			// if prefix is entirely contained in r.key, return r
			if len(r.key) >= len(commonPrefix) {
				return r
			}

			return nil
		}

		// if there is no child starting with the leftover key, abort
		firstRune, _ := utf8.DecodeRuneInString(prefix[posInPrefix:])
		r, ok = r.children[firstRune]
		if !ok {
			return nil
		}
	}

	return r
}

// Get returns the value associated with key or nil if there is no such key.
func (r *Radix) Get(key string) interface{} {
	r = r.SubTree(key)
	if r != nil {
		return r.value
	}

	return nil
}

func (r *Radix) getChildrenValues() (values []interface{}) {
	for _, c := range r.Children() {
		values = append(values, c.getChildrenValues()...)
	}

	if r.value != nil {
		values = append(values, r.value)
	}
	return
}

// Returns all values associated with keys that start with prefix as a slice. The slice is empty if there is no key with that prefix.
func (r *Radix) GetAllWithPrefix(prefix string) (values []interface{}) {
	r = r.SubTreeWithPrefix(prefix)
	if r == nil {
		return
	}

	values = r.getChildrenValues()

	return
}

// Remove removes any value set to key. It returns the removed node or nil if the
// node cannot be found.
func (r *Radix) Remove(key string) *Radix {

	r = r.SubTree(key)
	if r == nil {
		return nil
	}

	switch len(r.children) {
	case 0:
		// remove child from current node if child has no children on its own
		firstRune, _ := utf8.DecodeRuneInString(r.key)
		delete(r.parent.children, firstRune)
	case 1:
		// since len(r.parent.children) == 1, there is only one subchild; we have to use range to get the value, though, since we do not know the key
		for _, child := range r.children {
			// essentially moves the subchild up one level to replace the child we want to delete, while keeping the key of child
			r.key = r.key + child.key
			r.value = child.value
			r.children = child.children
		}
	default:
		// if there are >= 2 subchilds, we can only set the value to nil, thus delete any value set to key
		r.value = nil
	}

	return r
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
