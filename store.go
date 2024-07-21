package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	// "io/fs"
	"errors"
	"log"
	"os"
	"strings"
)

const defaultRoot = "fileNetwork"

func CasPathTransformFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])

	blockSize := 10
	sliceLen := len(hashStr) / blockSize

	paths := make([]string, sliceLen+1)

	for i := 0; i < sliceLen; i++ {
		from, to := i*blockSize, (i*blockSize)+blockSize
		paths[i] = hashStr[from:to]
	}

	return PathKey{
		Pathname: strings.Join(paths, "/"),
		Filename: hashStr,
	}
}

type PathTransformFunc func(string) PathKey

type PathKey struct {
	Pathname string
	Filename string
}

func (p PathKey) FullPath() string {
	return fmt.Sprintf("%s%s", p.Pathname, p.Filename)
}

type StoreOpts struct {
	// Root of the folder containing all the files of the system
	Root              string
	PathTransformFunc PathTransformFunc
}

var DefaultPathTransformFunc = func(key string) PathKey {
	return PathKey{
		Pathname: key,
		Filename: key,
	}
}

type Store struct {
	StoreOpts
}

func (p PathKey) FirstPathName() string {
	paths := strings.Split(p.Pathname, "/")
	if len(paths) == 0 {
		return ""
	}
	return paths[0]
}

func NewStore(opts StoreOpts) *Store {
	if opts.PathTransformFunc == nil {
		opts.PathTransformFunc = DefaultPathTransformFunc
	}
	if len(opts.Root) == 0 {
		opts.Root = defaultRoot
	}
	return &Store{
		StoreOpts: opts,
	}
}

func (s *Store) Has(key string) bool {
	pathKey := s.PathTransformFunc(key)
	fullPathWithRoot := s.Root + "/" + pathKey.FullPath()
	
	_, err := os.Stat(fullPathWithRoot)
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func (s *Store) Delete(key string) error {
	pathKey := s.PathTransformFunc(key)
	pathKeyWithRoot := s.Root + "/" + pathKey.FirstPathName()

	defer func() {
		log.Printf("Deleted %s", pathKeyWithRoot)
	}()

	return os.RemoveAll(pathKeyWithRoot)
}

func (s *Store) Read(key string) (io.Reader, error) {
	f, err := s.readStream(key)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, f)

	return buf, err
}

func (s *Store) readStream(key string) (io.ReadCloser, error) {
	pathKey := s.PathTransformFunc(key)
	pathKeyWithRoot := s.Root + "/" + pathKey.FullPath()
	return os.Open(pathKeyWithRoot)
}

func (s *Store) writeStream(key string, r io.Reader) error {
	pathKey := s.PathTransformFunc(key)
	pathNameWithRoot := s.Root + "/" + pathKey.Pathname

	if err := os.MkdirAll(pathNameWithRoot, os.ModePerm); err != nil {
		return err
	}

	fullpath := pathKey.FullPath()
	fullPathWithRoot := s.Root + "/" + fullpath

	f, err := os.Create(fullPathWithRoot)
	if err != nil {
		return err
	}

	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}

	log.Printf("Wrote %d bytes to %s in Disk", n, fullPathWithRoot)

	return nil
}
