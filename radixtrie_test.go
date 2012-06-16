package radixtrie

import (
	"testing"
	"fmt"
)

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
	trie.Print(0)

	// delete some strings, even strings not contained
	trie.Delete("c")
	trie.Delete("b")
	trie.Delete("ab")

	// print again, notice the changes:
	// 'b' is gone, 'ab' is no longer an end note, means it is no longer contained as a string
	trie.Print(0)

	// use Find() to check if a string is contained in the trie
	fmt.Printf("'a' holds: %v\n", trie.Find("a"))
	fmt.Printf("'x' holds: %v\n", trie.Find("x"))
	fmt.Printf("'abd' holds: %v\n", trie.Find("abd"))

	// Output:
	// ''  end: <nil>
	// 	'a'  end: value 2
	// 		'b'  end: value 1
	// 			'd'  end: [118 97 108 117 101 32 51]
	// 	'b'  end: 4
	// ''  end: <nil>
	// 	'a'  end: value 2
	// 		'bd'  end: [118 97 108 117 101 32 51]
	// 'a' holds: value 2
	// 'x' holds: <nil>
	// 'abd' holds: [118 97 108 117 101 32 51]

}	
