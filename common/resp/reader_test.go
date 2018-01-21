// Copyright 2018 The PigeonD Authors. All rights reserved.
// Use of this source code is governed by a GNU AGPL v3.0
// license that can be found in the AGPL V3 LICENSE file.

package resp_test

import (
	"bytes"
	"github.com/pigeond-io/pigeond/common/resp"
	"strconv"
	"testing"
)

func shouldBeParsedSuccessfully(t *testing.T, cmd string){
	t.Errorf("%q should be parsed successfully", cmd)
}

func shouldBeOk(t *testing.T, cmd string){
	t.Errorf("%q should be ok", cmd)
}

func shouldHaveNoArgs(t *testing.T, cmd string){
	t.Errorf("%q should have no args", cmd)
}

func shouldBeThis(t *testing.T, what string, expected interface{}, was interface{}){
	t.Errorf("Expected %s to be %v got this %v", what, expected, was)
}

func argLiteral(i int) string {
	var buffer bytes.Buffer
	buffer.WriteString("Arg#")
	buffer.WriteString(strconv.Itoa(i))
	return buffer.String()
}

func testRespCommand(t *testing.T, cmd *resp.Command, expectedAction string, expectedArgs ...string){
	if !cmd.Ok() {
		shouldBeOk(t, cmd.String())
	} else {
		if action := cmd.Action(); action != expectedAction {
			shouldBeThis(t, "Action", expectedAction, action)
		} else {
			if len(cmd.Args()) != len(expectedArgs) {
				shouldBeThis(t, "Args", len(expectedArgs), len(cmd.Args()))
			} else {
				for i, arg := range cmd.Args() {
					if string(arg.Bytes) != expectedArgs[i]  {
						shouldBeThis(t, argLiteral(i), expectedArgs[i], arg.Bytes)
					}
				}
			}
		}
	}
}

func testCommand(t *testing.T, command string, expectedAction string, expectedArgs ...string) {
	cmds, ok := resp.Read([]byte(command))
	if !ok {
		shouldBeParsedSuccessfully(t, command)
	}
	if len(cmds) != 1 {
		shouldBeThis(t, "commands", 1, len(cmds))
	} else {
		cmd := cmds[0]
		testRespCommand(t, cmd, expectedAction, expectedArgs...)
	}
}

func testMultiCommand(t *testing.T, command string, expectedCount int) ([]*resp.Command, bool) {
	cmds, ok := resp.Read([]byte(command))
	if !ok {
		shouldBeParsedSuccessfully(t, command)
	}
	if len(cmds) != expectedCount {
		shouldBeThis(t, "commands", expectedCount, len(cmds))
	}
	return cmds, ok
}

func TestCommandWithoutArgsUsingString(t *testing.T) {
	testCommand(t, "+SUBSCRIBE\r\n", "SUBSCRIBE")
}

func TestCommandWithoutArgsUsingBulkString(t *testing.T){
	testCommand(t, "$9\r\nSUBSCRIBE\r\n", "SUBSCRIBE")
}

func TestCommandWithoutArgsUsingArray(t *testing.T){
	testCommand(t, "*1\r\n$9\r\nSUBSCRIBE\r\n", "SUBSCRIBE")
}

func TestCommandWithArgs(t *testing.T) {
	testCommand(t, "*2\r\n$9\r\nSUBSCRIBE\r\n$7\r\nMyTopic\r\n", "SUBSCRIBE", "MyTopic")
	testCommand(t, "*2\r\n+SUBSCRIBE\r\n+MyTopic\r\n", "SUBSCRIBE", "MyTopic")
}

func TestMultipleCommands(t *testing.T){
	cmds, ok := testMultiCommand(t, "*2\r\n$9\r\nSUBSCRIBE\r\n$7\r\nMyTopic\r\n*3\r\n+SUBSCRIBE\r\n+MyTopic\r\n+MySubTopic\r\n*1\r\n$9\r\nSUBSCRIBE\r\n+UNSUBSCRIBE\r\n", 4)
	if ok {
		testRespCommand(t, cmds[0], "SUBSCRIBE", "MyTopic")
		testRespCommand(t, cmds[1], "SUBSCRIBE", "MyTopic", "MySubTopic")
		testRespCommand(t, cmds[2], "SUBSCRIBE")
		testRespCommand(t, cmds[3], "UNSUBSCRIBE")
	}
}

// TODO
func TestInvalidCommand(t *testing.T) {
	t.Skip("Pending: Test Invalid Commands")
}