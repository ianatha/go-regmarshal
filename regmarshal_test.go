package regmarshal

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"golang.org/x/sys/windows/registry"
)

type Datum struct {
	Text   string
	Number int
	Binary []byte
}

func exampleDatum() *Datum {
	return &Datum{
		Text:   "hello world",
		Number: 314,
		Binary: []byte{1, 2, 3, 4, 5},
	}
}

const RegistryRootKey = registry.CURRENT_USER
const RegistryRootPath = "SOFTWARE\\AthaIO\\Regmarshal\\"

func TestMarshal(t *testing.T) {
	TestRegistryPAth := fmt.Sprintf("%s%s_%d", RegistryRootPath, t.Name(), time.Now().Unix())

	datum := exampleDatum()

	_, _, err := registry.CreateKey(RegistryRootKey, TestRegistryPAth, registry.ALL_ACCESS)
	if err != nil {
		t.Error(err)
	}

	err = Marshal(datum, RegistryRootKey, TestRegistryPAth)
	if err != nil {
		t.Error(err)
	}

	key, err := registry.OpenKey(RegistryRootKey, TestRegistryPAth, registry.ALL_ACCESS)
	if err != nil {
		t.Error(err)
	}

	textVal, _, err := key.GetStringValue("Text")
	if err != nil {
		t.Error(err)
	}
	if textVal != datum.Text {
		t.Errorf("registry: %s; want: %s", textVal, datum.Text)
	}

	numVal, _, err := key.GetIntegerValue("Number")
	if err != nil {
		t.Error(err)
	}
	if numVal != uint64(datum.Number) {
		t.Errorf("registry: %d; want: %d", numVal, datum.Number)
	}

	binVal, _, err := key.GetBinaryValue("Binary")
	if err != nil {
		t.Error(err)
	}
	if bytes.Compare(binVal, datum.Binary) != 0 {
		t.Errorf("registry: %v; want: %v", binVal, datum.Binary)
	}

	key.Close()

	err = registry.DeleteKey(RegistryRootKey, TestRegistryPAth)
	if err != nil {
		t.Error(err)
	}
}

func TestUnmarshal(t *testing.T) {
	TestRegistryPath := fmt.Sprintf("%s%s_%d", RegistryRootPath, t.Name(), time.Now().Unix())

	datum := exampleDatum()

	_, _, err := registry.CreateKey(RegistryRootKey, TestRegistryPath, registry.ALL_ACCESS)
	if err != nil {
		t.Error(err)
	}

	key, err := registry.OpenKey(RegistryRootKey, TestRegistryPath, registry.ALL_ACCESS)
	if err != nil {
		t.Error(err)
	}

	key.SetStringValue("Text", datum.Text)
	key.SetQWordValue("Number", uint64(datum.Number))
	key.SetBinaryValue("Binary", []byte{1, 2, 3, 4, 5})

	readDatum := Datum{}
	err = Unmarshal(RegistryRootKey, TestRegistryPath, &readDatum)
	if err != nil {
		t.Error(err)
	}

	if readDatum.Text != datum.Text {
		t.Errorf("registry: %s; want: %s", readDatum.Text, datum.Text)
	}

	if readDatum.Number != datum.Number {
		t.Errorf("registry: %d; want: %d", readDatum.Number, datum.Number)
	}

	if bytes.Compare(readDatum.Binary, datum.Binary) != 0 {
		t.Errorf("registry: %v; want: %v", readDatum.Binary, datum.Binary)
	}

	key.Close()

	err = registry.DeleteKey(RegistryRootKey, TestRegistryPath)
	if err != nil {
		t.Error(err)
	}
}


func TestUnmarshalMissingValues(t *testing.T) {
	TestRegistryPath := fmt.Sprintf("%s%s_%d", RegistryRootPath, t.Name(), time.Now().Unix())

	datum := exampleDatum()
	datum.Text = ""
	datum.Number = 0
	datum.Binary = []byte{}

	_, _, err := registry.CreateKey(RegistryRootKey, TestRegistryPath, registry.ALL_ACCESS)
	if err != nil {
		t.Error(err)
	}

	key, err := registry.OpenKey(RegistryRootKey, TestRegistryPath, registry.ALL_ACCESS)
	if err != nil {
		t.Error(err)
	}

	readDatum := Datum{}
	err = Unmarshal(RegistryRootKey, TestRegistryPath, &readDatum)
	if err != nil {
		t.Error(err)
	}

	if readDatum.Text != datum.Text {
		t.Errorf("registry: %s; want: %s", readDatum.Text, datum.Text)
	}

	if readDatum.Number != datum.Number {
		t.Errorf("registry: %d; want: %d", readDatum.Number, datum.Number)
	}

	if bytes.Compare(readDatum.Binary, datum.Binary) != 0 {
		t.Errorf("registry: %v; want: %v", readDatum.Binary, datum.Binary)
	}

	key.Close()

	err = registry.DeleteKey(RegistryRootKey, TestRegistryPath)
	if err != nil {
		t.Error(err)
	}
}
