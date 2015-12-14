package discovery

import (
	"fmt"

	"github.com/pborman/uuid"
)

type Resolver interface {
	Resolve(string) (*Service, error)
}

type Announcer interface {
	Announce(string, Instance) error
}

type AnnounceResolver interface {
	Resolver
	Announcer
}

type Service struct {
	Name      string
	Instances []Instance
}

func (s *Service) String() string {
	return fmt.Sprintf("%s [%d]", s.Name, len(s.Instances))
}

type Instance struct {
	ID       string
	Endpoint string
	Meta     InstanceMeta
}

func NewInstance(endpoint string, meta InstanceMeta) Instance {
	// FIXME: How about if meta have "endpoint" ?
	meta["endpoint"] = endpoint
	return Instance{
		ID:       uuid.New(),
		Endpoint: endpoint,
		Meta:     meta,
	}
}

type InstanceMeta map[string]interface{}
