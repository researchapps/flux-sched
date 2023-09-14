# BUILD_PROGRAM builds a custom program
function(BUILD_PROGRAM NAME MAIN_SRC BUILD_SCRIPT)
  get_filename_component(MAIN_SRC_ABS ${MAIN_SRC} ABSOLUTE)
  get_filename_component(BUILD_SCRIPT_ABS ${BUILD_SCRIPT} ABSOLUTE)
  add_custom_target(${NAME})

  # Find bash to run custom command
  find_program (BASH_EXECUTABLE bash)
  add_custom_command(TARGET ${NAME}
                    COMMAND ${BASH_EXECUTABLE} ${BUILD_SCRIPT_ABS}
                    WORKING_DIRECTORY ${CMAKE_CURRENT_BINARY_DIR}
                    DEPENDS ${MAIN_SRC_ABS}
                    COMMENT "Building custom library")
  foreach(DEP ${ARGN})
    add_dependencies(${NAME} ${DEP})
  endforeach()
  
  add_custom_target(${NAME}_all ALL DEPENDS ${NAME})
  install(PROGRAMS ${CMAKE_CURRENT_BINARY_DIR}/${NAME} DESTINATION bin)
endfunction(BUILD_PROGRAM)