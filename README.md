# radixtrie

An implementation of the [radix trie data structure](http://en.wikipedia.org/wiki/Radix_tree) in Go

## Usage

Get the package:

	$ go get github.com/sauerbraten/radixtrie

Import the package:

	import (
		"github.com/sauerbraten/radixtrie"
	)

### Example
	package main
	
	import (
		"github.com/sauerbraten/radixtrie"
		"fmt"
	)
	
	func main() {
		// create new trie
		trie := radixtrie.New()
		
		// insert some strings
		trie.Insert("abc")
		trie.Insert("a")
		trie.Insert("abd")
		trie.Insert("b")
		
		// print trie structure, the parameter sets the initial level of indentation
		fmt.Println("Trie after inserting and before deleting")
		trie.Print(0)
		
		// delete some strings, even strings not contained
		trie.Delete("c")
		trie.Delete("b")
		trie.Delete("ab")
		
		// print again, notice the changes:
		// 'b' is gone, 'ab' is no longer an end note, means it is no longer contained as a string
		fmt.Println("Trie after deleting")
		trie.Print(0)
		
		// use Find() to check if a string is contained in the trie
		fmt.Printf("'a' is contained: %v\n", trie.Find("a"))
		fmt.Printf("'c' is contained: %v\n", trie.Find("c"))
		fmt.Printf("'abd' is contained: %v\n", trie.Find("abd"))
	}
