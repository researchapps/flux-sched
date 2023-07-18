set(TARGET main)
set(SRCS main.go)
set(CGO_CFLAGS "-I${CMAKE_BINARY_DIR}/resource/reapi/bindings/c")
set(CGO_LIBRARY_FLAGS "-L${CMAKE_BINARY_DIR}/resource/reapi/bindings;-lreapi_cli;-L${CMAKE_BINARY_DIR}/resource;-lresource;-lflux-idset;-lstdc++;-lczmq;-ljansson;-lhwloc;-lboost_system;-lflux-hostlist;-lboost_graph;-lyaml-cpp")

# This ensures the main binary is cleaned up
set_directory_properties(
    PROPERTIES
      ADDITIONAL_CLEAN_FILES "${CMAKE_CURRENT_BINARY_DIR}/main"
  )

# GO_GET(go_redis github.com/hoisie/redis)
#                          main      main.go  
ADD_GO_INSTALLABLE_PROGRAM(${TARGET} ${SRCS} "${CGO_CFLAGS}" "${CGO_LIBRARY_FLAGS}" reapi_module reapi_cli flux::core)