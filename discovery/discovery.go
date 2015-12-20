/*
package discovery defines common structs and interfaces for discovering and announcing
RPC services.
*/
package discovery

import (
	"fmt"

	"github.com/pborman/uuid"
)

// Resolver describes something that can resolve services from name.
type Resolver interface {
	Resolve(string) (*Service, error)
}

// Announcer describes something that can announce a service instance to make it
// available for discovery.
type Announcer interface {
	Announce(string, Instance) error
}

// AnnounceResolver groups Announcer and Resolver interfaces.
type AnnounceResolver interface {
	Resolver
	Announcer
}

// Service represents an RPC service, which is an abstract definition of
// RPC instances announced under the same name.
type Service struct {
	Name      string
	Instances []Instance
}

// String returns a representation string of service.
func (s *Service) String() string {
	return fmt.Sprintf("%s [%d]", s.Name, len(s.Instances))
}

// Instance represents an RPC instance, usually this map to one process.
type Instance struct {
	ID       string
	Endpoint string
	Meta     InstanceMeta
}

// NewInstance creates a new Instance.
func NewInstance(endpoint string, meta InstanceMeta) Instance {
	// FIXME: How about if meta have "endpoint" ?
	if meta == nil {
		meta = make(InstanceMeta)
	}
	meta["endpoint"] = endpoint
	return Instance{
		ID:       uuid.New(),
		Endpoint: endpoint,
		Meta:     meta,
	}
}

// InstanceMeta represents metadata associated with an instance.
type InstanceMeta map[string]interface{}
