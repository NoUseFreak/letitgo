
default: test build

test:
	sh -c "'$(CURDIR)/scripts/test.sh'"

build:
	sh -c "'$(CURDIR)/scripts/build.sh'"

test_e2e:
	sh -c "'$(CURDIR)/scripts/test_e2e.sh'"

install:
	# todo

.PHONY: build test