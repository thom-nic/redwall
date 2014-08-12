 

export GOPATH=$(CURDIR)

GC=go build
GOGET=go get
GOFMT=gofmt -w
MAIN=redwall

all: build

build: format
	$(GC) $(MAIN).go
	$(LD) -o $(MAIN).out $(MAIN).$O

format:
	$(GOFMT) $(MAIN).go

get:
	$(GOGET)

clean:
	rm *.8 *.o *.out *.6
