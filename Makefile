# ==============================================================================
# Vars
# ==============================================================================
DEVVERSION:=$(shell git describe --tags)
VERSION:=$(shell git describe --tags --abbr=0)

#BUILDCMD:=CGO_ENABLED=1 CC="zig cc" CPP="zig c++" go build
BUILDCMD:=CGO_ENABLED=1 go build

WINBUILDNAME:=RayTheater.exe

# ==============================================================================
# Directorys
# ==============================================================================
BUILDDIR:= .

# ==============================================================================
# Recipes
# ==============================================================================

debug.run: main.go
	$(BUILDCMD) -ldflags=-linkmode=internal -o $@
#	$(BUILDCMD) -o $@
#   $(BUILDCMD) -ldflags="-g" -o $@
#	$(BUILDCMD) -ldflags="-X main.Version=$(DEVVERSION)" -o $@


release.run: main.go
	$(BUILDCMD) -ldflags="-s -w" -o $@


$(WINBUILDNAME): main.go
	GOOS=windows GOARCH=amd64 CC="x86_64-w64-mingw32-gcc" $(BUILDCMD) -o $@

#tnt.linux.x86_64: main.go
#	GOOS=linux GOARCH=amd64 CC="zig cc -target x86_64-linux" CXX="zig c++ -target x86_64-linux" $(BUILDCMD) -ldflags="-w -X main.Version=$(VERSION)" -o $@

setup: go.sum
	@echo "setup done"

go.sum: go.mod
	GOPRIVATE="github.com/rocco-gossmann" go mod tidy



.phony: clean remake dev test tst run

run:
	go run -tags x11 .

dev:
	find . -type f -name "*.go" | entr make remake
	
test:
	find . -type f -name "*.go" | entr make tst 

tst:
	clear
	$(BUILDCMD) test

remake: 
	clear
	rm -f ./debug.run
	make debug.run
	
clean:
	rm -rf $(BUILDDIR)/debug.run	
	rm -rf $(BUILDDIR)/$(WINBUILDNAME)
