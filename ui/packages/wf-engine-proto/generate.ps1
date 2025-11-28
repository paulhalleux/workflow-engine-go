# powershell
$PROTOC_GEN_TS_PATH = "../../node_modules/.bin/protoc-gen-ts_proto.exe"
$OUT_DIR = "./src"

Write-Host "Generating WebSocket protocol buffers..."
Write-Host "Output directory: $OUT_DIR"

protoc --plugin="protoc-gen-ts_proto=$PROTOC_GEN_TS_PATH" --ts_proto_out="$OUT_DIR" --ts_proto_opt=outputIndex=true --proto_path ../../../proto/definition websocket.proto