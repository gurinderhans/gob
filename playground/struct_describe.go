package main

import "fmt"
import "encoding/json"

type Chat struct {
	Name      string   `json:"name"`
	UserCount int      `json:"userCount"`
	UserIDs   []string `json:"userIds"`
}

func main() {
	chat := Chat{
		Name:      "hello",
		UserCount: 23,
		UserIDs:   []string{"user-1", "user-2"},
	}
	b, err := json.Marshal(chat)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
}
