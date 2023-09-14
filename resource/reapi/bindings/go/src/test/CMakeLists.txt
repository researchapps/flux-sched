set(G0111MODULE "off")
set(GOOS "linux")
set(TARGET main)
set(SRCS main.go)
set(BUILD_SCRIPT "${CMAKE_CURRENT_BINARY_DIR}/go-build.sh")
set(CGO_CFLAGS "-I${CMAKE_BINARY_DIR}/resource/reapi/bindings/c")
set(CUSTOM_GO_PATH "${CMAKE_SOURCE_DIR}/resource/reapi/bindings/go")
set(GOPATH "${CMAKE_BINARY_DIR}/go")
file(MAKE_DIRECTORY ${GOPATH})

# This is currently passed but not used because when passed into add_custom_command the spaces are escaped
set(CGO_LIBRARY_FLAGS "-L${CMAKE_BINARY_DIR}/resource/reapi/bindings -L${CMAKE_BINARY_DIR}/resource/libjobspec -ljobspec_conv -lreapi_cli -L${CMAKE_BINARY_DIR}/resource -lresource -lflux-idset -lstdc++ -lczmq -ljansson -lhwloc -lboost_system -lflux-hostlist -lboost_graph -lyaml-cpp")

# This ensures the main binary is cleaned up
set_directory_properties(
    PROPERTIES
      ADDITIONAL_CLEAN_FILES "${CMAKE_CURRENT_BINARY_DIR}/main;${BUILD_SCRIPT}"
  )

message(STATUS "GOPATH: ${GOPATH}")
message(STATUS "CGO_LDFLAGS: ${CGO_LIBRARY_FLAGS}")
  
# Find bash to run custom command
find_program (BASH_EXECUTABLE bash)
file(WRITE "${BUILD_SCRIPT}" "GOPATH=${GOPATH}:${CUSTOM_GO_PATH} GOOS=${GOOS} G0111MODULE=${G0111MODULE} CGO_CFLAGS=${CGO_CFLAGS} CGO_LDFLAGS='${CGO_LIBRARY_FLAGS}' go build -ldflags -w -o ${CMAKE_CURRENT_SOURCE_DIR}/${TARGET}")
BUILD_PROGRAM(${TARGET} ${SRCS} "${BUILD_SCRIPT}" "${BUILD_SCRIPT}" jobspec_conv reapi_module reapi_cli resource flux::core)