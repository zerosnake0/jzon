package jzon

import (
	"log"
	"runtime"
	"strconv"
	"strings"
)

var (
	goVersion = newGoVersionInfo(runtime.Version())
)

func init() {
	log.Println("the current go version is:", runtime.Version())
}

type goVersionInfo struct {
	Major int
	Minor int
	Build int
}

func newGoVersionInfo(v string) (gv goVersionInfo) {
	if !strings.HasPrefix(v, "go") {
		return
	}
	arr := strings.Split(v[2:], ".")
	if len(arr) != 3 {
		return
	}
	major, err := strconv.Atoi(arr[0])
	if err != nil {
		return
	}
	minor, err := strconv.Atoi(arr[1])
	if err != nil {
		return
	}
	build, err := strconv.Atoi(arr[2])
	if err != nil {
		return
	}
	gv.Major = major
	gv.Minor = minor
	gv.Build = build
	return
}

func (gv goVersionInfo) LessEqual(v string) bool {
	other := newGoVersionInfo(v)
	if gv.Major > other.Major {
		return false
	}
	if gv.Major < other.Major {
		return true
	}
	if gv.Minor > other.Minor {
		return false
	}
	if gv.Minor < other.Minor {
		return true
	}
	return gv.Build <= other.Build
}
