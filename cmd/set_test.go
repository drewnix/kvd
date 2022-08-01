package kvcli

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"
)

func setHelper(sets []string) error {
	setCmd := SetCmd()
	buf := bytes.NewBufferString("")
	setCmd.SetOut(buf)
	setCmd.SetArgs(sets)
	out, err := ioutil.ReadAll(buf)
	if err != nil {
		return err
	}
	err = setCmd.Execute()
	if err != nil {
		return err
	}
	fmt.Println(out)
	return err
}

func TestSet(t *testing.T) {
	cmd := SetCmd()
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{"abc=bbd", "foo=bar"})
	out, err := ioutil.ReadAll(b)
	cmd.Execute()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(out))
	//if string(out) != "hi-via-args" {
	//	t.Fatalf("expected \"%s\" got \"%s\"", "hi-via-args", string(out))
	//}
	//args := os.Args[0:1]
	//args = append(args, "foo", "bar")
	//err := run(args)
	//if err != nil {
	//	t.Errorf("test yielded error: %s", err.Error())
	//}
}
