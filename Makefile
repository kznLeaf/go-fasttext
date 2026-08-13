# Build PIC static library for CGO: build/libgo_fasttext.a

CXX ?= c++
AR ?= ar
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
CXXFLAGS ?= -std=c++17 -O3 -fPIC -pthread -funroll-loops
CXXFLAGS += -IfastText/src -Icwrapper

BUILD_DIR := build
OBJ_DIR := $(BUILD_DIR)/obj
LIB_INSTALL_DIR := libs/$(GOOS)_$(GOARCH)

FT_SRCS := $(filter-out fastText/src/main.cc,$(wildcard fastText/src/*.cc))
FT_OBJS := $(patsubst fastText/src/%.cc,$(OBJ_DIR)/ft_%.o,$(FT_SRCS))

WRAPPER_SRC := cwrapper/fasttext_wrapper.cpp
WRAPPER_OBJ := $(OBJ_DIR)/fasttext_wrapper.o

LIB := $(BUILD_DIR)/libgo_fasttext.a

.PHONY: all lib lib-install clean test

all: lib

lib: $(LIB)

lib-install: $(LIB)
	mkdir -p $(LIB_INSTALL_DIR)
	cp $(LIB) $(LIB_INSTALL_DIR)/libgo_fasttext.a

$(LIB): $(FT_OBJS) $(WRAPPER_OBJ) | $(BUILD_DIR)
	$(AR) rcs $@ $^

$(OBJ_DIR)/ft_%.o: fastText/src/%.cc | $(OBJ_DIR)
	$(CXX) $(CXXFLAGS) -c $< -o $@

$(WRAPPER_OBJ): $(WRAPPER_SRC) cwrapper/fasttext_wrapper.h | $(OBJ_DIR)
	$(CXX) $(CXXFLAGS) -c $(WRAPPER_SRC) -o $@

$(BUILD_DIR) $(OBJ_DIR):
	mkdir -p $@

test:
	go test ./...

clean:
	rm -rf $(BUILD_DIR)
