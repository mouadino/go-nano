/*
package zookeeper contains definition of discovery logic to announce rpc
instances in zookeeper and resolve them.

*/
package zookeeper

import (
	"bytes"
	"strings"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/mouadino/go-nano/discovery"
	"github.com/mouadino/go-nano/serializer"

	"github.com/samuel/go-zookeeper/zk"
)

type zookeeperAnnounceResolver struct {
	logger    *log.Logger
	endpoints []string
	serial    serializer.Serializer
	timeout   time.Duration
	chroot    string
	conn      *zk.Conn
	perms     int32

	mu    sync.RWMutex
	cache map[string]*discovery.Service
}

// Timeout option configure connection timeout to zookeeper, default 3 seconds.
func Timeout(timeout time.Duration) func(z *zookeeperAnnounceResolver) {
	return func(z *zookeeperAnnounceResolver) {
		z.timeout = timeout
	}
}

// Chroot option configure zookeeper path chroot, default "nano-services".
func Chroot(chroot string) func(z *zookeeperAnnounceResolver) {
	return func(z *zookeeperAnnounceResolver) {
		z.chroot = chroot
	}
}

// Perms option configure zookeeper path permission, default to open for all.
func Perms(perms int32) func(z *zookeeperAnnounceResolver) {
	return func(z *zookeeperAnnounceResolver) {
		z.perms = perms
	}
}

// Serializer option configure serializer to use when reading values instance data
// from zookeeper, default to json serializer.
func Serializer(serial serializer.Serializer) func(z *zookeeperAnnounceResolver) {
	return func(z *zookeeperAnnounceResolver) {
		z.serial = serial
	}
}

// New creates a AnnounceResolver for zookeeper.
func New(endpoints []string, options ...func(z *zookeeperAnnounceResolver)) discovery.AnnounceResolver {
	z := &zookeeperAnnounceResolver{
		endpoints: endpoints,
		logger:    log.New(),
		chroot:    "nano-services",
		timeout:   3 * time.Second,
		perms:     zk.PermAll,
		serial:    serializer.JSONSerializer{},
		cache:     make(map[string]*discovery.Service),
	}

	for _, opt := range options {
		opt(z)
	}
	return z
}

// Resolve the given service name to a Service structure.
func (z *zookeeperAnnounceResolver) Resolve(name string) (*discovery.Service, error) {
	err := z.ensureConn()
	if err != nil {
		return nil, err
	}
	z.mu.RLock()
	s, ok := z.cache[name]
	z.mu.RUnlock()
	if ok {
		return s, nil
	}
	s, err = z.resolve(name)
	if err != nil {
		z.mu.Lock()
		z.cache[name] = s
		z.mu.Unlock()
	}
	return s, err
}

func (z *zookeeperAnnounceResolver) resolve(name string) (*discovery.Service, error) {
	children, _, events, err := z.conn.ChildrenW(z.getPath(name))
	if err != nil {
		return nil, err
	}

	service := &discovery.Service{
		Name:      name,
		Instances: z.getInstances(name, children),
	}

	go z.watchEvents(events, func(ev zk.Event) {
		z.onPathChange(ev, service)
	})

	return service, nil
}

func (z *zookeeperAnnounceResolver) getInstances(name string, children []string) []discovery.Instance {
	instances := []discovery.Instance{}
	for _, id := range children {
		meta, err := z.getInstanceData(name, id)
		if err != nil {
			continue
		}
		instances = append(instances, discovery.Instance{
			ID:       id,
			Endpoint: meta["endpoint"].(string),
			Meta:     meta,
		})
	}
	return instances
}

func (z *zookeeperAnnounceResolver) getInstanceData(name, id string) (discovery.InstanceMeta, error) {
	data, _, err := z.conn.Get(z.getPath(name, id))
	if err != nil {
		z.logger.Errorf("zookeeper get failed for %s: %s", id, err)
		return nil, err
	}

	var meta discovery.InstanceMeta
	err = z.serial.Decode(bytes.NewReader(data), &meta)
	if err != nil {
		z.logger.Errorf("zookeeper metadata parse failed for %s: %s", id, err)
		return nil, err
	}
	return meta, nil
}

func (z *zookeeperAnnounceResolver) onPathChange(ev zk.Event, service *discovery.Service) {
	z.logger.WithFields(log.Fields{
		"event": ev.Type,
		"name":  service.Name,
		"path":  ev.Path,
	}).Debug("zookeeper path changed")

	children, _, err := z.conn.Children(z.getPath(service.Name))
	if err != nil {
		z.logger.Debug("zookeeper error getting children ", err)
		return
	}
	service.Instances = z.getInstances(service.Name, children)
}

// Announce instance for discovery in zookeeper under given name.
func (z *zookeeperAnnounceResolver) Announce(name string, instance discovery.Instance) error {
	err := z.ensureConn()
	if err != nil {
		return err
	}
	metadata, err := z.serial.Encode(instance.Meta)
	if err != nil {
		return err
	}

	return z.createNode(z.getPath(name, instance.ID), metadata)
}

func (z *zookeeperAnnounceResolver) ensureConn() error {
	if z.conn != nil {
		return nil
	}
	conn, events, err := zk.Connect(z.endpoints, z.timeout)
	if err != nil {
		return err
	}

	go z.watchEvents(events, func(ev zk.Event) {
		z.logger.WithFields(log.Fields{
			"event": ev.Type,
		}).Debug("zookeeper changed state")
	})
	z.conn = conn
	z.conn.SetLogger(z.logger)
	return err
}

func (z *zookeeperAnnounceResolver) watchEvents(events <-chan zk.Event, callback func(zk.Event)) {
	for ev := range events {
		callback(ev)
	}
}

func (z *zookeeperAnnounceResolver) getPath(keys ...string) string {
	path := append([]string{z.chroot}, keys...)
	return "/" + strings.Join(path, "/")
}

func (z *zookeeperAnnounceResolver) createNode(path string, data []byte) error {
	acl := zk.WorldACL(z.perms)
	flags := int32(0)

	// TODO: Refactor me.
	keys := strings.Split(path, "/")
	d := []byte{}
	var p string
	for i, _ := range keys {
		p = strings.Join(keys[:i+1], "/")
		if p == "" {
			continue
		}
		if p == path {
			d = data
			flags = int32(zk.FlagEphemeral)
		}
		_, err := z.conn.Create(p, d, flags, acl)
		if err != nil && err != zk.ErrNodeExists {
			return err
		}
	}
	return nil
}
