package netmap

import (
	"fmt"

	"github.com/nspcc-dev/neo-go/pkg/vm/stackitem"
	"github.com/nspcc-dev/neofs-node/pkg/morph/client"
)

// ConfigArgs groups the arguments
// of get config value test invoke call.
type ConfigArgs struct {
	key []byte
}

// EpochValues groups the stack parameters
// returned by get epoch number test invoke.
type ConfigValues struct {
	val interface{}
}

// SetKey sets binary key to configuration parameter.
func (c *ConfigArgs) SetKey(v []byte) {
	c.key = v
}

// Value returns configuration value.
func (c ConfigValues) Value() interface{} {
	return c.val
}

// Config performs the test invoke of get config value
// method of NeoFS Netmap contract.
func (c *Client) Config(args ConfigArgs, assert func(stackitem.Item) (interface{}, error)) (*ConfigValues, error) {
	items, err := c.client.TestInvoke(
		c.configMethod,
		args.key,
	)
	if err != nil {
		return nil, fmt.Errorf("could not perform test invocation (%s): %w",
			c.configMethod, err)
	}

	if ln := len(items); ln != 1 {
		return nil, fmt.Errorf("unexpected stack item count (%s): %d",
			c.configMethod, ln)
	}

	val, err := assert(items[0])
	if err != nil {
		return nil, fmt.Errorf("value type assertion failed: %w", err)
	}

	return &ConfigValues{
		val: val,
	}, nil
}

func IntegerAssert(item stackitem.Item) (interface{}, error) {
	return client.IntFromStackItem(item)
}

func StringAssert(item stackitem.Item) (interface{}, error) {
	return client.StringFromStackItem(item)
}
