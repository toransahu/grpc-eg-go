//
// machine.go
// Copyright (C) 2020 Toran Sahu <toran.sahu@yahoo.com>
//
// Distributed under terms of the MIT license.
//

package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/toransahu/grpc-eg-go/machine"
	"google.golang.org/grpc"
)

var (
	serverAddr = flag.String("server_addr", "localhost:9111", "The server address in the format of host:port")
)

func runExecute(client machine.MachineClient, instructions *machine.InstructionSet) {
	log.Printf("Executing %v", instructions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := client.Execute(ctx, instructions)
	if err != nil {
		log.Fatalf("%v.Execute(_) = _, %v: ", client, err)
	}
	log.Println(result)
}

func main() {
	flag.Parse()
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithBlock())
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := machine.NewMachineClient(conn)

	// try Execute()
	instructions := []*machine.Instruction{}
	instructions = append(instructions, &machine.Instruction{Operand: 5, Operator: "PUSH"})
	instructions = append(instructions, &machine.Instruction{Operand: 6, Operator: "PUSH"})
	instructions = append(instructions, &machine.Instruction{Operator: "MUL"})
	runExecute(client, &machine.InstructionSet{Instructions: instructions})
}
