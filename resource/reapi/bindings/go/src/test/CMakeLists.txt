set(TARGET main)
set(SRCS main.go)

add_custom_command(OUTPUT ${CMAKE_CURRENT_SOURCE_DIR}/${TARGET}
  DEPENDS ${SRCS}
  WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
  COMMAND env GOPATH=${CMAKE_SOURCE_DIR}/resource/reapi/bindings/go GOOS=$(shell uname -s | tr A-Z a-z) GO111MODULE=off CGO_CFLAGS="-I${CMAKE_SOURCE_DIR}/resource/reapi/bindings/c" CGO_LDFLAGS="-L$(top_builddir)/resource/reapi/bindings/c/.libs -lreapi_cli -L$(top_builddir)/resource/.libs -lresource -lflux-idset -lstdc++ -lczmq -ljansson -lhwloc -lboost_system -lflux-hostlist -lboost_graph -lyaml-cpp" $(GOC) build -ldflags '-w' -o "${CMAKE_CURRENT_SOURCE_DIR}/${TARGET}"
  COMMENT "Building Go library")