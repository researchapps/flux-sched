/*****************************************************************************\
 * Copyright 2023 Lawrence Livermore National Security, LLC
 * (c.f. AUTHORS, NOTICE.LLNS, LICENSE)
 *
 * This file is part of the Flux resource manager framework.
 * For details, see https://github.com/flux-framework.
 *
 * SPDX-License-Identifier: LGPL-3.0
\*****************************************************************************/

package fluxcli

/*
#include "reapi_cli.h"
*/
import "C"
import (
	"fmt"
	"unsafe"
)

type (
	ReapiCtx C.struct_reapi_cli_ctx_t
)

// NewReapiCli creates a new resource API client
// reapi_cli_ctx_t *reapi_cli_new ();
func NewReapiCli() *ReapiCtx {
	return (*ReapiCtx)(C.reapi_cli_new())
}

// // ReapiCliGetNode retruns a resource node
// This function is currently not being used
// func ReapiCliGetNode(ctx *ReapiCtx) string {
// 	return C.GoString(C.reapi_cli_get_node((*C.struct_reapi_cli_ctx)(ctx)))
// }

// Given an integer return code, convert to go error
// Also provide a meaningful string to the developer user
func retvalToError(code int, message string) error {
	if code == 0 {
		return nil
	}
	return fmt.Errorf(message+" %d", code)
}

// ReapiCliDestroy destroys a resource API context
// void reapi_cli_destroy (reapi_cli_ctx_t *ctx);
func ReapiCliDestroy(ctx *ReapiCtx) {
	C.reapi_cli_destroy((*C.struct_reapi_cli_ctx)(ctx))
}

// ReapiCliInit initializes a new resource API context
// int reapi_cli_initialize (reapi_cli_ctx_t *ctx, const char *jgf);
func ReapiCliInit(ctx *ReapiCtx, jgf string, options string) (err error) {

	jobgraph := C.CString(jgf)
	defer C.free(unsafe.Pointer(jobgraph))

	opts := C.CString(options)
	defer C.free(unsafe.Pointer(opts))
	fluxerr := (int)(
		C.reapi_cli_initialize(
			(*C.struct_reapi_cli_ctx)(ctx), jobgraph, (opts),
		),
	)

	return retvalToError(fluxerr, "issue initializing resource api client")
}

// ReapiCliMatchAllocate matches a jobspec to the best resources, either
// allocating or reserved them. The best resources are determined by the
// match policy.
//
//	\param ctx       reapi_cli_ctx_t context object
//	\param or else_reserve
//	                Boolean: if false, only allocate; otherwise, first try
//	                 to allocate and if that fails, reserve.
//	\param jobspec   jobspec string.
//	\param jobid     jobid of the uint64_t type.
//	\param reserved  Boolean into which to return true if this job has been
//	                 reserved instead of allocated.
//	\param R         String into which to return the resource set either
//	                 allocated or reserved.
//	\param at        If allocated, 0 is returned; if reserved, actual time
//	                 at which the job is reserved.
//	\param ov        Double into which to return performance overhead
//	                 in terms of elapse time needed to complete
//	                 the match operation.
//	\return          0 on success; -1 on error.
//
// int reapi_module_match_allocate (reapi_module_ctx_t *ctx, bool orelse_reserve,
//
//	const char *jobspec, const uint64_t jobid,
//	bool *reserved,
//	char **R, int64_t *at, double *ov);
func ReapiCliMatchAllocate(
	ctx *ReapiCtx,
	orelse_reserve bool,
	jobspec string,
) (reserved bool, allocated string, at int64, overhead float64, jobid uint64, err error) {
	var r = C.CString("")
	defer C.free(unsafe.Pointer(r))

	spec := C.CString(jobspec)
	defer C.free(unsafe.Pointer(spec))

	fluxerr := (int)(C.reapi_cli_match_allocate((*C.struct_reapi_cli_ctx)(ctx),
		(C.bool)(orelse_reserve),
		spec,
		(*C.ulong)(&jobid),
		(*C.bool)(&reserved),
		&r,
		(*C.long)(&at),
		(*C.double)(&overhead)))

	allocated = C.GoString(r)

	err = retvalToError(fluxerr, "issue resource api client matching allocate")
	return reserved, allocated, at, overhead, jobid, err

}

// ReapiCliUpdateAllocate updates the resource state with R.
//
//	\param ctx       reapi_cli_ctx_t context object
//	\param jobid     jobid of the uint64_t type.
//	\param R         R string
//	\param at        return the scheduled time
//	\param ov        return the performance overhead
//	                 in terms of elapse time needed to complete
//	                 the match operation.
//	\param R_out     return the updated R string.
//	\return          0 on success; -1 on error.
//
// int reapi_cli_update_allocate (reapi_cli_ctx_t *ctx,
//
//	const uint64_t jobid, const char *R, int64_t *at,
//	double *ov, const char **R_out);
func ReapiCliUpdateAllocate(ctx *ReapiCtx, jobid int, r string) (at int64, overhead float64, r_out string, err error) {
	var tmp_rout = C.CString("")
	defer C.free(unsafe.Pointer(tmp_rout))

	var resource := C.CString(r)
	defer C.free(unsafe.Pointer(resource))

	fluxerr := (int)(C.reapi_cli_update_allocate((*C.struct_reapi_cli_ctx)(ctx),
		(C.ulong)(jobid),
		resource,
		(*C.long)(&at),
		(*C.double)(&overhead),
		&tmp_rout))

	r_out = C.GoString(tmp_rout)

	err = retvalToError(fluxerr, "issue resource api client updating allocate")
	return at, overhead, r_out, err
}

// ReapiCliCancel cancels the allocation or reservation corresponding to jobid.
//
//	\param ctx       reapi_cli_ctx_t context object
//	\param jobid     jobid of the uint64_t type.
//	\param noent_ok  don't return an error on nonexistent jobid
//	\return          0 on success; -1 on error.
//
// int reapi_cli_cancel (reapi_cli_ctx_t *ctx,
//
//	const uint64_t jobid, bool noent_ok);
func ReapiCliCancel(ctx *ReapiCtx, jobid int64, noent_ok bool) (err error) {
	fluxerr := (int)(C.reapi_cli_cancel((*C.struct_reapi_cli_ctx)(ctx),
		(C.ulong)(jobid),
		(C.bool)(noent_ok)))
	return retvalToError(fluxerr, "issue resource api client cancel")
}

// ReapiCliInfo gets the information on the allocation or reservation corresponding
//
//	to jobid.
//	\param ctx       reapi_cli_ctx_t context object
//	\param jobid     const jobid of the uint64_t type.
//	\param reserved  Boolean into which to return true if this job has been
//	                 reserved instead of allocated.
//	\param at        If allocated, 0 is returned; if reserved, actual time
//	                 at which the job is reserved.
//	\param ov        Double into which to return performance overhead
//	                 in terms of elapse time needed to complete
//	                 the match operation.
//	\return          0 on success; -1 on error.
//
// int reapi_cli_info (reapi_cli_ctx_t *ctx, const uint64_t jobid,
//
//	bool *reserved, int64_t *at, double *ov);
func ReapiCliInfo(ctx *ReapiCtx, jobid int64) (reserved bool, at int64, overhead float64, mode string, err error) {
	var tmp_mode = C.CString("")
	defer C.free(unsafe.Pointer(tmp_mode))

	fluxerr := (int)(C.reapi_cli_info((*C.struct_reapi_cli_ctx)(ctx),
		(C.ulong)(jobid),
		(&tmp_mode),
		(*C.bool)(&reserved),
		(*C.long)(&at),
		(*C.double)(&overhead)))

	err = retvalToError(fluxerr, "issue resource api client info")
	return reserved, at, overhead, C.GoString(tmp_mode), err
}

// ReapiCliStat gets the performance information about the resource infrastructure.
//
//	\param ctx       reapi_cli_ctx_t context object
//	\param V         Number of resource vertices
//	\param E         Number of edges
//	\param J         Number of jobs
//	\param load      Graph load time
//	\param min       Min match time
//	\param max       Max match time
//	\param avg       Avg match time
//	\return          0 on success; -1 on error.
//
// int reapi_cli_stat (reapi_cli_ctx_t *ctx, int64_t *V, int64_t *E,
//
//	int64_t *J, double *load,
//	double *min, double *max, double *avg);
func ReapiCliStat(ctx *ReapiCtx) (v int64, e int64,
	jobs int64, load float64, min float64, max float64, avg float64, err error) {
	fluxerr := (int)(C.reapi_cli_stat((*C.struct_reapi_cli_ctx)(ctx),
		(*C.long)(&v),
		(*C.long)(&e),
		(*C.long)(&jobs),
		(*C.double)(&load),
		(*C.double)(&min),
		(*C.double)(&max),
		(*C.double)(&avg)))

	err = retvalToError(fluxerr, "issue resource api client stat")
	return v, e, jobs, load, min, max, avg, err
}

// ReapiCliGetErrMsg returns a string error message from the resource api
func ReapiCliGetErrMsg(ctx *ReapiCtx) string {
	errmsg := C.reapi_cli_get_err_msg((*C.struct_reapi_cli_ctx)(ctx))
	return C.GoString(errmsg)
}

// ReapiCliGetErrMsg clears error messages
func ReapiCliClearErrMsg(ctx *ReapiCtx) {
	C.reapi_cli_clear_err_msg((*C.struct_reapi_cli_ctx)(ctx))
}

// ReapiCliSetHandleSet emulates setting the opaque handle to the reapi cli context.
// \param ctx       reapi_cli_ctx_t context object
// \param h         Opaque handle. How it is used is an implementation
//
//	                 detail. However, when it is used within a Flux's
//	                service cli, it is expected to be a pointer
//	                 to a flux_t object.
//	\return          0 on success; -1 on error.
//
// int reapi_cli_set_handle (reapi_cli_ctx_t *ctx, void *handle);
func ReapiCliSetHandle(ctx *ReapiCtx) int {
	return -1
}

// ReapiCliGetHandle emulates setting the opaque handle to the reapi cli context.
//
//	\param ctx       reapi_cli_ctx_t context object
//	\return          handle
//
// void *reapi_cli_get_handle (reapi_cli_ctx_t *ctx);
func ReapiCliGetHandle(ctx *ReapiCtx) int {
	return -1
}
