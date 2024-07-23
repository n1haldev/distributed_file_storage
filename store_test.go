package main

import (
	"bytes"
	"log"
	"testing"
	"fmt"

	"github.com/stretchr/testify/assert"
)

func TestPathTransformFunc(t *testing.T) {
	key := "samplekey"
	pathKey := CasPathTransformFunc(key)
	expectedOriginalKey := "9e3ddc4083856e058cf7ece035aab77749ea2ed6"
	expectedPathname := "9e3ddc4083/856e058cf7/ece035aab7/7749ea2ed6/"
	assert.Equal(t, pathKey.Pathname, expectedPathname)
	assert.Equal(t, pathKey.Filename, expectedOriginalKey)
}


func TestStore(t *testing.T) {
	s := newStore()

	defer tearDown(t, s)

	for i := 0;i < 50; i++ {
		key := fmt.Sprintf("samplekey_%d", i)
		data := bytes.NewReader([]byte("Some good bytes"))
		if err := s.writeStream(key, data); err != nil {
			t.Error(err)
		}
	
		res := assert.True(t, s.Has(key))
		if !res {
			t.Errorf("Key %s not found", key)
		}
		log.Println("Key found")
		
		r, err := s.Read(key)
		if err != nil {
			t.Errorf("Error reading file: %s", err)
		}
		
		buf := new(bytes.Buffer)
		n, err := buf.ReadFrom(r)
		if err != nil {
			t.Errorf("Error reading from buffer: %s", err)
		}
		
		log.Printf("Read %d bytes: %s", n, buf.String())
		assert.Equal(t, buf.String(), "Some good bytes")

		if err := s.Delete(key); err != nil {
			t.Error(err)
		}

		if ok := s.Has(key); ok {
			t.Errorf("Key %s not deleted", key)
		}
	}
}

func newStore() *Store {
	opts := StoreOpts {
		PathTransformFunc: CasPathTransformFunc,
	}
	return NewStore(opts)
}

func tearDown(t *testing.T, s *Store) {
	if err := s.Clear(); err != nil {
		t.Error(err)
	}
}