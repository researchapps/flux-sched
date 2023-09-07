
set(CUSTOM_GO_PATH "${CMAKE_SOURCE_DIR}/resource/reapi/bindings/go")
set(GOPATH "${CMAKE_CURRENT_BINARY_DIR}/go")

# This probably isn't necessary, although if we could build fluxcli into it maybe
file(MAKE_DIRECTORY ${GOPATH})

# GO_GET will retrieve a go module that needs to be installed
function(GO_GET TARG)
  add_custom_target(${TARG} env GOPATH=${GOPATH} go get ${ARGN})
endfunction(GO_GET)

# ADD_GO_INSTALLABLE_PROGRAM builds a custom go program (primarily for testing)
function(ADD_GO_INSTALLABLE_PROGRAM NAME MAIN_SRC CGO_CFLAGS CGO_LIBRARY_FLAGS)
  message(STATUS "GOPATH: ${GOPATH}")
  message(STATUS "CGO_LDFLAGS (before): ${CGO_LIBRARY_FLAGS}")
  message(STATUS "TEST_FLAGS: ${TEST_FLAGS}")
  get_filename_component(MAIN_SRC_ABS ${MAIN_SRC} ABSOLUTE)
  add_custom_target(${NAME})

  # string(REPLACE <match-string> <replace-string> <out-var> <input>...)
  STRING(REPLACE ";" " " CGO_LDFLAGS "${CGO_LIBRARY_FLAGS}")  

  # set(ENV{<variable>} [<value>]) as environment OR without CMake variable
  # Note that I couldn't get this to work (the spaces are always escaped) so I hard coded for now
  # We need a solution that takes the CGO_LIBRARY_FLAGS arg, and can pass (with spaces not escaped) to add_custom_command
  SET ($ENV{CGO_LDFLAGS} "${CGO_LDFLAGS}")
  message(STATUS "CGO_LDFLAGS (after): ${CGO_LDFLAGS}")
  # SET (CMAKE_GO_FLAGS "${CGO_LDFLAGS}")

  add_custom_command(TARGET ${NAME}
                    COMMAND GOPATH=${GOPATH}:${CUSTOM_GO_PATH} GOOS=linux G0111MODULE=off CGO_CFLAGS="${CGO_CFLAGS}" CGO_LDFLAGS='-L${CMAKE_BINARY_DIR}/resource/reapi/bindings -L${CMAKE_BINARY_DIR}/resource/libjobspec/ -ljobspec_conv -lreapi_cli -L${CMAKE_BINARY_DIR}/resource -lresource -lflux-idset -lstdc++ -lczmq -ljansson -lhwloc -lboost_system -lflux-hostlist -lboost_graph -lyaml-cpp' go build -ldflags '-w'
                    -o "${CMAKE_CURRENT_SOURCE_DIR}/${NAME}"
                    ${CMAKE_GO_FLAGS} ${MAIN_SRC}
                    WORKING_DIRECTORY ${CMAKE_CURRENT_LIST_DIR}
                    DEPENDS ${MAIN_SRC_ABS}
                    COMMENT "Building Go library")
  foreach(DEP ${ARGN})
    add_dependencies(${NAME} ${DEP})
  endforeach()
  
  add_custom_target(${NAME}_all ALL DEPENDS ${NAME})
  install(PROGRAMS ${CMAKE_CURRENT_BINARY_DIR}/${NAME} DESTINATION bin)
endfunction(ADD_GO_INSTALLABLE_PROGRAM)