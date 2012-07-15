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
	if r.Len() != 1 {
		t.Log("Len should be 1")
	}
	r.Insert("test", nil)
	r.Insert("slow", nil)
	r.Insert("water", nil)
	r.Insert("tester", nil)
	r.Insert("testering", nil)
	r.Insert("rewater", nil)
	r.Insert("waterrat", nil)
	if !validate(r) {
		t.Log("Tree does not validate")
		t.Fail()
	}
}

func TestRemove(t *testing.T) {
	r := New()
	r.Insert("test", "aa")
	r.Insert("slow", "bb")

	if r.Remove("slow").(string)!= "bb" {
		t.Log("should be bb")
		t.Fail()
	}

	if r.Remove("slow") != nil {
		t.Log("should be nil")
		t.Fail()
	}
}

// prefix tester
// prefix testering
// prefix testeringandmore
func ExampleTestFind(t *testing.T) {
	r := New()
	r.Insert("tester", nil)
	r.Insert("testering", nil)
	r.Insert("te", nil)
	r.Insert("testeringandmore", nil)
	for _, s := range r.FindPrefix("tester") {
		println("prefix", s)
	}
}

