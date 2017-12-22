package gob

import (
	"testing"
)

func TestTreeAddSimple(t *testing.T) {
	tree := NewTrie()
	tree.Add("ab", "hello")

	aTree, ok := tree.Children['a']
	if !ok {
		t.Errorf("No 'a' found in trie")
		return
	}

	bTree, ok := aTree.Children['b']
	if !ok {
		t.Errorf("No 'b' found in trie.a")
		return
	}

	if bTree.Value != "hello" {
		t.Errorf("No value set for trie.a.b")
		return
	}

	tree.Add("ac", "world")

	cTree, ok := aTree.Children['c']
	if !ok {
		t.Errorf("No 'c' found in trie.a")
		return
	}

	if cTree.Value != "world" {
		t.Errorf("No value set for trie.a.c")
		return
	}
}

func TestTreeAddComplex(t *testing.T) {
	tree := NewTrie()
	tree.Add("a/:bk/c", "val")

	aTree, ok := tree.Children['a']
	if !ok {
		t.Errorf("No 'a' found in trie")
		return
	}
	aTree = aTree.Children['/']

	colonTree, ok := aTree.Children[':']
	if !ok {
		t.Errorf("No ':' found in trie.a")
		return
	}

	if colonTree.Key != "bk" {
		t.Errorf("Incorrect key set for colon tree")
		return
	}
	colonTree = colonTree.Children['/']

	cTree, ok := colonTree.Children['c']
	if !ok {
		t.Errorf("No 'c' found in trie.a.:bk")
		return
	}

	if cTree.Value != "val" {
		t.Errorf("No value set for trie.a.:bk.c")
		return
	}

	tree.Add("a/:new/c", "other")

	colonTree, ok = tree.Children['a'].Children['/'].Children[':']
	if !ok {
		t.Errorf("No 'a/:' found in trie")
		return
	}

	if colonTree.Key != "new" {
		t.Errorf("Incorrect key on colonTree")
		return
	}
	colonTree = colonTree.Children['/']

	cTree, ok = colonTree.Children['c']
	if !ok {
		t.Errorf("No 'c' in trie.a/:new/")
		return
	}

	if cTree.Value != "other" {
		t.Errorf("Incorrect value on trie.a/:new/c")
		return
	}
}

func TestTreeFindSimple(t *testing.T) {
	tree := NewTrie()
	tree.Add("/hello", "val")

	res := tree.Find("/rand")
	if res != nil {
		t.Errorf("Didn't fail properly when searching for non-existent key")
		return
	}

	res = tree.Find("/hell")
	if res != nil {
		t.Errorf("Didn't fail properly when searching for sub-existent key")
		return
	}

	res = tree.Find("/hello")
	if res.Value != "val" {
		t.Errorf("Incorrect Value was set on trie key")
		return
	}

	res = tree.Find("/hello/")
	if res != nil {
		t.Errorf("Didn't fail properly when searching for non-existent key")
		return
	}
}

func TestTreeFindComplex(t *testing.T) {
	tree := NewTrie()
	tree.Add("/a/:id/c", "po")

	res := tree.Find("/a/me")
	if res != nil {
		t.Errorf("Incorrectly finding subkey")
		return
	}

	res = tree.Find("/a/you/c")
	if res == nil {
		t.Errorf("Not finding key with param")
		return
	}

	if res.Value != "po" {
		t.Errorf("Incorrect value returned for key")
		return
	}

	val, ok := res.Params["id"]
	if !ok {
		t.Errorf("Param key not found in result")
		return
	}
	if val != "you" {
		t.Errorf("Incorrect param value stored")
		return
	}

	tree.Add("/a/hey/c", "yo")

	res = tree.Find("/a/hey/c")

	if res == nil {
		t.Errorf("Not finding key w/o param")
		return
	}

	if res.Value != "yo" {
		t.Errorf("Incorrect value for key")
		return
	}
}
