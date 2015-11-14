package discovery

import (
	"fmt"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
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
}

func DefaultZooKeeperAnnounceResolver(endpoints []string) AnnounceResolver {
	return CustomZooKeeperAnnounceResolver(
		endpoints,
		"nano-services",
		3*time.Second,
		zk.PermAll,
		serializer.JSONSerializer{},
	)
}

func CustomZooKeeperAnnounceResolver(endpoints []string, chroot string, timeout time.Duration, perms int32, serial serializer.Serializer) AnnounceResolver {
	return &zookeeperAnnounceResolver{
		logger:    log.New(), // TODO: Pass logger (Namespacing ?)
		endpoints: endpoints,
		timeout:   timeout,
		chroot:    chroot,
		serial:    serial,
		perms:     perms,
	}
}

// TODO: Cache result.
func (z *zookeeperAnnounceResolver) Resolve(name string) (*Service, error) {
	err := z.ensureConn()
	if err != nil {
		return nil, err
	}
	children, _, err := z.conn.Children(z.getPath(name))

	if err != nil {
		return nil, err
	}

	instances := []Instance{}
	for _, id := range children {
		// TODO: Watch events and change Instances dynamically.
		data, _, err := z.conn.Get(z.getPath(name, id))
		if err != nil {
			z.logger.Error("zookeeper get failed", err)
			continue
		}
		var meta ServiceMetadata
		err = z.serial.Decode(data, &meta)
		if err != nil {
			z.logger.Error("zookeeper metadata parse failed", err)
			continue
		}
		instances = append(instances, Instance{
			ID:   id,
			Meta: meta,
		})
	}
	service := &Service{
		Name:      name,
		Instances: instances,
	}
	return service, nil
}

func (z *zookeeperAnnounceResolver) Announce(name string, instance Instance) error {
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
	go z.watchEvents(events)
	z.conn = conn
	z.conn.SetLogger(z.logger)
	return err
}

func (z *zookeeperAnnounceResolver) watchEvents(events <-chan zk.Event) {
	for ev := range events {
		z.logger.WithFields(log.Fields{
			"event": ev.Type,
		}).Debug("zookeeper changed state")
	}
}

func (z *zookeeperAnnounceResolver) getPath(keys ...string) string {
	path := append([]string{z.chroot}, keys...)
	return "/" + strings.Join(path, "/")
}

func (z *zookeeperAnnounceResolver) createNode(path string, data []byte) error {
	acl := zk.WorldACL(z.perms)
	flags := int32(0)

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
		fmt.Printf("Creating %s %s\n", p, d)
		_, err := z.conn.Create(p, d, flags, acl)
		if err != nil && err != zk.ErrNodeExists {
			return err
		}
	}
	return nil
}
