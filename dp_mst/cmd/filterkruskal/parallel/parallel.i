%module parallel

%ignore Edge;

%insert(go_wrapper) %{
import "dp_mst/internal/common"
%}

%typemap(gotype) Edge "common.Edge"
%typemap(gotype) Edge* "*common.Edge"

%{
#include "parallel.hpp"
%}

%include "parallel.hpp"