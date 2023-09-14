set(CUSTOM_GO_PATH "${CMAKE_SOURCE_DIR}/resource/reapi/bindings/go")
set(GOPATH "${CMAKE_CURRENT_BINARY_DIR}/go")
set(G0111MODULE "off")
set(GOOS "linux")
file(MAKE_DIRECTORY ${GOPATH})

# BUILD_GO_PROGRAM builds a custom go program (primarily for testing)
function(BUILD_GO_PROGRAM NAME MAIN_SRC GO_BUILD_SCRIPT CGO_CFLAGS CGO_LIBRARY_FLAGS)
  message(STATUS "GOPATH: ${GOPATH}")
  message(STATUS "CGO_LDFLAGS: ${CGO_LIBRARY_FLAGS}")
  get_filename_component(MAIN_SRC_ABS ${MAIN_SRC} ABSOLUTE)
  get_filename_component(GO_BUILD_SCRIPT_ABS ${GO_BUILD_SCRIPT} ABSOLUTE)
  add_custom_target(${NAME})

  # Find bash to run custom command
  find_program (BASH_EXECUTABLE bash)
  file(WRITE "${CMAKE_CURRENT_BINARY_DIR}/${GO_BUILD_SCRIPT}" "GOPATH=${GOPATH}:${CUSTOM_GO_PATH} GOOS=${GOOS} G0111MODULE=${G0111MODULE} CGO_CFLAGS=${CGO_CFLAGS} CGO_LDFLAGS='${CGO_LIBRARY_FLAGS}' go build -ldflags -w -o ${CMAKE_CURRENT_SOURCE_DIR}/${NAME}")
  add_custom_command(TARGET ${NAME}
                    COMMAND ${BASH_EXECUTABLE} ${GO_BUILD_SCRIPT_ABS}
                    WORKING_DIRECTORY ${CMAKE_CURRENT_BINARY_DIR}
                    DEPENDS ${MAIN_SRC_ABS}
                    COMMENT "Building Go library")
  foreach(DEP ${ARGN})
    add_dependencies(${NAME} ${DEP})
  endforeach()
  
  add_custom_target(${NAME}_all ALL DEPENDS ${NAME})
  install(PROGRAMS ${CMAKE_CURRENT_BINARY_DIR}/${NAME} DESTINATION bin)
endfunction(BUILD_GO_PROGRAM)