package kvcli

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"
)

func getHelper(gets []string) error {
	getCmd := GetCmd()
	buf := bytes.NewBufferString("")
	getCmd.SetOut(buf)
	getCmd.SetArgs(gets)
	out, err := ioutil.ReadAll(buf)
	if err != nil {
		return err
	}
	err = getCmd.Execute()
	if err != nil {
		return err
	}
	fmt.Println(out)
	return err
}

func TestGet(t *testing.T) {
	err := setHelper([]string{"test=true", "cat=meow", "dog=woof"})
	if err != nil {
		t.Fatal(err)
	}

	err = getHelper([]string{"test", "cat"})
	if err != nil {
		t.Fatal(err)
	}
}
