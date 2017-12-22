package gob

import (
	"bytes"
)

type Tree interface {
	Add(key string, val interface{})
	Find(key string) *Tree
}

type Trie struct {
	Value  interface{}
	Params map[string]string

	key      string
	children map[rune]*Trie
}

func NewTrie() *Trie {
	return &Trie{children: make(map[rune]*Trie)}
}

func (t *Trie) Add(key string, val interface{}) {
	var isParamKey bool
	var paramKey bytes.Buffer

	looper := t
	for _, r := range key {
		if isParamKey && r == '/' {
			isParamKey = false
			looper.key = paramKey.String() // TODO: maybe warn if `Key` already exists?
			paramKey.Reset()
		} else if isParamKey {
			paramKey.WriteRune(r)
			continue
		}
		if r == ':' {
			isParamKey = true
		}
		if _, exists := looper.children[r]; !exists {
			looper.children[r] = NewTrie()
		}
		looper = looper.children[r]
	}
  if isParamKey {
    looper.key = paramKey.String()
  }
	looper.Value = val
}

func (t *Trie) Find(key string) *Trie {
	var isParamValue bool
	var paramValue bytes.Buffer

	params := make(map[string]string)
	looper := t
	for _, r := range key {
		if isParamValue && r == '/' {
			isParamValue = false
			params[looper.key] = paramValue.String()
			paramValue.Reset()
		} else if isParamValue {
			paramValue.WriteRune(r)
			continue
		}

		if trie, exists := looper.children[r]; exists {
			looper = trie
		} else if trie, exists := looper.children[':']; exists {
			looper = trie
			isParamValue = true
			paramValue.WriteRune(r)
			continue
		} else {
			return nil
		}
	}
  if isParamValue {
    params[looper.key] = paramValue.String()
  }
	if looper.Value == nil {
		return nil
	}
	return &Trie{
		Params: params,
		Value:  looper.Value,
	}
}
