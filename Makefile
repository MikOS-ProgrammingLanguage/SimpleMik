GOC = go
CFLAGS = build
MIC_DIR = mik_compiler
MIP_DIR = mik_package_manager
MIC_OUTPUT_NAME = mic
MIP_OUTPUT_NAME = mip

export MPON=$(MIP_OUTPUT_NAME)

.PHONY: mic mip

all: mic mip

mic:
	$(GOC) $(CFLAGS) -o $(MIC_OUTPUT_NAME) ./$(MIC_DIR)
	@echo 'if [[ "$$OSTYPE" == "linux-gnu"* ]]; then mv $(MIC_OUTPUT_NAME) /usr/bin/;elif [[ "$$OSTYPE" == "darwin"* ]]; then mv $(MIC_OUTPUT_NAME) /opt/homebrew/bin/;else echo "OS not supported";fi' | bash

mip:
	$(GOC) $(CFLAGS) -o $(MIP_OUTPUT_NAME) ./$(MIP_DIR)
	@echo 'if [[ "$$OSTYPE" == "linux-gnu"* ]]; then mv $(MIP_OUTPUT_NAME) /usr/bin/; elif [[ "$$OSTYPE" == "darwin"* ]]; then mv $(MIP_OUTPUT_NAME) /opt/homebrew/bin/; else echo "OS not supported"; fi' | bash
