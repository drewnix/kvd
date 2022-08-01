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
	setCmd := SetCmd()
	b := bytes.NewBufferString("")
	setCmd.SetOut(b)
	setCmd.SetArgs([]string{"andrew=king", "echo=prince", "ann=queen"})
	out, err := ioutil.ReadAll(b)
	fmt.Println("Executing")
	setCmd.Execute()

	getCmd := GetCmd()
	getCmd.SetOut(b)
	getCmd.SetArgs([]string{"andrew", "echo"})
	out, err = ioutil.ReadAll(b)
	getCmd.Execute()
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
