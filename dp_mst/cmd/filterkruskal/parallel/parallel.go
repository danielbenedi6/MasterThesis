/* ----------------------------------------------------------------------------
 * This file was automatically generated by SWIG (https://www.swig.org).
 * Version 4.1.0
 *
 * Do not make changes to this file unless you know what you are doing - modify
 * the SWIG interface file instead.
 * ----------------------------------------------------------------------------- */

// source: parallel.i

package parallel

/*
#define intgo swig_intgo
typedef void *swig_voidp;

#include <stddef.h>
#include <stdint.h>


typedef ptrdiff_t intgo;
typedef size_t uintgo;



typedef struct { char *p; intgo n; } _gostring_;
typedef struct { void* array; intgo len; intgo cap; } _goslice_;


typedef long long swig_type_1;
typedef long long swig_type_2;
extern void _wrap_Swig_free_parallel_93f5361e022132e6(uintptr_t arg1);
extern uintptr_t _wrap_Swig_malloc_parallel_93f5361e022132e6(swig_intgo arg1);
extern swig_voidp _wrap_ParallelFilterKruskal_parallel_93f5361e022132e6(swig_voidp arg1, swig_type_1 arg2, swig_type_2 arg3);
#undef intgo
*/
import "C"

import "unsafe"
import _ "runtime/cgo"
import "sync"
import "dp_mst/internal/common"



type _ unsafe.Pointer



var Swig_escape_always_false bool
var Swig_escape_val interface{}


type _swig_fnptr *byte
type _swig_memberptr *byte


func getSwigcptr(v interface { Swigcptr() uintptr }) uintptr {
	if v == nil {
		return 0
	}
	return v.Swigcptr()
}


type _ sync.Mutex

func Swig_free(arg1 uintptr) {
	_swig_i_0 := arg1
	C._wrap_Swig_free_parallel_93f5361e022132e6(C.uintptr_t(_swig_i_0))
}

func Swig_malloc(arg1 int) (_swig_ret uintptr) {
	var swig_r uintptr
	_swig_i_0 := arg1
	swig_r = (uintptr)(C._wrap_Swig_malloc_parallel_93f5361e022132e6(C.swig_intgo(_swig_i_0)))
	return swig_r
}


func ParallelFilterKruskal(arg1 *common.Edge, arg2 int64, arg3 int64) (_swig_ret *common.Edge) {
	var swig_r *common.Edge
	_swig_i_0 := arg1
	_swig_i_1 := arg2
	_swig_i_2 := arg3
	swig_r = (*common.Edge)(C._wrap_ParallelFilterKruskal_parallel_93f5361e022132e6(C.swig_voidp(_swig_i_0), C.swig_type_1(_swig_i_1), C.swig_type_2(_swig_i_2)))
	return swig_r
}


