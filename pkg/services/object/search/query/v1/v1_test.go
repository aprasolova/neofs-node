package query

import (
	"crypto/rand"
	"crypto/sha256"
	"testing"

	"github.com/nspcc-dev/neofs-api-go/pkg/container"
	objectSDK "github.com/nspcc-dev/neofs-api-go/pkg/object"
	"github.com/nspcc-dev/neofs-api-go/pkg/owner"
	"github.com/nspcc-dev/neofs-node/pkg/core/object"
	"github.com/stretchr/testify/require"
)

func testID(t *testing.T) *objectSDK.ID {
	cs := [sha256.Size]byte{}

	_, err := rand.Read(cs[:])
	require.NoError(t, err)

	id := objectSDK.NewID()
	id.SetSHA256(cs)

	return id
}

func testCID(t *testing.T) *container.ID {
	cs := [sha256.Size]byte{}

	_, err := rand.Read(cs[:])
	require.NoError(t, err)

	id := container.NewID()
	id.SetSHA256(cs)

	return id
}

func testOwnerID(t *testing.T) *owner.ID {
	w := new(owner.NEO3Wallet)

	_, err := rand.Read(w.Bytes())
	require.NoError(t, err)

	id := owner.NewID()
	id.SetNeo3Wallet(w)

	return id
}

func TestQ_Match(t *testing.T) {
	t.Run("object identifier equal", func(t *testing.T) {
		obj := object.NewRaw()

		id := testID(t)
		obj.SetID(id)

		q := New(
			NewIDEqualFilter(id),
		)

		require.True(t, q.Match(obj.Object()))

		obj.SetID(testID(t))

		require.False(t, q.Match(obj.Object()))
	})

	t.Run("container identifier equal", func(t *testing.T) {
		obj := object.NewRaw()

		id := testCID(t)
		obj.SetContainerID(id)

		q := New(
			NewContainerIDEqualFilter(id),
		)

		require.True(t, q.Match(obj.Object()))

		obj.SetContainerID(testCID(t))

		require.False(t, q.Match(obj.Object()))
	})

	t.Run("owner identifier equal", func(t *testing.T) {
		obj := object.NewRaw()

		id := testOwnerID(t)
		obj.SetOwnerID(id)

		q := New(
			NewOwnerIDEqualFilter(id),
		)

		require.True(t, q.Match(obj.Object()))

		obj.SetOwnerID(testOwnerID(t))

		require.False(t, q.Match(obj.Object()))
	})

	t.Run("attribute equal", func(t *testing.T) {
		obj := object.NewRaw()

		k, v := "key", "val"
		a := new(objectSDK.Attribute)
		a.SetKey(k)
		a.SetValue(v)

		obj.SetAttributes(a)

		q := New(
			NewFilterEqual(k, v),
		)

		require.True(t, q.Match(obj.Object()))

		a.SetKey(k + "1")

		require.False(t, q.Match(obj.Object()))
	})
}