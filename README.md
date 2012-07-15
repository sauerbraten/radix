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

### Documentation

For full package documentation, visit http://go.pkgdoc.org/github.com/sauerbraten/radixtrie.

## License

This code is licensed under the MIT License:

	Copyright (c) 2012 Alexander Willing. All rights reserved.
	
	Redistribution and use in source and binary forms, with or without
	modification, are permitted provided that the following conditions are
	met:
	
		* Redistributions of source code must retain the above copyright
	notice, this list of conditions and the following disclaimer.
		* Redistributions in binary form must reproduce the above
	copyright notice, this list of conditions and the following disclaimer
	in the documentation and/or other materials provided with the
	distribution.
	
	THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
	"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
	LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
	A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
	OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
	SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
	LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
	DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
	THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
	(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
	OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.