package main

import (
	"bytes"
	"sync"
)

// enclaveKeys holds key material for nitriding itself (the HTTPS certificate)
// and for the enclave application (whatever the application wants to "store"
// in nitriding).  These keys are meant to be managed by a leader enclave and --
// if horizontal scaling is required -- synced to worker enclaves.  The struct
// implements getters and setters that allow for thread-safe setting and getting
// of members.
type enclaveKeys struct {
	sync.RWMutex
	NitridingKey  []byte `json:"nitriding_key"`
	NitridingCert []byte `json:"nitriding_cert"`
	AppKeys       []byte `json:"app_keys"`
}

func (e1 *enclaveKeys) equal(e2 *enclaveKeys) bool {
	e1.RLock()
	e2.RLock()
	defer e1.RUnlock()
	defer e2.RUnlock()

	return bytes.Equal(e1.NitridingCert, e2.NitridingCert) &&
		bytes.Equal(e1.NitridingKey, e2.NitridingKey) &&
		bytes.Equal(e1.AppKeys, e2.AppKeys)
}

func (e *enclaveKeys) setAppKeys(appKeys []byte) {
	e.Lock()
	defer e.Unlock()

	e.AppKeys = appKeys
}

func (e *enclaveKeys) setNitridingKeys(key, cert []byte) {
	e.Lock()
	defer e.Unlock()

	e.NitridingKey = key
	e.NitridingCert = cert
}

func (e *enclaveKeys) set(newKeys *enclaveKeys) {
	e.Lock()
	defer e.Unlock()

	e.NitridingKey = newKeys.NitridingKey
	e.NitridingCert = newKeys.NitridingCert
	e.AppKeys = newKeys.AppKeys
}

func (e *enclaveKeys) get() *enclaveKeys {
	e.RLock()
	defer e.RUnlock()

	return &enclaveKeys{
		NitridingKey:  e.NitridingKey,
		NitridingCert: e.NitridingCert,
		AppKeys:       e.AppKeys,
	}
}

func (e *enclaveKeys) getAppKeys() []byte {
	e.RLock()
	defer e.RUnlock()

	return e.AppKeys
}
