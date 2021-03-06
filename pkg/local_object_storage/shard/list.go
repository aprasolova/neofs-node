package shard

import (
	"fmt"

	cid "github.com/nspcc-dev/neofs-api-go/pkg/container/id"
	"github.com/nspcc-dev/neofs-api-go/pkg/object"
	meta "github.com/nspcc-dev/neofs-node/pkg/local_object_storage/metabase"
	"go.uber.org/zap"
)

type ListContainersPrm struct{}

type ListContainersRes struct {
	containers []*cid.ID
}

func (r *ListContainersRes) Containers() []*cid.ID {
	return r.containers
}

func (s *Shard) List() (*SelectRes, error) {
	lst, err := s.metaBase.Containers()
	if err != nil {
		return nil, fmt.Errorf("can't list stored containers: %w", err)
	}

	res := new(SelectRes)
	filters := object.NewSearchFilters()

	for i := range lst {
		ids, err := meta.Select(s.metaBase, lst[i], filters) // consider making List in metabase
		if err != nil {
			s.log.Debug("can't select all objects",
				zap.Stringer("cid", lst[i]),
				zap.String("error", err.Error()))

			continue
		}

		res.addrList = append(res.addrList, ids...)
	}

	return res, nil
}

func (s *Shard) ListContainers(_ *ListContainersPrm) (*ListContainersRes, error) {
	containers, err := s.metaBase.Containers()
	if err != nil {
		return nil, fmt.Errorf("could not get list of containers: %w", err)
	}

	return &ListContainersRes{
		containers: containers,
	}, nil
}

func ListContainers(s *Shard) ([]*cid.ID, error) {
	res, err := s.ListContainers(&ListContainersPrm{})
	if err != nil {
		return nil, err
	}

	return res.Containers(), nil
}
