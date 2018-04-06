package main

import (
	"testing"
	"os/exec"
	"fmt"
	"os"
)

func TestCalledByMain(t *testing.T) {
	for i := 0; i < 100; i++ {
		clifile := CalledByMain()
		t.Log("client file name " + clifile)
		var err error
		shcmd := "diff " + "./" + clifile + " ../server/" + clifile
		cmd := exec.Command("sh", "-c", shcmd)
		var output []byte
		if output, err = cmd.Output(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if string(output) != "" {
			t.Error("diff!")
		}
	}

	shcmd := "rm 127*"
	cmd := exec.Command("sh", "-c", shcmd)
	var err error
	if _, err = cmd.Output(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	shcmd = "rm ../server/127*"
	cmd = exec.Command("sh", "-c", shcmd)

	if _, err = cmd.Output(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
