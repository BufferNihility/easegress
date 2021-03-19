package storage

import (
	"fmt"

	"github.com/megaease/easegateway/pkg/cluster"
	"github.com/megaease/easegateway/pkg/logger"
)

type (
	// Storage is the interface to contain storage APIs.
	Storage interface {
		Lock() error
		Unlock() error

		Get(key string) (*string, error)
		GetPrefix(prefix string) (map[string]string, error)

		Put(key, value string) error
		PutUnderLease(key, value string) error
		PutAndDelete(map[string]*string) error
		PutAndDeleteUnderLease(map[string]*string) error

		Delete(key string) error
		DeletePrefix(prefix string) error

		Watcher() (Watcher, error)
	}

	// Watcher is the interface to contain watcher APIs.
	Watcher interface {
		Watch(key string) (<-chan *string, error)
		WatchPrefix(prefix string) (<-chan map[string]*string, error)
		Close()
	}

	clusterStorage struct {
		name  string
		cls   cluster.Cluster
		mutex cluster.Mutex
	}
)

// New creates a storage.
func New(name string, cls cluster.Cluster) Storage {
	cs := &clusterStorage{
		name: name,
		cls:  cls,
	}

	err := cs.mutexGoReady()
	if err != nil {
		logger.Errorf(err.Error())

		return cs
	}

	return cs
}

func (cs *clusterStorage) mutexGoReady() error {
	if cs.mutex != nil {
		return nil
	}

	mutex, err := cs.cls.Mutex(cs.name)
	if err != nil {
		return fmt.Errorf("create mutex for %s failed: %v", cs.name, err)
	}

	cs.mutex = mutex

	return nil
}

func (cs *clusterStorage) Lock() error {
	err := cs.mutexGoReady()
	if err != nil {
		return err
	}

	return cs.mutex.Lock()
}

func (cs *clusterStorage) Unlock() error {
	err := cs.mutexGoReady()
	if err != nil {
		return err
	}

	return cs.mutex.Unlock()
}

func (cs *clusterStorage) Get(key string) (*string, error) {
	return cs.cls.Get(key)
}

func (cs *clusterStorage) GetPrefix(prefix string) (map[string]string, error) {
	return cs.cls.GetPrefix(prefix)
}

func (cs *clusterStorage) Put(key, value string) error {
	return cs.cls.Put(key, value)
}

func (cs *clusterStorage) PutUnderLease(key, value string) error {
	return cs.cls.PutUnderLease(key, value)
}

func (cs *clusterStorage) PutAndDelete(kvs map[string]*string) error {
	return cs.cls.PutAndDelete(kvs)
}

func (cs *clusterStorage) PutAndDeleteUnderLease(kvs map[string]*string) error {
	return cs.cls.PutAndDeleteUnderLease(kvs)
}

func (cs *clusterStorage) Delete(key string) error {
	return cs.cls.Delete(key)
}

func (cs *clusterStorage) DeletePrefix(prefix string) error {
	return cs.cls.DeletePrefix(prefix)
}

func (cs *clusterStorage) Watcher() (Watcher, error) {
	return cs.cls.Watcher()
}
