# Changeset server

## For installing nodejs
```
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.1/install.sh | bash
export NVM_DIR="$HOME/.local/share/nvm"
[ -s "$NVM_DIR" ] && \. "$NVM_DIR/nvm.sh"
nvm install node
```

## For running the server
```
cd server
npm install
node server.js
```

## For regenerating the static codegen files
```
npm install -g grpc-tools
grpc_tools_node_protoc --js_out=import_style=commonjs,binary:./server/changesetproto/ --grpc_out=grpc_js:./server/changesetproto changeset.proto
```


# Changeset client

## For installing go
```
wget https://go.dev/dl/go1.20.6.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.20.6.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
rm go1.20.6.linux-amd64.tar.gz
go version
```

## For running the client
```
cd client
go run client.go
```

## For regenerating the static codegen files
```
sudo apt update
sudo apt install protobuf-compiler -y
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
export PATH=$PATH:$(go env GOPATH)/bin
source ~/.bashrc
protoc --go_out=./changesetproto --go_opt=paths=source_relative --go-grpc_out=./changesetproto --go-grpc_opt=paths=source_relative changeset.proto
```