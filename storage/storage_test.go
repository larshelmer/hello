package storage

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var testPath = filepath.Join(os.TempDir(), "motd_test_storage.json")

func TestParseFile(t *testing.T) {
	//	t.SkipNow()
	motd := "a motd"
	want := data{[]string{motd}}
	json := []byte("{\"messages\":[\"" + motd + "\"]}")
	got, err := parseFile(json)
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}
	if len(want.Messages) != len((*got).Messages) {
		t.Errorf("parseFile(%q) == %q, want %q", json, *got, want)
	} else {
		for ix := range want.Messages {
			if want.Messages[ix] != got.Messages[ix] {
				t.Errorf("parseFile(%q) == %q, want %q", json, got, want)
			}
		}
	}
}

func TestParseFileFail(t *testing.T) {
	//	t.SkipNow()
	_, err := parseFile(nil)
	if err != nil {
		t.Log(err.Error())
	} else {
		t.Error("parseFile(nil) expected error")
	}
}

func TestRead(t *testing.T) {
	//	t.SkipNow()
	InitData(testPath)
	want := data{[]string{initialData}}
	defer os.Remove(testPath)
	got, err := Read()
	if err != nil {
		t.Fatal(err.Error())
	}
	for ix := range want.Messages {
		if want.Messages[ix] != (*got)[ix] {
			t.Errorf("Read() == %q, want %q", got, want)
		}
	}
}

func TestReadGarbage(t *testing.T) {
	//	t.SkipNow()
	json := []byte("asdf#!=\\")
	ioutil.WriteFile(testPath, json, 0)
	defer os.Remove(testPath)
	InitData(testPath)
	_, err := Read()
	if err == nil {
		t.Fatal("Read() expected error")
	}
}

func TestReadEmptyFile(t *testing.T) {
	InitData(testPath)
	os.Remove(testPath)
	f, _ := os.Create(testPath)
	defer os.Remove(testPath)
	f.Close()
	var want []string
	got, _ := Read()
	if len(*got) != len(want) {
		t.Errorf("Read() == %q, want %q", got, want)
	}
}

func TestMakeJSON(t *testing.T) {
	//	t.SkipNow()
	in := []string{"new motd"}
	want := []byte(fmt.Sprintf("{\"messages\":%q}", in))
	got, err := makeJSON(in)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !bytes.Equal(got, want) {
		t.Errorf("makeJSON(%q) == %q, want %q", in, got, want)
	}
}

func TestAddEmpty(t *testing.T) {
	//	t.SkipNow()
	err := Add("")
	if err == nil {
		t.Error("Add(\"\") expected error")
	}
}

func TestAdd(t *testing.T) {
	in := "new motd"
	InitData(testPath)
	os.Remove(testPath)
	f, _ := os.Create(testPath)
	defer os.Remove(testPath)
	f.Close()
	err := Add(in)
	if err != nil {
		t.Errorf("Add(%q) == %q, want nil", in, err.Error())
	}
}

func TestInitData(t *testing.T) {
	if _, err := os.Stat(testPath); err == nil {
		os.Remove(testPath)
	}
	err := InitData(testPath)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(testPath)
	want := []string{"quidquid Latine dictum sit altum videtur"}
	dat, _ := Read()
	if (*dat)[0] != want[0] {
		t.Errorf("Read() == %q, want %q", *dat, want)
	}
}

func TestInitDataGarbage(t *testing.T) {
	//	t.SkipNow()
	json := []byte("asdf#!=\\")
	ioutil.WriteFile(testPath, json, 0)
	defer os.Remove(testPath)
	err := InitData(testPath)
	if err == nil {
		t.Fatal("Read() expected error")
	}
}
