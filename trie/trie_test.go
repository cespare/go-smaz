package trie

import (
	"testing"
)

func TestPutGet(t *testing.T) {
	tr := New()
	ks := []string{"foo", "fool", "fools", "blah"}
	for i, s := range ks {
		if ok := tr.Put([]byte(s), i); !ok {
			t.Fatal("expected false when putting a fresh key")
		}
		if ok := tr.Put([]byte("foo"), 100); ok {
			t.Fatal("expected true when putting an existing key")
		}
	}
	for i, s := range ks {
		n, ok := tr.Get([]byte(s))
		if !ok {
			t.Fatalf("expected to find %s in trie, but did not", s)
		}
		want := i
		if i == 0 {
			want = 100
		}
		if n != want {
			t.Fatalf("want %d; got %d", want, n)
		}
	}
	for _, s := range []string{"f", "fo", "b", "fooll"} {
		if _, ok := tr.Get([]byte(s)); ok {
			t.Fatalf("did not expect to find %s in trie", s)
		}
	}
}
