set(TARGET main)
set(SRCS main.go)
set(CGO_CFLAGS "-I${CMAKE_BINARY_DIR}/resource/reapi/bindings/c")
set(BUILD_SCRIPT "go-build.sh")

# This is currently passed but not used because when passed into add_custom_command the spaces are escaped
set(CGO_LIBRARY_FLAGS "-L${CMAKE_BINARY_DIR}/resource/reapi/bindings -L${CMAKE_BINARY_DIR}/resource/libjobspec -ljobspec_conv -lreapi_cli -L${CMAKE_BINARY_DIR}/resource -lresource -lflux-idset -lstdc++ -lczmq -ljansson -lhwloc -lboost_system -lflux-hostlist -lboost_graph -lyaml-cpp")

# This ensures the main binary is cleaned up
set_directory_properties(
    PROPERTIES
      ADDITIONAL_CLEAN_FILES "${CMAKE_CURRENT_BINARY_DIR}/main;${CMAKE_CURRENT_BINARY_DIR}/${BUILD_SCRIPT}"
  )

# Build a go program       main      main.go
BUILD_GO_PROGRAM(${TARGET} ${SRCS} "${BUILD_SCRIPT}" "${CGO_CFLAGS}" "${CGO_LIBRARY_FLAGS}" jobspec_conv reapi_module reapi_cli resource flux::core)