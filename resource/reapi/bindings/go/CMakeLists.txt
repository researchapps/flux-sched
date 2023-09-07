set(go_files
  src/test/main.go
)

add_cgo_executable(main go_files ${CMAKE_CURRENT_BINARY_DIR})