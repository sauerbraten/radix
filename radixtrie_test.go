package radixtrie

import (
	"fmt"
	"testing"
)

func (r *Radix) printit(level int) {
	x := level
	for x > 0 {
		fmt.Print("\t")
		x--
	}
	fmt.Printf("'%v'  end: %v\n", r.chars, r.Value)
	if len(r.children) != 0 {
		// iterate over children, print each                                                        
		for _, child := range r.children {
			child.Print(level + 1)
		}
	}

}

func TestRadixTrie(t *testing.T) {
	trie := New()
	if trie.value != nil {
		t.Fail()
	}
}

func Example() {
	// create new trie
	trie := New()

	// insert some strings
	trie.Insert("ab", "value 1")
	trie.Insert("a", "value 2")
	trie.Insert("abd", []byte("value 3"))
	trie.Insert("b", 4)

	// print trie structure, the parameter sets the initial level of indentation
	// not safe for tests, since the output can differ, depending of whether the 'a' or 'b' node comes first in the "range children" in Print()
	//trie.Print(0)

	// delete some strings, even strings not contained
	trie.Delete("c")
	trie.Delete("b")
	trie.Delete("ab")

	// print again, notice the changes:
	// 'b' is gone, 'ab' is no longer an end note, means it is no longer contained as a string
	// not safe for tests, since the output can differ, depending of whether the 'a' or 'b' node comes first in the "range children" in Print()
	//trie.Print(0)

	// use Find() to check if a string is contained in the trie
	fmt.Printf("'a' holds: %v\n", trie.Find("a"))
	fmt.Printf("'x' holds: %v\n", trie.Find("x"))
	fmt.Printf("'abd' holds: %v\n", trie.Find("abd"))

	// Output:
	// 'a' holds: value 2
	// 'x' holds: <nil>
	// 'abd' holds: [118 97 108 117 101 32 51]

}
