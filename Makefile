OUT_DIR = bin
FILES = $(shell ls -d cmd/*/ | cut -d/ -f2)
FILES_OUT = $(addprefix ${OUT_DIR}/,${FILES})
PROTO_FILE = $(shell ls -d pb/*/*.proto)
PROTO_GEN = $(shell echo $(PROTO_FILE) | sed -e "s/.proto$\/.pb.go/g" | xargs -n 1 printf "pkg/core/%s\n")

ifeq ($(OS),Windows_NT)
	PLATFORM ?= windows
	DEST ?= windows
else ifeq ($(UNAME_M),x86_64)
	PLATFORM ?= linux
endif

ARCH =

all:
ifeq ($(OS),Windows_NT)
	make win
else
	make linux
endif

linux: PLATFORM = linux
linux: ${OUT_DIR} ${FILES_OUT}

win: PLATFORM = windows
win: ${OUT_DIR} ${FILES_OUT:=.exe}

arm32: PLATFORM = linux
arm32: ARCH = arm
arm32: ${OUT_DIR} ${FILES_OUT}

dist:
ifeq ($(OS),Windows_NT)
	@sh ./scripts/build_artifact_win.sh ${DEST}
endif

.FORCE:
${OUT_DIR}/%: .FORCE
	@echo compiling $(@)...
	@GOOS=$(PLATFORM) GOARCH=$(ARCH) go build -o $(@) -tags $(PLATFORM) ./cmd/$(basename ${@F})

clean:
	$(foreach file,$(shell ls -d pkg/core/pb/*/*.pb.go 2>/dev/null),rm $(file);)
	rm ${OUT_DIR} -rf

${OUT_DIR}:
	@echo create output dir...
	@mkdir ${OUT_DIR}

.PHONY: all clean win linux arm32 .FORCE