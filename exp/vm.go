package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
)

import (
	"github.com/sternix/commands/lib/sysctl"
)

type (
	VirtualMemory struct {
		Available     uint32
		Used          uint32
		PageCount     uint32
		FreeCount     uint32
		ActiveCount   uint32
		InactiveCount uint32
		CacheCount    uint32
		BufSpace      uint32
		WireCount     uint32
	}

	// load per 1,5,15 minutes
	LoadAvg struct {
		OneM     float32
		FiveM    float32
		FifteenM float32
	}
)

func getSysctlValue(name string) (uint32, error) {
	if val, err := sysctl.Uint32(name); err != nil {
		return 0, fmt.Errorf("cannot get sysctl value of %q: %v", name, err)
	} else {
		return val, nil
	}
}

func VirtualMemoryStat() (VirtualMemory, error) {
	var vm VirtualMemory
	vmsysctls := map[string]uint32{
		"vm.stats.vm.v_page_size":      0,
		"vm.stats.vm.v_page_count":     0,
		"vm.stats.vm.v_free_count":     0,
		"vm.stats.vm.v_active_count":   0,
		"vm.stats.vm.v_inactive_count": 0,
		"vm.stats.vm.v_cache_count":    0,
		"vfs.bufspace":                 0,
		"vm.stats.vm.v_wire_count":     0,
	}

	for k, _ := range vmsysctls {
		if val, err := getSysctlValue(k); err != nil {
			return vm, err
		} else {
			vmsysctls[k] = val
		}
	}
	pageSize := vmsysctls["vm.stats.vm.v_page_size"]
	vm.PageCount = vmsysctls["vm.stats.vm.v_page_count"] * pageSize
	vm.FreeCount = vmsysctls["vm.stats.vm.v_free_count"] * pageSize
	vm.ActiveCount = vmsysctls["vm.stats.vm.v_active_count"] * pageSize
	vm.InactiveCount = vmsysctls["vm.stats.vm.v_inactive_count"] * pageSize
	vm.CacheCount = vmsysctls["vm.stats.vm.v_cache_count"] * pageSize
	vm.BufSpace = vmsysctls["vfs.bufspace"]
	vm.WireCount = vmsysctls["vm.stats.vm.v_wire_count"] * pageSize
	vm.Available = vm.InactiveCount + vm.CacheCount + vm.FreeCount
	vm.Used = vm.PageCount - vm.Available
	return vm, nil
}

func LoadAvarage() (LoadAvg, error) {
	var loadAvg LoadAvg

	la := struct {
		LdAvg  [4]uint32 // normally fixed point number fixpt_t[3]
		FScale int64
	}{}

	loadInfo, err := sysctl.Raw("vm.loadavg")
	if err != nil {
		return loadAvg, fmt.Errorf("cannot get sysctl vm.loadavg: %v", err)
	}

	br := bytes.NewReader(loadInfo)
	binary.Read(br, binary.LittleEndian, &la)

	fs := float32(la.FScale)
	loadAvg.OneM = float32(la.LdAvg[0]) / fs
	loadAvg.FiveM = float32(la.LdAvg[1]) / fs
	loadAvg.FifteenM = float32(la.LdAvg[2]) / fs

	return loadAvg, nil
}

func main() {
	if vm, err := VirtualMemoryStat(); err != nil {
		log.Fatalln(err)
	} else {
		fmt.Printf("%+v\n", vm)
	}

	if la, err := LoadAvarage(); err != nil {
		log.Fatal(err)
	} else {
		// same as sysctl vm.loadavg
		fmt.Printf("%.2f %.2f %.2f\n", la.OneM, la.FiveM, la.FifteenM)
	}
}
