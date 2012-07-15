package radixtrie

import (
	"fmt"
	"testing"
)

func (r *Radix) printit(level int) {
	for i:=0; i < level; i++ {
		fmt.Print("\t")
	}
	fmt.Printf("'%v'  value: %v\n", r.chars, r.Value)
	if len(r.children) != 0 {
		// iterate over children, print each                                                        
		for _, child := range r.children {
			child.printit(level + 1)
		}
	}
}

func TestInsert(t *testing.T) {
	r := New()
	r.Insert("test", nil)
	r.Insert("slow", nil)
	r.Insert("water", nil)
}

func Example() {
	// create new trie
	trie := New()

	// insert some strings
	trie.Insert("ab", "1")
	trie.Insert("a", "2")
	trie.Insert("abd", "3"))
	trie.Insert("b", 4)

	trie.printit(0)

	trie.Remove("c")
	trie.Remove("b")
	trie.Remove("ab")
	trie.printit(0)

	// print again, notice the changes:
	// 'b' is gone, 'ab' is no longer an end note, means it is no longer contained as a string
	// not safe for tests, since the output can differ, depending of whether the 'a' or 'b' node comes first in the "range children" in Print()
	trie.printit(0)

	// use Find() to check if a string is contained in the trie
	fmt.Printf("'a' holds: %v\n", trie.Find("a"))
	fmt.Printf("'x' holds: %v\n", trie.Find("x"))
	fmt.Printf("'abd' holds: %v\n", trie.Find("abd"))

	// Output:
	// 'a' holds: value 2
	// 'x' holds: <nil>
	// 'abd' holds: [118 97 108 117 101 32 51]

}
