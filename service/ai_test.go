package service

import (
	"context"
	"fmt"
	"testing"
)

// test ai chat
func TestChat(t *testing.T) {
	l, err := Chat(context.Background(), "https://www.github.com/pocketbase/pocketbase")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("%v\n\n", l.What)
	for _, v := range l.Why {
		fmt.Println("==========")
		fmt.Printf("%v\n", v)
		fmt.Println("==========")
	}
	for _, v := range l.How {
		fmt.Println("==========")
		fmt.Printf("%v\n", v)
		fmt.Println("==========")
	}
	fmt.Printf("%+v", l.Other)
}
