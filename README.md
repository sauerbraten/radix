# radixtrie

An implementation of the [radix trie data structure](http://en.wikipedia.org/wiki/Radix_tree) in Go

## Usage

Get the package:

	$ go get github.com/sauerbraten/radixtrie

Import the package:

	import (
		"github.com/sauerbraten/radixtrie"
	)

You can use the radixtrie as a key-value structure, where every node's can have its own value (as shown in the exmaple below), or you can of course just use it to look up strings, like so:

	trie := radixtrie.New()
	trie.Insert("foo", true)
	fmt.Printf("foo is contained: %v\n", trie.Find("foo"))


### Example

This example code is taken from the radixtrie_test.go file

	package main
	
	import (
		"github.com/sauerbraten/radixtrie"
		"fmt"
	)
	
	func main() {
		// create new trie
		trie := radixtrie.New()
		
		// insert some strings
		trie.Insert("abc", "value 1")
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
		fmt.Printf("'c' holds: %v\n", trie.Find("c"))
		fmt.Printf("'abd' holds: %v\n", trie.Find("abd"))
	}

This example should print the following:

	''  end: <nil>
	'a'  end: value 2
		'b'  end: <nil>
			'c'  end: value 1
			'd'  end: [118 97 108 117 101 32 51]
	'b'  end: 4
	''  end: <nil>
		'a'  end: value 2
			'b'  end: <nil>
				'c'  end: value 1
				'd'  end: [118 97 108 117 101 32 51]
	'a' holds: value 2
	'c' holds: <nil>
	'abd' holds: [118 97 108 117 101 32 51]
