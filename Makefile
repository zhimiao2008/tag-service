#NAME := $(shell basename $(shell git config --get remote.origin.url) | sed 's/\.git//')
#BRANCH := $(shell git symbolic-ref --short HEAD 2>/dev/null)

TMPDIR=tmp
RELEASESDIR=releases

default: all

all: pack ## all

check: ## 检查编译环境
	if ! [ -x "`command -v go`" ]; then  echo 'Error: go is not installed.' >&2 ;exit 1; fi
	if ! [ -x "`command -v protoc`" ]; then  echo 'Error: protoc is not installed.' >&2 ;exit 1; fi
	if ! [ -x "`command -v go-bindata`" ]; then  echo 'Error: go-bindata is not installed.' >&2 ;exit 1; fi

test: ## 测试

gen-proto: ## proto 相关文件处理
	protoc --go_out=plugins=grpc:.  ./proto/*.proto
	protoc -I/usr/local/include -I. -I$GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.14.4/third_party/googleapis
    --grpc-gateway_out=logtostderr=true:. ./proto/*.proto

docs: ## 更新文档
	go-bindata --nocompress -pkg swagger -o pkg/swagger/data.go third_party/swagger/...
	protoc -I/usr/local/include -I. -I$GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.14.4/third_party/googleapis --swagger_out=logtostderr=true:. ./proto/*.proto

build: gen-proto ## 构建
	mkdir -p $(TMPDIR)/bin
	go build -o $(TMPDIR)/bin/rpc -ldflags "-X main.buildTime=`date +%Y-%m-%d,%H:%M:%S` -X main.buildVersion=1.0.0 -X main.gitCommitID=`git rev-parse HEAD` "
	cp -r configs $(TMPDIR)/
	cp -r proto $(TMPDIR)/
pack: clean check test docs build ## 打包项目
	mkdir -p $(RELEASESDIR)
	tar czf $(RELEASESDIR)/$(TMPDIR).tar.gz $(TMPDIR)/
	rm -rf $(TMPDIR)/
clean: ## 清除目录
	rm -rf $(TMPDIR)
	rm -rf $(RELEASESDIR)
	go clean

.DEFAULT: all

help: ## Show this help.
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'