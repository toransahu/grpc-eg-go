//
// machine_test.go
// Copyright (C) 2020 Toran Sahu <toran.sahu@yahoo.com>
//
// Distributed under terms of the MIT license.
//

package server

import (
	"io"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/toransahu/grpc-eg-go/machine"
	"github.com/toransahu/grpc-eg-go/mock_machine"
)

func TestExecute(t *testing.T) {
	s := MachineServer{}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockServerStream := mock_machine.NewMockMachine_ExecuteServer(ctrl)

	mockResults := []*machine.Result{}
	callRecv1 := mockServerStream.EXPECT().Recv().Return(&machine.Instruction{Operand: 1, Operator: "PUSH"}, nil)
	callRecv2 := mockServerStream.EXPECT().Recv().Return(&machine.Instruction{Operand: 2, Operator: "PUSH"}, nil).After(callRecv1)
	callRecv3 := mockServerStream.EXPECT().Recv().Return(&machine.Instruction{Operator: "MUL"}, nil).After(callRecv2)
	callRecv4 := mockServerStream.EXPECT().Recv().Return(&machine.Instruction{Operand: 3, Operator: "PUSH"}, nil).After(callRecv3)
	callRecv5 := mockServerStream.EXPECT().Recv().Return(&machine.Instruction{Operator: "ADD"}, nil).After(callRecv4)
	callRecv6 := mockServerStream.EXPECT().Recv().Return(&machine.Instruction{Operator: "FIB"}, nil).After(callRecv5)
	mockServerStream.EXPECT().Recv().Return(nil, io.EOF).After(callRecv6)
	mockServerStream.EXPECT().Send(gomock.Any()).DoAndReturn(
		func(result *machine.Result) error {
			mockResults = append(mockResults, result)
			return nil
		}).AnyTimes()
	wants := []float32{2, 5, 0, 1, 1, 2, 3, 5}

	err := s.Execute(mockServerStream)
	if err != nil {
		t.Errorf("Execute(%v) got unexpected error: %v", mockServerStream, err)
	}
	for i, result := range mockResults {
		got := result.GetOutput()
		want := wants[i]
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	}
}
