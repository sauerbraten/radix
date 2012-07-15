// Package radix implements a radix tree.                                                           
//                                                                                                  
// A radix tree is defined in:                                                                      
//    Donald R. Morrison. "PATRICIA -- practical algorithm to retrieve                              
//    information coded in alphanumeric". Journal of the ACM, 15(4):514-534,                        
//    October 1968                                                                                  
//
// Also see http://en.wikipedia.org/wiki/Radix_tree for more information.
package radixtrie

// Radix represents a radix tree.
type Radix struct {
	// children maps the first letter of each child to the child itself, e.g. "a" -> "ab", "x" -> "xyz", "y" -> "yza", ...
	children map[byte]*Radix
	key    string

	// The contents of the radix node.
	Value interface{}
}

func longestCommonPrefix(key, bar string) (string, int) {
	x := 0
	for key[x] == bar[x] {
		x = x + 1
		if x == len(key) || x == len(bar) {
			break
		}
	}
	return key[:x], x // == bar[:x]
}

// Insert inserts the value into the tree with the specified key.
func (r *Radix) Insert(key string, value interface{}) {
	// look up the child starting with the same letter as key
	// if there is no child with the same starting letter, insert a new one
	child, ok := r.children[key[0]]
	if !ok {
		r.children[key[0]] = &Radix{make(map[byte]*Radix), key, value}
		return
	}

	// if key == child.key, don't have to create a new child, but only have to set the (maybe new) value
	if key == child.key {
		child.Value = value
		return
	}

	// commonPrefix is now the longest common substring of key and child.key [e.g. only "ab" from "abab" is contained in "abba"]
	commonPrefix, prefixEnd := longestCommonPrefix(key, child.key)

	// insert keys not shared if commonPrefix == child.key [e.g. child is "ab", key is "abc". we only want to insert "c" below "ab"]
	if commonPrefix == child.key {
		child.Insert(key[prefixEnd:], value)
		return
	}

	// if current child is "abc" and key is "abx", we need to create a new child "ab" with two sub children "c" and "x"

	// create new child node to replace current child
	newChild := &Radix{make(map[byte]*Radix), commonPrefix, nil}

	// replace child of current node with new child: map first letter of common prefix to new child
	r.children[commonPrefix[0]] = newChild

	// shorten old keys to the non-shared part
	child.key = child.key[prefixEnd:]

	// map old child's new first letter to old child as a child of the new child
	newChild.children[child.key[0]] = child

	// if there are keys left of key, insert them into our new child
	if key != newChild.key {
		// insert keys left of key into new child [insert "abba" into "abab" -> "ab" with "ab" as child. now go into node "ab" and create child node "ba"]
		newChild.Insert(key[prefixEnd:], value)

		// else, set new.Child.Value to the value to insert and return
	} else {
		newChild.Value = value
		return
	}
}

// Find returns the value associated with key, or nil if there is no value set to key
func (r *Radix) Find(key string) interface{} {
	// look up the child starting with the same letter as key
	// if there is no child with the same starting letter, return false
	child, ok := r.children[key[0]]
	if !ok {
		return nil
	}

	// check if the end of our string is found and return .isEnd
	if key == child.key {
		return child.Value
	}

	// commonPrefix is now the longest common substring of key and child.key [e.g. only "ab" from "abab" is contained in "abba"]
	commonPrefix, prefixEnd := longestCommonPrefix(key, child.key)

	// if child.key is not completely contained in key, abort [e.g. trying to find "ab" in "abc"]
	if child.key != commonPrefix {
		return nil
	}

	// find the keys left of key in child
	return child.Find(key[prefixEnd:])
}

// Remove unsets any value set to key. If no value exists for key, nothing happens.
// TODO(mg): remove entire node??
func (r *Radix) Remove(key string) {
	// look up the child starting with the same letter as key
	// if there is no child with the same starting letter, return
	child, ok := r.children[key[0]]
	if !ok {
		return
	}

	// if the correct end node is found...
	if key == child.key {
		if len(child.children) == 0 {
			// remove child from current node if child has no children on its own
			delete(r.children, key[0])
		} else if len(child.children) == 1 {
			// since len(child.children) == 1, there is only one subchild; we have to use range to get the value, though, since we do not know the key
			for _, subchild := range child.children {
				// essentially moves the subchild up one level to replace the child we want to delete, while keeping the keys of child
				child.key = child.key + subchild.key
				child.Value = subchild.Value
				child.children = subchild.children
			}
		} else {
			// if there are >= 2 subchilds, we can only set the value to nil, thus delete any value set to key
			child.Value = nil
		}
		return
	}

	// key != child.keys

	// commonPrefix is now the longest common substring of key and child.keys [e.g. only "ab" from "abab" is contained in "abba"]
	commonPrefix, prefixEnd := longestCommonPrefix(key, child.key)

	// if child.keys is not completely contained in key, abort [e.g. trying to delete "ab" from "abc"]
	if child.key != commonPrefix {
		return
	}

	// else: cut off common prefix and delete left string in child
	child.Remove(key[prefixEnd:])
}

// New returns an initialized radix trie.
func New() *Radix {
	return &Radix{make(map[byte]*Radix), "", nil}
}
