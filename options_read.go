package gorocksdb

// #include <stdlib.h>
// #include "rocksdb/c.h"
import "C"
import "unsafe"

// ReadTier controls fetching of data during a read request.
// An application can issue a read request (via Get/Iterators) and specify
// if that read should process data that ALREADY resides on a specified cache
// level. For example, if an application specifies BlockCacheTier then the
// Get call will process data that is already processed in the memtable or
// the block cache. It will not page in data from the OS cache or data that
// resides in storage.
type ReadTier uint

const (
	// ReadAllTier reads data in memtable, block cache, OS cache or storage.
	ReadAllTier = ReadTier(0)
	// BlockCacheTier reads data in memtable or block cache.
	BlockCacheTier = ReadTier(1)
)

type IterBound struct {
	bound *C.char
	size  C.size_t
}

func NewIterBound(v []byte) *IterBound {
	cv := C.CBytes(v)
	return &IterBound{bound: (*C.char)(cv), size: C.size_t(len(v))}
}

func (ib *IterBound) Destroy() {
	C.free(unsafe.Pointer(ib.bound))
}

// ReadOptions represent all of the available options when reading from a
// database.
type ReadOptions struct {
	c *C.rocksdb_readoptions_t
}

// NewDefaultReadOptions creates a default ReadOptions object.
func NewDefaultReadOptions() *ReadOptions {
	return NewNativeReadOptions(C.rocksdb_readoptions_create())
}

// NewNativeReadOptions creates a ReadOptions object.
func NewNativeReadOptions(c *C.rocksdb_readoptions_t) *ReadOptions {
	return &ReadOptions{c}
}

// UnsafeGetReadOptions returns the underlying c read options object.
func (opts *ReadOptions) UnsafeGetReadOptions() unsafe.Pointer {
	return unsafe.Pointer(opts.c)
}

// SetVerifyChecksums speciy if all data read from underlying storage will be
// verified against corresponding checksums.
// Default: false
func (opts *ReadOptions) SetVerifyChecksums(value bool) {
	C.rocksdb_readoptions_set_verify_checksums(opts.c, boolToChar(value))
}

// SetFillCache specify whether the "data block"/"index block"/"filter block"
// read for this iteration should be cached in memory?
// Callers may wish to set this field to false for bulk scans.
// Default: true
func (opts *ReadOptions) SetFillCache(value bool) {
	C.rocksdb_readoptions_set_fill_cache(opts.c, boolToChar(value))
}

func (opts *ReadOptions) SetIterUpperBound(value *IterBound) {
	C.rocksdb_readoptions_set_iterate_upper_bound(opts.c, value.bound, value.size)
}

func (opts *ReadOptions) SetIterLowerBound(value *IterBound) {
	C.rocksdb_readoptions_set_iterate_lower_bound(opts.c, value.bound, value.size)
}

func (opts *ReadOptions) SetPrefixSameAsStart(value bool) {
	C.rocksdb_readoptions_set_prefix_same_as_start(opts.c, boolToChar(value))
}

func (opts *ReadOptions) SetIgnoreRangeDeletions(value bool) {
	C.rocksdb_readoptions_set_ignore_range_deletions(opts.c, boolToChar(value))
}

// SetSnapshot sets the snapshot which should be used for the read.
// The snapshot must belong to the DB that is being read and must
// not have been released.
// Default: nil
func (opts *ReadOptions) SetSnapshot(snap *Snapshot) {
	if snap == nil {
		C.rocksdb_readoptions_set_snapshot(opts.c, nil)
		return
	}
	C.rocksdb_readoptions_set_snapshot(opts.c, snap.c)
}

// SetReadTier specify if this read request should process data that ALREADY
// resides on a particular cache. If the required data is not
// found at the specified cache, then Status::Incomplete is returned.
// Default: ReadAllTier
func (opts *ReadOptions) SetReadTier(value ReadTier) {
	C.rocksdb_readoptions_set_read_tier(opts.c, C.int(value))
}

// SetTailing specify if to create a tailing iterator.
// A special iterator that has a view of the complete database
// (i.e. it can also be used to read newly added data) and
// is optimized for sequential reads. It will return records
// that were inserted into the database after the creation of the iterator.
// Default: false
func (opts *ReadOptions) SetTailing(value bool) {
	C.rocksdb_readoptions_set_tailing(opts.c, boolToChar(value))
}

// Destroy deallocates the ReadOptions object.
func (opts *ReadOptions) Destroy() {
	C.rocksdb_readoptions_destroy(opts.c)
	opts.c = nil
}
