package cago

type Store []byte

// create new store
func NewStore(data []byte) *Store {
	s := make(Store, 0)
	s = data
	return &s
}

// get size of store
// this function include size of headers data
// the headers length is 13 bytes
func (s *Store) SizeAll() int {
	return len(*s)
}

// clear all data in store
func (s *Store) Reset() int {
	return len(*s)
}

// set data in store
// if data exists this function will return an error
func (s *Store) Set(name string, value []byte) (int, error) {
	// v := make([]byte, len(name) + len(value))
	return 0, nil
}

// Cek data exist or not
// this function will return boolean
func (s *Store) Exist(name string) bool {

	return false
}

// put data is set or replace data if exist
// this function will return an error if the storage reaches the maximum memory limit specified in the configuration
func (s *Store) Put(name string, value []byte) error {

	return nil
}

// Size used for check size of size value
// if value not found this function will be returned -1
func (s *Store) Size(name string) int {
	return len(*s)
}

// check remaining age of value
func (s *Store) TimeLeft(name string) int {

	return len(*s)
}

// remove key and value in store
// return true if value exits and removed and return false if value not exist or didn't work removed
func (s *Store) Remove() bool {

	return false
}
