package radix

import (
	"fmt"
	"testing"
)

func printit(r *Radix, level int) {
	for i:=0; i < level; i++ {
		fmt.Print("\t")
	}
	fmt.Printf("'%v'  value: %v\n", r.key, r.Value)
	for _, child := range r.children {
		printit(child, level + 1)
	}
}

func validate(r *Radix) bool {
	// None, of the childeren must have a prefix incommon with r.key
	for _, child := range r.children {
		s, i := longestCommonPrefix(r.key, child.key)
		println("s", s, "i", i)
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
	r.Insert("test", nil)
	printit(r, 0)


	r.Insert("slow", nil)
	r.Insert("water", nil)
	r.Insert("tester", nil)
	r.Insert("testering", nil)
	r.Insert("rewater", nil)
	r.Insert("waterrat", nil)
	printit(r, 0)
	validate(r)
	t.Fail()
}

//	r.Insert("ab", "1")
//	r.Insert("a", "2")
//	r.Insert("abd", "3")
//	r.Insert("b", 4)
