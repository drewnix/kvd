package kvcli

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"
)

func serveHelper(dels []string) error {
	serveCmd := ServeCmd()
	buf := bytes.NewBufferString("")
	serveCmd.SetOut(buf)
	serveCmd.SetArgs(dels)
	out, err := ioutil.ReadAll(buf)
	if err != nil {
		return err
	}
	err = serveCmd.Execute()
	if err != nil {
		return err
	}
	fmt.Println(out)
	return err
}

func TestServe(t *testing.T) {
	err := serveHelper([]string{""})
	if err != nil {
		t.Fatal(err)
	}
}
