GO=go
GOCOVER=$(GO) tool cover
GOTEST=$(GO) test
COVERFILE=coverage.out

.PHONY: test test/cover

test:
	$(GOTEST) -v -race ./...
	
test/cover:
	$(GOTEST) -v -coverpkg=./... -covermode=atomic -coverprofile=$(COVERFILE) ./...
	$(GOCOVER) -func=$(COVERFILE)
	$(GOCOVER) -html=$(COVERFILE)
	rm $(COVERFILE)