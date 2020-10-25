.PHONY: up down clean

ROOT = $(PWD)

down:
	bash nodes/127.0.0.1/stop_all.sh

up:
	rm -rf keys_certs
	bash nodes/127.0.0.1/start_all.sh

clean:
	rm -rf $(ROOT)/nodes/127.0.0.1/node0/data
	rm -rf $(ROOT)/nodes/127.0.0.1/node1/data
	rm -rf $(ROOT)/nodes/127.0.0.1/node2/data
	rm -rf $(ROOT)/nodes/127.0.0.1/node3/data

define compile
	mkdir -p ./build/$(1)
	docker run --rm -v $(ROOT)/contracts:/sources -v $(ROOT)/build/:/output ethereum/$(SOLC) --overwrite --abi --bin -o /output /sources/$(1).sol
	$(ROOT)/bin/abigen --bin=./build/$(1).bin --abi=./build/$(1).abi --pkg=$(1) --out=./build/$(1)/$(1).go
endef
