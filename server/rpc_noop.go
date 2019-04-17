package server

import (
	"../rpc/types"
)

func init() {
	// init on module load
	registerEndpoint(NewNoOpEndpoint())
}

type NoOpEndpoint struct {
	server *Instance
}

func NewNoOpEndpoint() *NoOpEndpoint {
	return &NoOpEndpoint{}
}

func (endpoint *NoOpEndpoint) Execute(args *types.ReadRequest, resp *types.ReadResponse) error {
	return nil
}

func (endpoint *NoOpEndpoint) register(opts *EndpointOpts) error {
	if err := opts.server.rpc.RegisterName(endpoint.name().String(), endpoint); err != nil {
		return err
	}
	endpoint.server = opts.server
	return nil
}

func (endpoint *NoOpEndpoint) name() EndpointName {
	return EndpointName(types.EndpointNoOp)
}
