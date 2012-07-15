package radix

import (
	"fmt"
	"testing"
)

func printit(r *Radix, level int) {
	for i := 0; i < level; i++ {
		fmt.Print("\t")
	}
	fmt.Printf("'%v'  value: %v\n", r.key, r.Value)
	for _, child := range r.children {
		printit(child, level+1)
	}
}

// None, of the childeren must have a prefix incommon with r.key
func validate(r *Radix) bool {
	for _, child := range r.children {
		_, i := longestCommonPrefix(r.key, child.key)
		if i != 0 {
			return false
		}
		validate(child)
	}
	return true
}

func TestInsert(t *testing.T) {
	r := New()
	if !validate(r) {
		t.Log("Tree does not validate")
		t.Fail()
	}
	println(r.Len())
	r.Insert("test", nil)
	println(r.Len())
	printit(r, 0)

	r.Insert("slow", nil)
	r.Insert("water", nil)
	r.Insert("tester", nil)
	r.Insert("testering", nil)
	r.Insert("rewater", nil)
	r.Insert("waterrat", nil)
	println(r.Len())
	printit(r, 0)
	validate(r)
	t.Fail()
}

//	r.Insert("ab", "1")
//	r.Insert("a", "2")
//	r.Insert("abd", "3")
//	r.Insert("b", 4)
