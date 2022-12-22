


echo "生成rpc代码"

OUT=../pb
protoc \
--go_out=${OUT} \
--go-grpc_out=${OUT} \
--go-grpc_opt=require_unimplemented_servers=false \
request_service.proto push_service.proto

echo "生成完毕"

