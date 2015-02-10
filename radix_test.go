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
	r.Set("østern", nil)
	r.Set("øpen", nil)
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

func TestGetAllWithPrefix(t *testing.T) {
	r := New()
	r.Set("ðom", "ðom")
	r.Set("ðomum", "ðomum")
	r.Set("ðomulus", "ðomulus")

	roms := r.GetAllWithPrefix("ðom")

	if len(roms) != 3 {
		t.Error("not 3 results for ðom")
		t.Fail()
	}

	romus := r.GetAllWithPrefix("ðomu")

	if len(romus) != 2 {
		t.Error("not 2 results for ðomu")
		t.Fail()
	}

	c := 0
	for _, s := range []string{"ðomum", "ðomulus"} {
		for _, romu := range romus {
			if romu.(string) == s {
				c++
			}
		}
	}

	if c != 2 {
		t.Error("ðomum or ðomulus not contained in results for ðomu")
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
