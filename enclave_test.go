package nitriding

import (
	"testing"
)

var defaultCfg = Config{
	FQDN:          "example.com",
	ExtPort:       50000,
	IntPort:       50001,
	HostProxyPort: 1024,
	UseACME:       false,
	Debug:         false,
	FdCur:         1024,
	FdMax:         4096,
	WaitForApp:    true,
}

func createEnclave(cfg *Config) *Enclave {
	e, err := NewEnclave(cfg)
	if err != nil {
		panic(err)
	}
	return e
}

func TestValidateConfig(t *testing.T) {
	var err error
	var c Config

	if err = c.Validate(); err == nil {
		t.Fatalf("Validation of invalid config did not return an error.")
	}

	// Set one required field but leave others unset.
	c.FQDN = "example.com"
	if err = c.Validate(); err == nil {
		t.Fatalf("Validation of invalid config did not return an error.")
	}

	// Set the remaining required fields.
	c.ExtPort = 1
	c.IntPort = 1
	c.HostProxyPort = 1
	if err = c.Validate(); err != nil {
		t.Fatalf("Validation of valid config returned an error.")
	}
}

func TestGenSelfSignedCert(t *testing.T) {
	e := createEnclave(&defaultCfg)
	if err := e.genSelfSignedCert(); err != nil {
		t.Fatalf("Failed to create self-signed certificate: %s", err)
	}
}

func TestKeyMaterial(t *testing.T) {
	e := createEnclave(&defaultCfg)
	k := struct{ Foo string }{"foobar"}

	if _, err := e.KeyMaterial(); err != errNoKeyMaterial {
		t.Fatal("Expected error because we're trying to retrieve non-existing key material.")
	}

	e.SetKeyMaterial(k)
	r, err := e.KeyMaterial()
	if err != nil {
		t.Fatalf("Failed to retrieve key material: %s", err)
	}
	if r != k {
		t.Fatal("Retrieved key material is unexpected.")
	}
}
