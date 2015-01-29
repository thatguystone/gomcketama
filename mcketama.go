// Package gomcketama implements a ServerSelector for gomemcache that provides
// ketama hashing that's compatible with SpyMemcached's ketama hashing.
package gomcketama

import (
	"crypto/md5"
	"fmt"
	"net"
	"sort"
	"strings"
	"sync"

	mc "github.com/bradfitz/gomemcache/memcache"
)

const (
	numReps = 40
)

type info struct {
	hash uint32
	addr net.Addr
}

type infoSlice []info

// KetamaServerSelector implements gomemcache's ServerSelector
type KetamaServerSelector struct {
	lock  sync.RWMutex
	addrs []net.Addr
	ring  infoSlice
}

func (s infoSlice) Len() int {
	return len(s)
}

func (s infoSlice) Less(i, j int) bool {
	return s[i].hash < s[j].hash
}

func (s infoSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// AddServer adds a server to the ketama continuum.
func (ks *KetamaServerSelector) AddServer(server string) error {
	addr, keys, err := prepServer(server)
	if err != nil {
		return err
	}

	ks.lock.Lock()
	defer ks.lock.Unlock()

	for _, k := range keys {
		h := md5.Sum([]byte(k))

		for i := 0; i < 4; i++ {
			k := (uint32(h[3+i*4]&0xFF) << 24) |
				(uint32(h[2+i*4]&0xFF) << 16) |
				(uint32(h[1+i*4]&0xFF) << 8) |
				(uint32(h[0+i*4]&0xFF) << 0)
			ks.ring = append(ks.ring, info{
				hash: k,
				addr: addr,
			})
		}
	}

	sort.Sort(ks.ring)
	ks.addrs = append(ks.addrs, addr)

	return nil
}

// PickServer returns the server address that a given item should be written
// to.
func (ks *KetamaServerSelector) PickServer(key string) (net.Addr, error) {
	if len(ks.addrs) == 0 {
		return nil, mc.ErrNoServers
	}

	if len(ks.addrs) == 1 {
		return ks.addrs[0], nil
	}

	hm := md5.Sum([]byte(key))
	hk := (uint32(hm[3]&0xFF) << 24) |
		(uint32(hm[2]&0xFF) << 16) |
		(uint32(hm[1]&0xFF) << 8) |
		(uint32(hm[0]&0xFF) << 0)
	i := sort.Search(len(ks.ring), func(i int) bool {
		return ks.ring[i].hash >= hk
	})

	if i == len(ks.ring) {
		i = 0
	}

	return ks.ring[i].addr, nil
}

// Each loops through all registered servers, calling the given function.
func (ks *KetamaServerSelector) Each(f func(net.Addr) error) error {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	for _, a := range ks.addrs {
		if err := f(a); nil != err {
			return err
		}
	}
	return nil
}

// New creates a new memcache client with the given servers, using ketama as
// the ServerSelector. This functions exactly like gomemcache's New().
func New(server ...string) *mc.Client {
	ks := &KetamaServerSelector{}

	for _, s := range server {
		ks.AddServer(s)
	}

	return mc.NewFromSelector(ks)
}

func lookup(server string) (a net.Addr, err error) {
	if strings.Contains(server, "/") {
		a, err = net.ResolveUnixAddr("unix", server)
	} else {
		a, err = net.ResolveTCPAddr("tcp", server)
	}

	return
}

func prepServer(server string) (a net.Addr, keys []string, err error) {
	host, port, err := net.SplitHostPort(server)
	if err != nil {
		return
	}

	a, err = lookup(server)
	if err != nil {
		return
	}

	key := fmt.Sprintf("%s:%s", host, port)
	if !strings.HasPrefix(a.String(), host) {
		key = fmt.Sprintf("%s/%s", host, a.String())
	}

	for i := 0; i < numReps; i++ {
		keys = append(keys, fmt.Sprintf("%s-%d", key, i))
	}

	return
}
