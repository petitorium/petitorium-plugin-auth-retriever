# Detect OS
ifeq ($(OS),Windows_NT)
    DETECTED_OS := Windows
    PLUGIN_DIR := $(USERPROFILE)/.config/petitorium/plugins/available
    RM := del /F
    CP := copy
    PLUGIN_EXT := .dll
else
    UNAME_S := $(shell uname -s)
    ifeq ($(UNAME_S),Linux)
        DETECTED_OS := Linux
        PLUGIN_DIR := $(HOME)/.config/petitorium/plugins/available
        RM := rm -f
        CP := cp
        PLUGIN_EXT := .so
    else ifeq ($(UNAME_S),Darwin)
        DETECTED_OS := macOS
        PLUGIN_DIR := $(HOME)/.config/petitorium/plugins/available
        RM := rm -f
        CP := cp
        PLUGIN_EXT := .so
    endif
endif

PLUGIN_NAME := auth-retriever$(PLUGIN_EXT)
PLUGIN_PATH := $(PLUGIN_DIR)/$(PLUGIN_NAME)

clean:
ifeq ($(DETECTED_OS),Windows)
	@if exist "$(PLUGIN_PATH)" $(RM) "$(PLUGIN_PATH)"
	@if exist "$(PLUGIN_NAME)" $(RM) "$(PLUGIN_NAME)"
else
	$(RM) $(PLUGIN_PATH) 2>/dev/null || true
	$(RM) $(PLUGIN_NAME) 2>/dev/null || true
endif

build:
	go build -buildmode=plugin -o $(PLUGIN_NAME) .

install: build
ifeq ($(DETECTED_OS),Windows)
	@if not exist "$(PLUGIN_DIR)" mkdir "$(PLUGIN_DIR)"
	$(CP) $(PLUGIN_NAME) "$(PLUGIN_DIR)\"
else
	mkdir -p $(PLUGIN_DIR)
	$(CP) $(PLUGIN_NAME) $(PLUGIN_DIR)/
endif

.PHONY: build clean install
