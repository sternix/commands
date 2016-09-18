package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"runtime"
)

import (
	"github.com/sternix/commands/lib/sysctl"
)

// /usr/include/sys/resource.h
const (
	CP_USER   = 0
	CP_NICE   = 1
	CP_SYS    = 2
	CP_INTR   = 3
	CP_IDLE   = 4
	CPUSTATES = 5
)

type CpuStat struct {
	Name string
	User int64
	Nice int64
	Sys  int64
	Intr int64
	Idle int64
}

func main() {
	lp, err := LastPID()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Last Pid:%d\n", lp)

	cpus, err := CpuTimes()
	if err != nil {
		log.Fatal(err)
	}

	for _, c := range cpus {
		fmt.Printf("%+v\n", c)
	}
}

func CpuTimes() ([]CpuStat, error) {
	var cpuStats []CpuStat
	ncpu := runtime.NumCPU()
	times := make([]int64, ncpu*CPUSTATES)

	ret, err := sysctl.Raw("kern.cp_times")
	if err != nil {
		return cpuStats, err
	}

	br := bytes.NewReader(ret)
	binary.Read(br, binary.LittleEndian, &times)

	for i := 0; i < ncpu; i++ {
		offset := i * CPUSTATES

		cpuStat := CpuStat{
			Name: fmt.Sprintf("CPU %d", i),
			User: times[CP_USER+offset],
			Nice: times[CP_NICE+offset],
			Sys:  times[CP_SYS+offset],
			Intr: times[CP_INTR+offset],
			Idle: times[CP_IDLE+offset],
		}
		cpuStats = append(cpuStats, cpuStat)
	}
	return cpuStats, nil
}

// NOT TESTED
// from 
// https://github.com/freebsd/freebsd/blob/master/contrib/top/utils.c
func Percentages(cnt int, out []int64, newv []int64, old []int64, diffs []int64) int64 {
	var (
		total_change int64
		dp           []int64 = diffs
		half_total   int64
		change       int64
	)

	/* calculate changes for each state and the overall change */
	for i := 0; i < cnt; i++ {
		if change = newv[i] - old[i]; change < 0 {
			/* this only happens when the counter wraps */
			change = newv[i] - old[i]
		}
		dp[i] = change
		total_change += change
		old[i] = newv[i]
	}

	/* avoid divide by zero potential */
	if total_change == 0 {
		total_change = 1
	}

	/* calculate percentages based on overall change, rounding up */
	half_total = total_change / 2

	for i := 0; i < cnt; i++ {
		out[i] = ((diffs[i]*1000 + half_total) / total_change)
	}

	/* return the total in case the caller wants to use it */
	return total_change
}

func LastPID() (uint32, error) {
	lp, err := sysctl.Uint32("kern.lastpid")
	if err != nil {
		return 0, err
	}

	return lp, nil
}
