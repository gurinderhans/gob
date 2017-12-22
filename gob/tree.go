package gob

import (
	"bytes"
)

type Tree interface {
	Add(key string, val interface{})
	Find(key string) *Tree
}

type Trie struct {
	Value    interface{}
	Key      string
	Children map[rune]*Trie
	Params   map[string]interface{}
}

func NewTrie() *Trie {
	return &Trie{Children: make(map[rune]*Trie)}
}

func (t *Trie) Add(key string, val interface{}) {
	var isParamKey bool
	var paramKey bytes.Buffer

	looper := t
	for _, r := range key {
		if isParamKey && r == '/' {
			isParamKey = false
			// TODO: maybe warn if `Key` already exists?
			looper.Key = paramKey.String()
			paramKey.Reset()
		} else if isParamKey {
			paramKey.WriteRune(r)
			continue
		}
		if r == ':' {
			isParamKey = true
		}
		if _, exists := looper.Children[r]; !exists {
			looper.Children[r] = NewTrie()
		}
		looper = looper.Children[r]
	}
	looper.Value = val
}

func (t *Trie) Find(key string) *Trie {
	var isParamValue bool
	var paramValue bytes.Buffer

	params := make(map[string]interface{})

	looper := t
	for _, r := range key {
		if isParamValue && r == '/' {
			isParamValue = false
			params[looper.Key] = paramValue.String()
			paramValue.Reset()
		} else if isParamValue {
			paramValue.WriteRune(r)
			continue
		}

		if trie, exists := looper.Children[r]; exists {
			looper = trie
		} else if trie, exists := looper.Children[':']; exists {
			looper = trie
			isParamValue = true
			paramValue.WriteRune(r)
			continue
		} else {
			return nil
		}
	}
	if looper.Value == nil {
		return nil
	}
	return &Trie{
		Params: params,
		Value:  looper.Value,
	}
}
