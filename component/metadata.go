package component

import (
	"reflect"
	"unsafe"
)

type MetadataType int

const (
	TStruct MetadataType = iota
	TString
	TSlice
)

type Metadata struct {
	sz int
	dt unsafe.Pointer
	tp MetadataType
}

func MakeString(dt unsafe.Pointer, sz int) string {
	var result string
	hdr := (*reflect.StringHeader)(unsafe.Pointer(&result))
	hdr.Data = uintptr(dt)
	hdr.Len = sz
	return result
}

func MakeSlice[T any](dt unsafe.Pointer, sz int) []T {
	var result []T
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&result))
	hdr.Data = uintptr(dt)
	hdr.Len = sz
	hdr.Cap = sz
	return result
}

func NewMetadataString(s string) *Metadata {
	md := &Metadata{}
	MetadataPutString(md, s)
	return md
}

func NewMetadataStruct[T any](s T) *Metadata {
	md := &Metadata{}
	MetadataPutStruct(md, s)
	return md
}

func NewMetadataSlice[T any](s []T) *Metadata {
	md := &Metadata{}
	MetadataPutSlice(md, s)
	return md
}

func MetadataPutString(self *Metadata, s string) {
	hdr := (*reflect.StringHeader)(unsafe.Pointer(&s))
	self.dt = unsafe.Pointer(hdr.Data)
	self.sz = hdr.Len
	self.tp = TString
}

func MetadataPutSlice[T any](self *Metadata, s []T) {
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	self.dt = unsafe.Pointer(hdr.Data)
	self.sz = hdr.Len
	self.tp = TSlice
}

func MetadataPutStruct[T any](self *Metadata, s T) {
	self.dt = unsafe.Pointer(&s)
	self.sz = 0
	self.tp = TStruct
}

func MetadataString(self *Metadata) string {
	return MakeString(self.dt, self.sz)
}

func MetadataStruct[T any](self *Metadata) T {
	return *(*T)(self.dt)
}

func MetadataSlice[T any](self *Metadata) []T {
	return MakeSlice[T](self.dt, self.sz)
}
