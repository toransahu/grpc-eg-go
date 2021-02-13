//
// machine_test.go
// Copyright (C) 2020 Toran Sahu <toran.sahu@yahoo.com>
//
// Distributed under terms of the MIT license.
//

package server

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/toransahu/grpc-eg-go/machine"
	"github.com/toransahu/grpc-eg-go/mock_machine"
)

func TestExecute(t *testing.T) {
	s := MachineServer{}

	// set up test cases
	instruction_set_1 := []*machine.Instruction{}
	instruction_set_1 = append(instruction_set_1, &machine.Instruction{Operand: 5, Operator: "PUSH"})
	instruction_set_1 = append(instruction_set_1, &machine.Instruction{Operand: 6, Operator: "PUSH"})
	instruction_set_1 = append(instruction_set_1, &machine.Instruction{Operator: "ADD"})

	instruction_set_2 := []*machine.Instruction{}
	instruction_set_2 = append(instruction_set_2, &machine.Instruction{Operand: 5, Operator: "PUSH"})
	instruction_set_2 = append(instruction_set_2, &machine.Instruction{Operand: 6, Operator: "PUSH"})
	instruction_set_2 = append(instruction_set_2, &machine.Instruction{Operator: "MUL"})

	tests := []struct {
		instructions []*machine.Instruction
		want         float32
	}{
		{
			instructions: instruction_set_1,
			want:         11,
		},
		{
			instructions: instruction_set_2,
			want:         30,
		},
	}

	for _, tt := range tests {
		req := &machine.InstructionSet{Instructions: tt.instructions}
		resp, err := s.Execute(context.Background(), req)
		if err != nil {
			t.Errorf("ExecuteTest(%v) got unexpected error", tt.instructions)
		}
		if resp.Output != tt.want {
			t.Errorf("got Execute(%v)=%v, wanted %v", tt.instructions, resp.Output, tt.want)
		}
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
