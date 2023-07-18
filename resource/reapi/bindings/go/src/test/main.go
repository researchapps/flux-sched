/*****************************************************************************\
 * Copyright 2023 Lawrence Livermore National Security, LLC
 * (c.f. AUTHORS, NOTICE.LLNS, LICENSE)
 *
 * This file is part of the Flux resource manager framework.
 * For details, see https://github.com/flux-framework.
 *
 * SPDX-License-Identifier: LGPL-3.0
\*****************************************************************************/

package main

import (
	"flag"
	"fluxcli"
	"fmt"
	"io/ioutil"
)

func main() {
	ctx := fluxcli.NewReapiCli()
	jgfPtr := flag.String("jgf", "", "path to jgf")
	jobspecPtr := flag.String("jobspec", "", "path to jobspec")
	reserve := flag.Bool("reserve", false, "or else reserve?")
	flag.Parse()

	jgf, err := ioutil.ReadFile(*jgfPtr)
	if err != nil {
		fmt.Println("Error reading JGF file")
		return
	}
	err = fluxcli.ReapiCliInit(ctx, string(jgf), "{}")
	if err != nil {
		fmt.Printf("Error init ReapiCli: %v\n", err)
		return
	}
	fmt.Printf("Errors so far: %s\n", fluxcli.ReapiCliGetErrMsg(ctx))

	jobspec, err := ioutil.ReadFile(*jobspecPtr)
	if err != nil {
		fmt.Printf("Error reading jobspec file: %v\n", err)
		return
	}
	fmt.Printf("Jobspec:\n %s\n", jobspec)

	reserved, allocated, at, overhead, jobid, err := fluxcli.ReapiCliMatchAllocate(ctx, *reserve, string(jobspec))
	if err != nil {
		fmt.Printf("Error in ReapiCliMatchAllocate: %v\n", err)
		return
	}
	printOutput(reserved, allocated, at, jobid, err)
	reserved, allocated, at, overhead, jobid, err = fluxcli.ReapiCliMatchAllocate(ctx, *reserve, string(jobspec))
	fmt.Println("Errors so far: \n", fluxcli.ReapiCliGetErrMsg(ctx))

	if err != nil {
		fmt.Printf("Error in ReapiCliMatchAllocate: %v\n", err)
		return
	}
	printOutput(reserved, allocated, at, jobid, err)
	err = fluxcli.ReapiCliCancel(ctx, 1, false)
	if err != nil {
		fmt.Printf("Error in ReapiCliCancel: %v\n", err)
		return
	}
	fmt.Printf("Cancel output: %v\n", err)

	reserved, at, overhead, mode, err := fluxcli.ReapiCliInfo(ctx, 1)
	if err != nil {
		fmt.Printf("Error in ReapiCliInfo: %v\n", err)
		return
	}
	fmt.Printf("Info output jobid 1: %t, %d, %f, %s, %v\n", reserved, at, overhead, mode, err)

	reserved, at, overhead, mode, err = fluxcli.ReapiCliInfo(ctx, 2)
	if err != nil {
		fmt.Println("Error in ReapiCliInfo: %v\n", err)
		return
	}
	fmt.Printf("Info output jobid 2: %t, %d, %f, %v\n", reserved, at, overhead, err)

}

func printOutput(reserved bool, allocated string, at int64, jobid uint64, err error) {
	fmt.Println("\n\t----Match Allocate output---")
	fmt.Printf("jobid: %d\nreserved: %t\nallocated: %s\nat: %d\nerror: %v\n", jobid, reserved, allocated, at, err)
}
