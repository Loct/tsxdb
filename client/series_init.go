package client

import (
	"github.com/RobinUS2/tsxdb/rpc/types"
	"github.com/RobinUS2/tsxdb/tools"
	"github.com/pkg/errors"
	"sync/atomic"
)

// init series metadata on the server side (only transmits names, tags, etc. once instead of each time)

func (series *Series) Create() (id uint64, err error) {
	id = series.Id()
	if id > 0 {
		return id, nil
	}
	conn, err := series.client.GetConnection()
	if err != nil {
		return 0, err
	}
	defer panicOnErrorClose(conn.Close)
	return series.Init(conn)
}

func (series *Series) ResetInit() {
	// setting the ID to 0 means metadata will be exchanged with server again after next call to series.Init()
	atomic.StoreUint64(&series.id, 0)
}

func (series *Series) Init(conn *ManagedConnection) (id uint64, err error) {
	// Note: this function does not close the connection, need to do in function that uses it

	// fast path
	id = atomic.LoadUint64(&series.id)
	if id > 0 {
		// already initialised
		return id, nil
	}

	// verify connection
	if conn == nil {
		return 0, errors.New("missing connection")
	}

	// max 1 init at a time for 1 series
	series.initMux.Lock()
	defer series.initMux.Unlock()

	// check again already sent? (could be done during waiting of the lock)
	id = atomic.LoadUint64(&series.id)
	if id > 0 {
		// already initialised
		return id, nil
	}

	// request
	request := types.SeriesMetadataRequest{
		SeriesCreateMetadata: types.SeriesCreateMetadata{
			SeriesMetadata: types.SeriesMetadata{
				Namespace: series.Namespace(),
				Tags:      series.Tags(),
				Name:      series.Name(),
				Ttl:       series.TTL(),
			},
			SeriesCreateIdentifier: types.SeriesCreateIdentifier(tools.RandomInsecureIdentifier()),
		},
		SessionTicket: conn.getSessionTicket(),
	}

	// execute
	var response *types.SeriesMetadataResponse
	if err := conn.client.Call(types.EndpointSeriesMetadata.String()+"."+types.MethodName, request, &response); err != nil {
		return 0, err
	}
	if response.Error != nil {
		return 0, response.Error.Error()
	}

	// store id
	atomic.StoreUint64(&series.id, response.Id)

	return response.Id, nil
}
