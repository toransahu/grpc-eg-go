//
// machine.go
// Copyright (C) 2020 Toran Sahu <toran.sahu@yahoo.com>
//
// Distributed under terms of the MIT license.
//

package server

import (
	"context"

	"github.com/toransahu/grpc-eg-go/machine"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MachineServer struct{}

// Execute runs the set of instructions given.
func (s *MachineServer) Execute(ctx context.Context, instructions *machine.InstructionSet) (*machine.Result, error) {
	return nil, status.Error(codes.Unimplemented, "Execute() not implemented yet")
}
