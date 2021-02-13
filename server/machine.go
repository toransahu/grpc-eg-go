//
// machine.go
// Copyright (C) 2020 Toran Sahu <toran.sahu@yahoo.com>
//
// Distributed under terms of the MIT license.
//

package server

import (
	"fmt"
	"io"
	"log"

	"github.com/toransahu/grpc-eg-go/machine"
	"github.com/toransahu/grpc-eg-go/utils"
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
	FIB               = "FIB"
)

type MachineServer struct{}

// Execute runs the set of instructions given.
func (s *MachineServer) Execute(stream machine.Machine_ExecuteServer) error {
	var stack stack.Stack
	for {
		instruction, err := stream.Recv()
		if err == io.EOF {
			log.Println("EOF")
			output, popped := stack.Pop()
			if !popped {
				return status.Error(codes.Aborted, "Invalid sets of instructions. Execution aborted")
			}

			if err := stream.SendAndClose(&machine.Result{
				Output: output,
			}); err != nil {
				return err
			}

			return nil
		}
		if err != nil {
			return err
		}

		operand := instruction.GetOperand()
		operator := instruction.GetOperator()
		op_type := OperatorType(operator)

		fmt.Printf("Operand: %v, Operator: %v\n", operand, operator)

		switch op_type {
		case PUSH:
			stack.Push(float32(operand))
		case POP:
			stack.Pop()
		case ADD, SUB, MUL, DIV:
			item2, popped := stack.Pop()
			item1, popped := stack.Pop()

			if !popped {
				return status.Error(codes.Aborted, "Invalid sets of instructions. Execution aborted")
			}

			if op_type == ADD {
				stack.Push(item1 + item2)
			} else if op_type == SUB {
				stack.Push(item1 - item2)
			} else if op_type == MUL {
				stack.Push(item1 * item2)
			} else if op_type == DIV {
				stack.Push(item1 / item2)
			}

		default:
			return status.Errorf(codes.Unimplemented, "Operation '%s' not implemented yet", operator)
		}

	}
}

// ServerStreamingExecute runs the set of instructions given and streams a sequence of Results.
func (s *MachineServer) ServerStreamingExecute(instructions *machine.InstructionSet, stream machine.Machine_ServerStreamingExecuteServer) error {
	if len(instructions.GetInstructions()) == 0 {
		return status.Error(codes.InvalidArgument, "No valid instructions received")
	}

	var stack stack.Stack

	for _, instruction := range instructions.GetInstructions() {
		operand := instruction.GetOperand()
		operator := instruction.GetOperator()
		op_type := OperatorType(operator)

		log.Printf("Operand: %v, Operator: %v\n", operand, operator)

		switch op_type {
		case PUSH:
			stack.Push(float32(operand))
		case POP:
			stack.Pop()
		case FIB:
			n, popped := stack.Pop()

			if !popped {
				return status.Error(codes.Aborted, "Invalid sets of instructions. Execution aborted")
			}

			if op_type == FIB {
				for f := range utils.FibonacciRange(int(n)) {
					log.Println(float32(f))
					stream.Send(&machine.Result{Output: float32(f)})
				}
			}
		default:
			return status.Errorf(codes.Unimplemented, "Operation '%s' not implemented yet", operator)
		}
	}
	return nil
}
