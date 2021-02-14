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

	mockResult := &machine.Result{}
	callRecv1 := mockServerStream.EXPECT().Recv().Return(&machine.Instruction{Operand: 5, Operator: "PUSH"}, nil)
	callRecv2 := mockServerStream.EXPECT().Recv().Return(&machine.Instruction{Operand: 6, Operator: "PUSH"}, nil).After(callRecv1)
	callRecv3 := mockServerStream.EXPECT().Recv().Return(&machine.Instruction{Operator: "MUL"}, nil).After(callRecv2)
	mockServerStream.EXPECT().Recv().Return(nil, io.EOF).After(callRecv3)
	mockServerStream.EXPECT().SendAndClose(gomock.Any()).DoAndReturn(
		func(result *machine.Result) error {
			mockResult = result
			return nil
		})

	err := s.Execute(mockServerStream)
	if err != nil {
		t.Errorf("Execute(%v) got unexpected error: %v", mockServerStream, err)
	}
	got := mockResult.GetOutput()
	want := float32(30)
	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}
}

func TestServerStreamingExecute(t *testing.T) {
	s := MachineServer{}

	// set up test table
	tests := []struct {
		instructions []*machine.Instruction
		want         []float32
	}{
		{
			instructions: []*machine.Instruction{
				{Operand: 5, Operator: "PUSH"},
				{Operator: "FIB"},
			},
			want: []float32{0, 1, 1, 2, 3, 5},
		},
		{
			instructions: []*machine.Instruction{
				{Operand: 6, Operator: "PUSH"},
				{Operator: "FIB"},
			},
			want: []float32{0, 1, 1, 2, 3, 5, 8},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockServerStream := mock_machine.NewMockMachine_ServerStreamingExecuteServer(ctrl)
	for _, tt := range tests {
		mockResults := []*machine.Result{}
		mockServerStream.EXPECT().Send(gomock.Any()).DoAndReturn(
			func(result *machine.Result) error {
				mockResults = append(mockResults, result)
				return nil
			}).AnyTimes()

		req := &machine.InstructionSet{Instructions: tt.instructions}

		err := s.ServerStreamingExecute(req, mockServerStream)
		if err != nil {
			t.Errorf("ServerStreamingExecute(%v) got unexpected error: %v", req, err)
		}
		for i, result := range mockResults {
			got := result.GetOutput()
			want := tt.want[i]
			if got != want {
				t.Errorf("got %v, want %v", got, want)
			}
		}
	}
}
