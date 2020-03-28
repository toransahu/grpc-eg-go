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

	"github.com/toransahu/grpc-eg-go/machine"
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
