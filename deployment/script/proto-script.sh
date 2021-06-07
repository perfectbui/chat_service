DIR=$GOPATH/src/github.com/perfectbui/chat
TEMP_DIR=$DIR/proto
rm -rf $DIR/pb


  protoc --proto_path=proto proto/*.proto --go_out=plugins=grpc:$GOPATH/src/. --grpc-gateway_out=:$GOPATH/src/.
  protoc --proto_path=proto proto/dto/*.proto --go_out=plugins=grpc:$GOPATH/src/. 
  protoc --proto_path=proto proto/types/*.proto --go_out=plugins=grpc:$GOPATH/src/. 

