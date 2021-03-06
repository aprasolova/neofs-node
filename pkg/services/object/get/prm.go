package getsvc

import (
	"hash"

	"github.com/nspcc-dev/neofs-api-go/pkg/client"
	objectSDK "github.com/nspcc-dev/neofs-api-go/pkg/object"
	"github.com/nspcc-dev/neofs-node/pkg/core/object"
	"github.com/nspcc-dev/neofs-node/pkg/services/object/util"
)

// Prm groups parameters of Get service call.
type Prm struct {
	commonPrm
}

// RangePrm groups parameters of GetRange service call.
type RangePrm struct {
	commonPrm

	rng *objectSDK.Range
}

// RangeHashPrm groups parameters of GetRange service call.
type RangeHashPrm struct {
	commonPrm

	hashGen func() hash.Hash

	rngs []*objectSDK.Range

	salt []byte
}

type RequestForwarder func(client.Client) (*objectSDK.Object, error)

// HeadPrm groups parameters of Head service call.
type HeadPrm struct {
	commonPrm
}

type commonPrm struct {
	objWriter ObjectWriter

	common *util.CommonPrm

	client.GetObjectParams

	forwarder RequestForwarder
}

// ChunkWriter is an interface of target component
// to write payload chunk.
type ChunkWriter interface {
	WriteChunk([]byte) error
}

// HeaderWriter is an interface of target component
// to write object header.
type HeaderWriter interface {
	WriteHeader(*object.Object) error
}

// ObjectWriter is an interface of target component to write object.
type ObjectWriter interface {
	HeaderWriter
	ChunkWriter
}

// SetObjectWriter sets target component to write the object.
func (p *Prm) SetObjectWriter(w ObjectWriter) {
	p.objWriter = w
}

// SetChunkWriter sets target component to write the object payload range.
func (p *RangePrm) SetChunkWriter(w ChunkWriter) {
	p.objWriter = &partWriter{
		chunkWriter: w,
	}
}

// SetRange sets range of the requested payload data.
func (p *RangePrm) SetRange(rng *objectSDK.Range) {
	p.rng = rng
}

// SetRangeList sets list of object payload ranges.
func (p *RangeHashPrm) SetRangeList(rngs []*objectSDK.Range) {
	p.rngs = rngs
}

// SetHashGenerator sets constructor of hashing algorithm.
func (p *RangeHashPrm) SetHashGenerator(v func() hash.Hash) {
	p.hashGen = v
}

// SetSalt sets binary salt to XOR object's payload ranges before hash calculation.
func (p *RangeHashPrm) SetSalt(salt []byte) {
	p.salt = salt
}

// SetCommonParameters sets common parameters of the operation.
func (p *commonPrm) SetCommonParameters(common *util.CommonPrm) {
	p.common = common
}

func (p *commonPrm) SetRequestForwarder(f RequestForwarder) {
	p.forwarder = f
}

// SetHeaderWriter sets target component to write the object header.
func (p *HeadPrm) SetHeaderWriter(w HeaderWriter) {
	p.objWriter = &partWriter{
		headWriter: w,
	}
}
