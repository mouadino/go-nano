package discovery

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/nu7hatch/gouuid"
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

type LoadBalancer interface {
	// TODO: Endpoint([]Instance) ?
	Endpoint(*Service) (Endpoint, error)
}

var NoEndpointError = errors.New("No Endpoint")

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

func NewInstance(meta ServiceMetadata) (Instance, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return Instance{}, err
	}
	instance := Instance{ID: id.String(), Meta: meta}
	return instance, nil
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

func (m ServiceMetadata) Endpoint() Endpoint {
	return Endpoint(m["endpoint"].(string))
}

type Endpoint string

func (e Endpoint) Type() (string, error) {
	url, err := url.Parse(string(e))
	if err != nil {
		return "", err
	}
	return url.Scheme, nil
}
