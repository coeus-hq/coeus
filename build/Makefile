ifeq ($(OS), Windows_NT)
	TARGET_MAKEFILE := Makefile.win
else
	UNAME_S := $(shell uname -s)
	ifeq ($(UNAME_S),Linux)
	TARGET_MAKEFILE := Makefile.lin
	endif
	ifeq ($(UNAME_S),Darwin)
	TARGET_MAKEFILE := Makefile.mac
	endif
endif

build:
	@make -f ${TARGET_MAKEFILE} build

clean:
	@make -f ${TARGET_MAKEFILE} clean
