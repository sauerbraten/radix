package radix

import (
	"fmt"
	"testing"
)

func printit(r *Radix, level int) {
	for i := 0; i < level; i++ {
		fmt.Print("\t")
	}
	fmt.Printf("'%v'  value: %v\n", r.key, r.value)
	println()
	for _, child := range r.children {
		printit(child, level+1)
	}
}

// None of the children must have a prefix in common with r.key
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
	if r.Len() != 0 {
		t.Log("Len should be 0", r.Len())
	}
	r.Set("test", nil)
	r.Set("slow", nil)
	r.Set("water", nil)
	r.Set("tester", nil)
	r.Set("testering", nil)
	r.Set("rewater", nil)
	r.Set("waterrat", nil)
	if !validate(r) {
		t.Log("Tree does not validate")
		t.Fail()
	}
}

func TestRemove(t *testing.T) {
	r := New()
	r.Set("test", "aa")
	r.Set("slow", "bb")

	if k := r.Remove("slow").value; k != "bb" {
		t.Log("should be bb", k)
		t.Fail()
	}

	if r.Remove("slow") != nil {
		t.Log("should be nil")
		t.Fail()
	}
}

func ExampleSubTree() {
	r := New()
	r.Set("tester", nil)
	r.Set("testering", nil)
	r.Set("te", nil)
	r.Set("testeringandmore", nil)
	f := r.SubTree("tester")
	fmt.Printf("'%v'  value: %v\n", f.key, f.value)
	// Output:
	// 'ster'  value: <nil>
}

func BenchmarkSubTree(b *testing.B) {
	b.StopTimer()
	r := New()
	r.Set("tester", nil)
	r.Set("testering", nil)
	r.Set("te", nil)
	r.Set("testeringandmore", nil)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_ = r.SubTree("tester")
	}
}
