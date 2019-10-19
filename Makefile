
default: test build

test:
	sh -c "'$(CURDIR)/scripts/test.sh'"

build:
	sh -c "'$(CURDIR)/scripts/build.sh'"

test_e2e:
	sh -c "'$(CURDIR)/scripts/test_e2e.sh'"

install:
	# todo

report:
	mkdir -p build/report
	goreporter -p . -r build/report

.PHONY: build test