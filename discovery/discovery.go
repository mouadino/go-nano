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
	ID   string
	Meta ServiceMetadata
}

func NewInstance(meta ServiceMetadata) Instance {
	instance := Instance{ID: uuid.New(), Meta: meta}
	return instance
}

func (inst *Instance) String() string {
	return fmt.Sprintf("%s", inst.ID)
}

type ServiceMetadata map[string]interface{}

func NewServiceMetadata(endpoint string, meta map[string]interface{}) ServiceMetadata {
	res := ServiceMetadata{"endpoint": endpoint}
	for k, v := range meta {
		if k == "endpoint" {
			continue // TODO: Return error ?
		}
		res[k] = v
	}
	return res
}

func (m ServiceMetadata) Endpoint() string {
	return m["endpoint"].(string)
}
