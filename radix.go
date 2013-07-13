// Package radix implements a radix tree.
//
// A radix tree is defined in:
//    Donald R. Morrison. "PATRICIA -- practical algorithm to retrieve
//    information coded in alphanumeric". Journal of the ACM, 15(4):514-534,
//    October 1968
//
// Also see http://en.wikipedia.org/wiki/Radix_tree for more information.
package radix

// Radix represents a radix tree.
type Radix struct {
	parent *Radix
	// children maps the first letter of each child to the child itself, e.g. "a" -> "ab", "x" -> "xyz", "y" -> "yza", ...
	children map[byte]*Radix
	key      string

	// The contents of the radix node.
	value interface{}
}

// Value returns the value stored udner the key ending at the node r.
func (r *Radix) Value() interface{} {
	return r.value
}

// Children returns the children of r or nil if there are none.
func (r *Radix) Children() map[byte]*Radix {
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
	return &Radix{nil, make(map[byte]*Radix), "", nil}
}

// helper function, used in Set() and SubTree()
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

// Set inserts the value into the tree with the specified key. It returns the radix node
// it just inserted.
func (r *Radix) Set(key string, value interface{}) *Radix {
	// look up the child starting with the same letter as key
	// if there is no child with the same starting letter, insert a new one

	newR, ok := r.children[key[0]]

	for ok {
		key = key[len(r.key):]
		r = newR
		if len(key) <= len(r.key) {
			break
		} else {
			newR, ok = r.children[key[len(r.key)]]
		}
	}

	// r now is the deepest we can go

	if key == r.key {
		r.value = value
		return r
	}

	// commonPrefix is now the longest common substring of key and child.key [e.g. only "ab" from "abab" is contained in "abba"]
	cP, prefixEnd := longestCommonPrefix(key, r.key)

	// key: 'abcd', cP: 'abc'
	if len(cP) < len(key) {
		// key: 'abcd', cP: 'abc', r.key: 'abcx'
		if len(cP) < len(r.key) {
			// newOldR := &Radix{r, r.children, cP, r.Value}
			newOldR := &Radix{r, r.children, r.key[prefixEnd:], r.value}
			// newR: r, empty children, uncommmon part, value
			newR := &Radix{r, make(map[byte]*Radix), key[prefixEnd:], value}
			// r: r.parent, [newOldR, newR], cP, nil
			r.children = make(map[byte]*Radix)
			r.children[newOldR.key[0]] = newOldR
			r.children[newR.key[0]] = newR
			r.key = cP
			r.value = nil
			// go into newly created r for return statement at the end
			r = newR
			// key: 'abcd', cP: 'abc', r.key: 'abc'
		} else { //len(cP) == len(r.key)
			// newR: r, empty children, uncommmon part, value
			newR := &Radix{r, make(map[byte]*Radix), key[prefixEnd:], value}
			// r: r.parent, r.children + newR, r.key, r.Value
			r.children[newR.key[0]] = newR
			// go into newly created r for return statement at the end
			r = newR
		}
		// key: 'abc', cP: 'abc', r.key: 'abcd'
	} else { // len(cP) == len(key)
		// newOldR := &Radix{r, r.children, uncommmon part, r.Value}
		newOldR := &Radix{r, r.children, r.key[prefixEnd:], r.value}
		// r: r.parent, [newOldR], key, value
		r.children = make(map[byte]*Radix)
		r.children[newOldR.key[0]] = newOldR
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
	r, ok := r.children[key[0]]
	if !ok {
		return nil
	}

	posInKey := 0

	for r.key != key[posInKey:] {
		// commonPrefix is now the longest common substring of key and child.key [e.g. only "ab" from "abab" is contained in "abba"]
		commonPrefix, prefixEnd := longestCommonPrefix(key[posInKey:], r.key)
		posInKey = posInKey + prefixEnd

		// if child.key is not completely contained in key, abort [e.g. trying to find "ab" in "abc"]
		if r.key != commonPrefix {
			return nil
		}

		// if there is no child starting with the leftover key, abort
		r, ok = r.children[key[posInKey]]
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
	r, ok := r.children[prefix[0]]
	if !ok {
		return nil
	}

	posInPrefix := 0

	for posInPrefix != len(prefix) {
		// commonPrefix is now the longest common substring of key and child.key [e.g. only "ab" from "abab" is contained in "abba"]
		commonPrefix, prefixEnd := longestCommonPrefix(prefix[posInPrefix:], r.key)
		posInPrefix += prefixEnd

		if posInPrefix > len(prefix)-1 {
			// if prefix is entirely contained in r.key, return r
			if len(r.key) >= len(commonPrefix) {
				return r
			}

			return nil
		}

		// if there is no child starting with the leftover key, abort
		r, ok = r.children[prefix[posInPrefix]]
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
		delete(r.parent.children, r.key[0])
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
