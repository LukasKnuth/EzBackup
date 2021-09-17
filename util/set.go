package util

import (
	"fmt"

	"github.com/LukasKnuth/EzBackup/k8s"
)

/// Go has no Set type, and we can't implement custom equality for a struct.
/// Therefor we implement custom equality as a unique string and map to the actual interface type.
type SpecificSet map[string]k8s.BlockingMountOwner

func MakeSet(capacity int) SpecificSet {
	return make(SpecificSet, capacity)
}

func (set SpecificSet) Put(entry k8s.BlockingMountOwner) {
	key := toKey(entry)
	set[key] = entry
}

func (set SpecificSet) PutAll(entries []k8s.BlockingMountOwner) {
	for _, bmo := range entries {
		set.Put(bmo)
	}
}

func toKey(entry k8s.BlockingMountOwner) string {
	return fmt.Sprintf("key_%s_%s", entry.Kind(), entry.Name())
}