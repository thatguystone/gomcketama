package gomcketama

import (
	"fmt"
	"testing"

	mc "github.com/bradfitz/gomemcache/memcache"
)

func TestNodeKeys(t *testing.T) {
	for h, nKeys := range nodeKeys {
		_, keys, err := prepServer(h)
		if err != nil {
			t.Fatalf("prep server failed: %s", err)
		}

		if len(keys) != len(nKeys) {
			t.Fatalf("generated keys length==%d, expected %d",
				len(keys),
				len(nKeys))
		}

		for i := 0; i < len(keys); i++ {
			if keys[i] != nKeys[i] {
				t.Fatalf("invalid key: got %s, expected %s",
					keys[i],
					nKeys[i])
			}
		}
	}
}

func TestKVToNode(t *testing.T) {
	ks := &KetamaServerSelector{}

	for h, _ := range nodeKeys {
		ks.AddServer(h)
	}

	for k, h := range kvToNode {
		a, err := ks.PickServer(k)
		if err != nil {
			t.Fatalf("got error instead of server: %s", err)
		}

		ae, err := lookup(h)
		if err != nil {
			t.Fatalf("expected server lookup failed: %s", err)
		}

		if a.String() != ae.String() {
			t.Fatalf("for key %s, got server %s, expected %s: %s",
				k,
				a.String(),
				ae.String())
		}
	}
}

func TestMemcache(t *testing.T) {
	m := NewClient("localhost:11211", "localhost:11212")

	for i := 0; i < 100; i++ {
		err := m.Set(&mc.Item{
			Key:   fmt.Sprintf("%d", i),
			Value: []byte(fmt.Sprintf("%d", i+1)),
		})

		if err != nil {
			t.Fatalf("failed to set key: %s", err)
		}
	}

	for i := 0; i < 100; i++ {
		k := fmt.Sprintf("%d", i)
		it, err := m.Get(k)
		if err != nil {
			t.Fatalf("failed to get key %s: %s", k, err)
		}

		ev := fmt.Sprintf("%d", i+1)
		if string(it.Value) != ev {
			t.Fatalf("wrong value for %s, expected %s, got %s: %s",
				k,
				ev,
				string(it.Value))
		}
	}
}
