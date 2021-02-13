//
// machine_mock_test.go
// Copyright (C) 2020 Toran Sahu <toran.sahu@yahoo.com>
//
// Distributed under terms of the MIT license.
//

package mock_machine_test

import (
	context "context"
	"fmt"
	"log"
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
	"github.com/toransahu/grpc-eg-go/machine"
	mock_machine "github.com/toransahu/grpc-eg-go/mock_machine"
)

func testExecute(t *testing.T, client machine.MachineClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	instructions := []*machine.Instruction{}
	instructions = append(instructions, &machine.Instruction{Operand: 5, Operator: "PUSH"})
	instructions = append(instructions, &machine.Instruction{Operand: 6, Operator: "PUSH"})
	instructions = append(instructions, &machine.Instruction{Operator: "MUL"})

	result, err := client.Execute(ctx, &machine.InstructionSet{Instructions: instructions})
	if err != nil {
		log.Fatalf("%v.Execute(_) = _, %v: ", client, err)
	}
	log.Println(result)
	got := result.GetOutput()
	want := float32(30)
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestExecute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockMachineClient := mock_machine.NewMockMachineClient(ctrl)

	instructions := []*machine.Instruction{}
	instructions = append(instructions, &machine.Instruction{Operand: 5, Operator: "PUSH"})
	instructions = append(instructions, &machine.Instruction{Operand: 6, Operator: "PUSH"})
	instructions = append(instructions, &machine.Instruction{Operator: "MUL"})

	instructionSet := &machine.InstructionSet{Instructions: instructions}

	mockMachineClient.EXPECT().Execute(
		gomock.Any(),   // context
		instructionSet, // rpc uniary message
	).Return(&machine.Result{Output: 30}, nil)

	testExecute(t, mockMachineClient)
}

func testServerStreamingExecute(t *testing.T, mockMachineClient machine.MachineClient, instructionSet *machine.InstructionSet) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, err := mockMachineClient.ServerStreamingExecute(ctx, instructionSet)
	if err != nil {
		log.Fatalf("%v", err)
	}
	result, err := stream.Recv()
	if err != nil {
		return err
	}
	got := result.GetOutput()
	want := float32(0)
	if got != want {
		return fmt.Errorf("stream.Recv() = %v, want %v", got, want)
	}
	return nil
}

func TestServerStreamingExecute(t *testing.T) {
	instructions := []*machine.Instruction{}
	instructions = append(instructions, &machine.Instruction{Operand: 1, Operator: "PUSH"})
	instructions = append(instructions, &machine.Instruction{Operator: "FIB"})
	instructionSet := &machine.InstructionSet{Instructions: instructions}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockMachineClient := mock_machine.NewMockMachineClient(ctrl)
	clientStream := mock_machine.NewMockMachine_ServerStreamingExecuteClient(ctrl)

	clientStream.EXPECT().Recv().Return(&machine.Result{Output: 0}, nil)

	mockMachineClient.EXPECT().ServerStreamingExecute(
		gomock.Any(),   // context
		instructionSet, // rpc uniary message
	).Return(clientStream, nil)

	if err := testServerStreamingExecute(t, mockMachineClient, instructionSet); err != nil {
		t.Fatalf("Test failed: %v", err)
	}
}
