SRC_PATH:= ${PWD}

gen:
	@echo "Protobuf generate"
	${SRC_PATH}/deployment/script/proto-script.sh
run:
	go run *.go