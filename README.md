# Radix

An implementation of a radix tree in Go. See

> Donald R. Morrison. "PATRICIA -- practical algorithm to retrieve              
> information coded in alphanumeric". Journal of the ACM, 15(4):514-534,        
> October 1968    

Or the [wikipedia article](http://en.wikipedia.org/wiki/Radix_tree).

## Usage

Get the package:

	$ go get github.com/sauerbraten/radix

Import the package:

	import (
		"github.com/sauerbraten/radix"
	)

You can use the tree as a key-value structure, where every node's value can have 
a different type:

	r := radix.New()
	r.Set("one", "1")
	r.Set("twoAndThree", []int{2, 3})
	fmt.Printf("%v, %v, %v\n", r.Get("one"), r.Get("twoAndThree").([]int)[0], r.Get("twoAndThree").([]int)[1])
	
	// prints 1, 2, 3

Or you can of course just use it to look up strings, like so:

	r := radix.New()
	r.Set("foo", true)
	fmt.Printf("foo is contained: %v\n", r.Get("foo"))

### Documentation

For full package documentation, visit http://go.pkgdoc.org/github.com/sauerbraten/radix.

## License

This code is licensed under a BSD License:

    Copyright (c) 2012 Alexander Willing and Miek Gieben. All rights reserved.

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
