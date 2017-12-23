package gob

import (
	"testing"
)

func TestTreeAddSimple(t *testing.T) {
	tree := NewTrie()
	tree.Add("ab", "hello")

	aTree, ok := tree.children['a']
	if !ok {
		t.Errorf("No 'a' found in trie")
		return
	}

	bTree, ok := aTree.children['b']
	if !ok {
		t.Errorf("No 'b' found in trie.a")
		return
	}

	if bTree.Value != "hello" {
		t.Errorf("No value set for trie.a.b")
		return
	}

	tree.Add("ac", "world")

	cTree, ok := aTree.children['c']
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

	aTree, ok := tree.children['a']
	if !ok {
		t.Errorf("No 'a' found in trie")
		return
	}
	aTree = aTree.children['/']

	colonTree, ok := aTree.children[':']
	if !ok {
		t.Errorf("No ':' found in trie.a")
		return
	}

	if colonTree.key != "bk" {
		t.Errorf("Incorrect key set for colon tree")
		return
	}
	colonTree = colonTree.children['/']

	cTree, ok := colonTree.children['c']
	if !ok {
		t.Errorf("No 'c' found in trie.a.:bk")
		return
	}

	if cTree.Value != "val" {
		t.Errorf("No value set for trie.a.:bk.c")
		return
	}

	tree.Add("a/:new/c", "other")

	colonTree, ok = tree.children['a'].children['/'].children[':']
	if !ok {
		t.Errorf("No 'a/:' found in trie")
		return
	}

	if colonTree.key != "new" {
		t.Errorf("Incorrect key on colonTree")
		return
	}
	colonTree = colonTree.children['/']

	cTree, ok = colonTree.children['c']
	if !ok {
		t.Errorf("No 'c' in trie.a/:new/")
		return
	}

	if cTree.Value != "other" {
		t.Errorf("Incorrect value on trie.a/:new/c")
		return
	}

	tree.Add("r/:rid", "special")

	rTree, ok := tree.children['r']
	if !ok {
		t.Errorf("No 'r' found in trie")
		return
	}
	rTree = rTree.children['/']

	colonTree, ok = rTree.children[':']
	if !ok {
		t.Errorf("No colon found")
		return
	}
	if colonTree.Value != "special" {
		t.Errorf("Incorrect value set on key")
		return
	}
	if colonTree.key != "rid" {
		t.Errorf("Incorrect key set on trie")
		return
	}
	if len(colonTree.children) != 0 {
		t.Errorf("children map is NOT empty!")
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

	tree.Add("r/:rid", "reddit")
	res = tree.Find("r/meirl")
	if res == nil {
		t.Errorf("Not finding key with param")
		return
	}
	if res.Value != "reddit" {
		t.Errorf("Incorrect value set for key")
		return
	}

	val, ok = res.Params["rid"]
	if !ok {
		t.Errorf("Param key not found in result")
		return
	}
	if val != "meirl" {
		t.Errorf("Incorrect param value stored")
		return
	}
}
