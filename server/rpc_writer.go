package server

import (
	"../rpc/types"
	"log"
)

type WriterEndpoint struct {
}

func NewWriterEndpoint() *WriterEndpoint {
	return &WriterEndpoint{}
}

func (endpoint *WriterEndpoint) Execute(args *types.WriteRequest, resp *types.WriteResponse) error {
	log.Printf("args %+v", args)
	numTimes := len(args.Times)
	numValues := len(args.Values)
	if numTimes != numValues {
		resp.Error = &types.RpcErrorNumTimeValuePairsMisMatch
		return nil
	}
	resp.Num = numTimes
	resp.Error = &types.RpcErrorNotImplemented
	// @todo implement real write
	return nil
}

func (endpoint *WriterEndpoint) register(opts *EndpointOpts) error {
	if err := opts.server.rpc.RegisterName(endpoint.name().String(), endpoint); err != nil {
		return err
	}
	return nil
}

func (endpoint *WriterEndpoint) name() EndpointName {
	return EndpointName(types.EndpointWriter)
}

func init() {
	endpoint := NewWriterEndpoint()
	endpointsMux.Lock()
	endpoints = append(endpoints, endpoint)
	endpointsMux.Unlock()
}
