//
// machine.go
// Copyright (C) 2020 Toran Sahu <toran.sahu@yahoo.com>
//
// Distributed under terms of the MIT license.
//

package server

import (
	"context"
	"fmt"

	"github.com/toransahu/grpc-eg-go/machine"
	"github.com/toransahu/grpc-eg-go/utils/stack"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OperatorType string

const (
	PUSH OperatorType = "PUSH"
	POP               = "POP"
	ADD               = "ADD"
	SUB               = "SUB"
	MUL               = "MUL"
	DIV               = "DIV"
)

type MachineServer struct{}

// Execute runs the set of instructions given.
func (s *MachineServer) Execute(ctx context.Context, instructions *machine.InstructionSet) (*machine.Result, error) {
	if len(instructions.GetInstructions()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "No valid instructions received")
	}

	var stack stack.Stack

	for _, instruction := range instructions.GetInstructions() {
		operand := instruction.GetOperand()
		operator := instruction.GetOperator()
		op_type := OperatorType(operator)

		fmt.Printf("Operand: %v, Operator: %v\n", operand, operator)

		switch op_type {
		case PUSH:
			stack.Push(float32(operand))
		case POP:
			stack.Pop()
		case ADD, MUL, DIV:
			item2, popped := stack.Pop()
			item1, popped := stack.Pop()

			if !popped {
				return &machine.Result{}, status.Error(codes.Aborted, "Invalide sets of instructions. Execution aborted")
			}

			if op_type == ADD {
				stack.Push(item1 + item2)
			} else if op_type == MUL {
				stack.Push(item1 * item2)
			} else if op_type == DIV {
				stack.Push(item1 / item2)
			}

		default:
			return nil, status.Errorf(codes.Unimplemented, "Operation '%s' not implemented yet", operator)
		}

	}

	item, popped := stack.Pop()
	if !popped {
		return &machine.Result{}, status.Error(codes.Aborted, "Invalide sets of instructions. Execution aborted")
	}
	return &machine.Result{Output: item}, nil
}
