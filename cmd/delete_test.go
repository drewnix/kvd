package kvcli

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"
)

func deleteHelper(dels []string) error {
	deleteCmd := DeleteCmd()
	buf := bytes.NewBufferString("")
	deleteCmd.SetOut(buf)
	deleteCmd.SetArgs(dels)
	out, err := ioutil.ReadAll(buf)
	if err != nil {
		return err
	}
	err = deleteCmd.Execute()
	if err != nil {
		return err
	}
	fmt.Println(out)
	return err
}

func TestDelete(t *testing.T) {
	err := setHelper([]string{"test=true", "cat=meow", "dog=woof"})
	if err != nil {
		t.Fatal(err)
	}

	err = getHelper([]string{"test", "cat"})
	if err != nil {
		t.Fatal(err)
	}

	err = deleteHelper([]string{"cat", "dog"})
	if err != nil {
		t.Fatal(err)
	}
}
