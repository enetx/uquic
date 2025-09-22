#!/bin/bash -eu

export CXX="${CXX} -lresolv" # required by Go 1.20

compile_go_fuzzer github.com/enetx/uquic/fuzzing/frames Fuzz frame_fuzzer
compile_go_fuzzer github.com/enetx/uquic/fuzzing/header Fuzz header_fuzzer
compile_go_fuzzer github.com/enetx/uquic/fuzzing/transportparameters Fuzz transportparameter_fuzzer
compile_go_fuzzer github.com/enetx/uquic/fuzzing/tokens Fuzz token_fuzzer
compile_go_fuzzer github.com/enetx/uquic/fuzzing/handshake Fuzz handshake_fuzzer
