// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import ()

// Frames may be used to get function/file/line information for a
// slice of PC values returned by Callers.
type Frames struct {
	callers []uintptr

	// If previous caller in iteration was a panic, then
	// ci.callers[0] is the address of the faulting instruction
	// instead of the return address of the call.
	wasPanic bool

	// Frames to return for subsequent calls to the Next method.
	// Used for non-Go frames.
	frames *[]Frame
}

// Frame is the information returned by Frames for each call frame.
type Frame struct {
	// Program counter for this frame; multiple frames may have
	// the same PC value.
	PC uintptr

	// Func for this frame; may be nil for non-Go code or fully
	// inlined functions.
	Func *Func

	// Function name, file name, and line number for this call frame.
	// May be the empty string or zero if not known.
	// If Func is not nil then Function == Func.Name().
	Function string
	File     string
	Line     int

	// Entry point for the function; may be zero if not known.
	// If Func is not nil then Entry == Func.Entry().
	Entry uintptr
}

// NOTE: Func does not expose the actual unexported fields, because we return *Func
// values to users, and we want to keep them from being able to overwrite the data
// with (say) *f = Func{}.
// All code operating on a *Func must call raw to get the *_func instead.

// A Func represents a Go function in the running binary.
type Func struct {
	opaque struct{} // unexported field to disallow conversions
}

// funcdata.h
const (
	_PCDATA_StackMapIndex       = 0
	_FUNCDATA_ArgsPointerMaps   = 0
	_FUNCDATA_LocalsPointerMaps = 1
	_ArgsSizeUnknown            = -0x80000000
)

// moduledata records information about the layout of the executable
// image. It is written by the linker. Any changes here must be
// matched changes to the code in cmd/internal/ld/symtab.go:symtab.
// moduledata is stored in read-only memory; none of the pointers here
// are visible to the garbage collector.
type moduledata [393]byte
/*
type moduledata struct {
	pclntable    []byte
	ftab         []functab
	filetab      []uint32
	findfunctab  uintptr
	minpc, maxpc uintptr

	text, etext           uintptr
	noptrdata, enoptrdata uintptr
	data, edata           uintptr
	bss, ebss             uintptr
	noptrbss, enoptrbss   uintptr
	end, gcdata, gcbss    uintptr
	types, etypes         uintptr

	textsectmap []textsect
	typelinks []int32 // offsets from types
	itablinks []*itab

	ptab []ptabEntry

	pluginpath string

	modulename   string
	modulehashes []modulehash

	gcdatamask, gcbssmask bitvector

	typemap map[typeOff]*_type // offset to *_rtype in previous module

	next *moduledata
}*/

// For each shared library a module links against, the linker creates an entry in the
// moduledata.modulehashes slice containing the name of the module, the abi hash seen
// at link time and a pointer to the runtime abi hash. These are checked in
// moduledataverify1 below.
type modulehash struct {
	modulename   string
	linktimehash string
	runtimehash  *string
}

var firstmoduledata moduledata  // linker symbol
var lastmoduledatap *moduledata // linker symbol

type functab struct {
	entry   uintptr
	funcoff uintptr
}

const minfunc = 16                 // minimum function size
const pcbucketsize = 256 * minfunc // size of bucket in the pc->func lookup table

// findfunctab is an array of these structures.
// Each bucket represents 4096 bytes of the text segment.
// Each subbucket represents 256 bytes of the text segment.
// To find a function given a pc, locate the bucket and subbucket for
// that pc. Add together the idx and subbucket value to obtain a
// function index. Then scan the functab array starting at that
// index to find the target function.
// This table uses 20 bytes for every 4096 bytes of code, or ~0.5% overhead.
type findfuncbucket struct {
	idx        uint32
	subbuckets [16]byte
}

const debugPcln = false

type pcvalueCache struct {
	entries [16]pcvalueCacheEnt
}

type pcvalueCacheEnt struct {
	// targetpc and off together are the key of this cache entry.
	targetpc uintptr
	off      int32
	// val is the value of this cached pcvalue entry.
	val int32
}

type stackmap struct {
	n        int32   // number of bitmaps
	nbit     int32   // number of bits in each bitmap
	bytedata [1]byte // bitmaps, each starting on a 32-bit boundary
}
