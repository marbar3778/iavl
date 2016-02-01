package test

import (
	"bytes"
	"testing"

	. "github.com/tendermint/go-common"
	"github.com/tendermint/merkleeyes/app"
	mecli "github.com/tendermint/merkleeyes/client/golang"
	"github.com/tendermint/tmsp/server"
)

func TestClient(t *testing.T) {

	addr := "tcp://127.0.0.1:46659"

	// Start the listener
	mApp := app.NewMerkleEyesApp()
	ln, err := server.StartListener(addr, mApp)
	if err != nil {
		t.Fatal(err.Error())
		return
	}
	defer ln.Close()

	// Create client
	conn, err := Connect(addr)
	if err != nil {
		t.Fatal(err.Error())
		return
	}
	cli := mecli.NewMEClient(conn, 100)
	defer conn.Close()

	// Empty
	getHash(t, cli, "")
	get(t, cli, "foo", "", "")
	get(t, cli, "bar", "", "")
	// Set foo=FOO
	set(t, cli, "foo", "FOO")
	getHash(t, cli, "68DECA470D80183B5E979D167E3DD0956631A952")
	get(t, cli, "foo", "FOO", "")
	get(t, cli, "foa", "", "")
	get(t, cli, "foz", "", "")
	rem(t, cli, "foo")
	// Empty
	get(t, cli, "foo", "", "")
	getHash(t, cli, "")
	// Set foo1, foo2, foo3...
	set(t, cli, "foo1", "1")
	set(t, cli, "foo2", "2")
	set(t, cli, "foo3", "3")
	set(t, cli, "foo1", "4")
	get(t, cli, "foo1", "4", "")
	get(t, cli, "foo2", "2", "")
	get(t, cli, "foo3", "3", "")
	rem(t, cli, "foo3")
	rem(t, cli, "foo2")
	rem(t, cli, "foo1")
	// Empty
	getHash(t, cli, "")

}

func get(t *testing.T, cli *mecli.MEClient, key string, value string, err string) {
	_value, _err := cli.GetSync([]byte(key))
	if !bytes.Equal([]byte(value), _value) {
		t.Errorf("Expected value 0x%X (%v) but got 0x%X", []byte(value), value, _value)
	}
	if _err == nil {
		if err != "" {
			t.Errorf("Expected error %v but got no error", err)
		}
	} else {
		if err == "" {
			t.Errorf("Expected no error but got error %v", _err.Error())
		}
	}
}

func set(t *testing.T, cli *mecli.MEClient, key string, value string) {
	cli.SetSync([]byte(key), []byte(value))
}

func rem(t *testing.T, cli *mecli.MEClient, key string) {
	cli.RemSync([]byte(key))
}

func getHash(t *testing.T, cli *mecli.MEClient, hash string) {
	_hash, _, err := cli.GetHashSync()
	if err != nil {
		t.Error("Unexpected error getting hash", err.Error())
	}
	if hash != Fmt("%X", _hash) {
		t.Errorf("Expected hash 0x%v but got 0x%X", hash, _hash)
	}
}
