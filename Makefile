

procs-account = $(shell ps -ef | grep  smdx_account | grep -v grep | awk  '{ print $$2 ; }')
kill-cmd-account = $(if $(procs-account), "kill" "-9" $(procs-account), "echo" "no matching processes")

#-----------------------------------------EVN----------------------------------------------
CUR_DIR = $(shell pwd)
ETCD_HOSTS = 127.0.0.1
KAFKA_CA_FILE = $(CUR_DIR)/docs/aliyun_orog_kafka_ssl.pem

#-----------------------------------------Gen rpcx proto-----------------------------------
gen_common_proto:
	protoc --proto_path=./proto/ ./proto/common/tx.proto --go_out=./rpcx/chains --go-grpc_out=./rpcx/chains
	protoc --proto_path=./proto/ ./proto/common/trade.proto --go_out=./rpcx/chains --go-grpc_out=./rpcx/chains
	protoc --proto_path=./proto/ ./proto/common/kline.proto --go_out=./rpcx/chains --go-grpc_out=./rpcx/chains
	cp -r ./rpcx/chains/github.com/simance-ai/smdx/rpcx/chains/common ./rpcx/chains
	rm -rf ./rpcx/chains/github.com

gen_account_proto: gen_common_proto
	goctl rpc protoc --proto_path ./proto/ ./proto/account.proto --go_out=./rpcx/account --go-grpc_out=./rpcx/account --zrpc_out=./rpcx/account

gen_order_proto: gen_common_proto
	goctl rpc protoc --proto_path ./proto/ ./proto/order.proto --go_out=./rpcx/order --go-grpc_out=./rpcx/order --zrpc_out=./rpcx/order


gen_order_consumer_proto: gen_common_proto
	goctl rpc protoc --proto_path ./proto/ ./proto/order_consumer.proto --go_out=./rpcx/order_consumer --go-grpc_out=./rpcx/order_consumer --zrpc_out=./rpcx/order_consumer


gen_ws_proto: gen_common_proto
	goctl rpc protoc --proto_path ./proto/ ./proto/ws.proto --go-grpc_out=./rpcx/ws --go_out=./rpcx/ws --zrpc_out=./rpcx/ws

gen_rebate_proto: gen_common_proto
	goctl rpc protoc --proto_path ./proto/ ./proto/rebate.proto --go-grpc_out=./rpcx/rebate --go_out=./rpcx/rebate --zrpc_out=./rpcx/rebate

gen_chains_eth_proto: gen_common_proto
	goctl rpc protoc --proto_path ./proto/ ./proto/chains/eth.proto --go-grpc_out=./rpcx/chains/eth --go_out=./rpcx/chains/eth --zrpc_out=./rpcx/chains/eth

gen_chains_sol_proto: gen_common_proto
	goctl rpc protoc --proto_path ./proto/ ./proto/chains/sol.proto --go-grpc_out=./rpcx/chains/sol --go_out=./rpcx/chains/sol --zrpc_out=./rpcx/chains/sol

gen_chains_bsc_proto: gen_common_proto
	goctl rpc protoc --proto_path ./proto/ ./proto/chains/bsc.proto --go-grpc_out=./rpcx/chains/bsc --go_out=./rpcx/chains/bsc --zrpc_out=./rpcx/chains/bsc


gen_proto_all: gen_common_proto gen_account_proto gen_order_proto gen_ws_proto gen_chains_eth_proto gen_chains_sol_proto gen_chains_bsc_proto gen_rebate_proto


gen_app_api:
	goctl api go -api ./app/app.api -dir ./app

#-----------------------------------------Gen rpcx model---------------------------------

#gen_sol_mongo_model:
#	goctl model mongo --type Market --dir ./rpcx/chains/sol/internal/model --cache --easy
#	goctl model mongo --type MarketTx --dir ./rpcx/chains/sol/internal/model --cache --easy
#	goctl model mongo --type MarketToken --dir ./rpcx/chains/sol/internal/model --cache --easy
#	goctl model mongo --type MarketKline1m --dir ./rpcx/chains/sol/internal/model --cache --easy
#	goctl model mongo --type MarketTime5m --dir ./rpcx/chains/sol/internal/model --cache --easy

gen_sol_pg_model:
	gentool -db postgres \
	-dsn "postgresql://[替换自己的]:5432/user_3zwEat" \
	-tables="market,market_tx,market_real_time_data,market_audit_media,market_kline_1s,market_kline_1m,market_kline_5m,market_kline_15m,market_kline_30m, \
	market_kline_1h,market_kline_4h,market_kline_6h,market_kline_12h,market_kline_1d" \
	-outPath ./rpcx/chains/sol/internal/dao/dbx


gen_order_sql_model:
	goctl model pg datasource --url="postgresql://[替换自己的]" --cache --schema public --table="order" --dir ./rpcx/order/internal/model/

gen_order_consumer_sql_model:
	goctl model pg datasource --url="postgresql://[替换自己的]" --cache --schema public --table="order" --dir ./rpcx/order_consumer/internal/model/

gen_account_pg_model:
	goctl model pg datasource --url="postgresql://[替换自己的]" --cache --schema public --table="account" --dir ./rpcx/account/internal/model/

gen_account_user_token_follow_sql_model:
	goctl model pg datasource --url="postgresql://[替换自己的]" --cache --schema public --table="user_token_follow" --dir ./rpcx/account/internal/model/


#-----------------------------------------Run Backend rpcx---------------------------------
kill_account:
	@echo "--> Kill account"
	@$(kill-cmd-account)

build_account:
	@echo "--> Build account"
	@go build -mod=readonly $(BUILD_FLAGS) -o ./build/smdx_account ./rpcx/account/

run_backend_account: kill_account build_account
	export ETCD_HOSTS=$(ETCD_HOSTS)
	echo "--> Installing run_account install and run"
	nohup ./build/smdx_account -f ./rpcx/account/etc/account.yaml >./build/smdx_account.log &

#-----------------------------------------Run rpcx-------------------------------------
run_account:
	export ETCD_HOSTS=$(ETCD_HOSTS)
	go run ./rpcx/account/account.go -f ./rpcx/account/etc/account.yaml

run_ws:
	env KAFKA_CA_FILE=$(KAFKA_CA_FILE) ETCD_HOSTS=$(ETCD_HOSTS) \
	go run ./rpcx/ws/ws.go -f ./rpcx/ws/etc/ws.yaml

run_chains_sol:
	env KAFKA_CA_FILE=$(KAFKA_CA_FILE) ETCD_HOSTS=$(ETCD_HOSTS) \
	go run ./rpcx/chains/sol/sol.go -f ./rpcx/chains/sol/etc/sol.yaml

run_chains_eth:
	export ETCD_HOSTS=$(ETCD_HOSTS)
	export KAFKA_CA_FILE=$(KAFKA_CA_FILE)
	go run ./rpcx/chains/eth/eth.go -f ./rpcx/chains/eth/etc/eth.yaml

run_chains_bsc:
	export ETCD_HOSTS=$(ETCD_HOSTS)
	export KAFKA_CA_FILE=$(KAFKA_CA_FILE)
	go run ./rpcx/chains/bsc/bsc.go -f ./rpcx/chains/bsc/etc/bsc.yaml

run_chains_all: run_chains_sol run_chains_eth run_chains_bsc

run_app:
	export ETCD_HOSTS=$(ETCD_HOSTS)
	export KAFKA_CA_FILE=$(KAFKA_CA_FILE)
	go run ./app/app.go -f ./app/etc/app-api.yaml

run_order:
	export ETCD_HOSTS=$(ETCD_HOSTS)
	export KAFKA_CA_FILE=$(KAFKA_CA_FILE)
	go run ./rpcx/order/order.go -f ./rpcx/order/etc/order.yaml

run_order_consumer:
	export ETCD_HOSTS=$(ETCD_HOSTS)
	export KAFKA_CA_FILE=$(KAFKA_CA_FILE)
	go run ./rpcx/order_consumer/orderconsumer.go -f ./rpcx/order_consumer/etc/orderconsumer.yaml


run_rebate:
	export ETCD_HOSTS=$(ETCD_HOSTS)
	export KAFKA_CA_FILE=$(KAFKA_CA_FILE)
	go run ./rpcx/rebate/rebate.go -f ./rpcx/rebate/etc/rebate.yaml


#----------------------------------------Docker Build-------------------------------------------
build_ws_docker:
	docker buildx build -[替换自己的]/smdx_hub/ws:v1.0 -f ./scripts/dk/Dockerfile_ws .

build_account_docker:
	docker buildx build -[替换自己的]/smdx_hub/account:v1.0 -f ./scripts/dk/Dockerfile_account .

build_chains_sol_docker:
	docker buildx build -[替换自己的]/smdx_hub/chains_sol:v1.0 -f ./scripts/dk/Dockerfile_sol .

build_order_docker:
	docker buildx build -[替换自己的]/smdx_hub/order:v1.0 -f ./scripts/dk/Dockerfile_order .

build_app_docker:
	docker buildx build -[替换自己的]/smdx_hub/app:v1.0 -f ./scripts/dk/Dockerfile_app .

#----------------------------------------Docker Build AMD64-------------------------------------------

build_ws_docker_amd64:
	docker buildx build -[替换自己的]/smdx_hub/ws_amd64:v1.0 -f ./scripts/dk/Dockerfile_ws --platform linux/amd64 .

build_account_docker_amd64:
	docker buildx build -[替换自己的]/smdx_hub/account_amd64:v1.0 -f ./scripts/dk/Dockerfile_account --platform linux/amd64 .

build_order_docker_amd64:
	docker buildx build -[替换自己的]/smdx_hub/order_amd64:v1.0 -f ./scripts/dk/Dockerfile_order --platform linux/amd64 .

build_chains_sol_docker_amd64:
	docker buildx build -[替换自己的]/smdx_hub/chains_sol_amd64:v1.0 -f ./scripts/dk/Dockerfile_sol --platform linux/amd64 .

build_app_docker_amd64:
	docker buildx build -[替换自己的]/smdx_hub/app_amd64:v1.0 -f ./scripts/dk/Dockerfile_app --platform linux/amd64 .


swagger:
	goctl api plugin -plugin goctl-swagger="swagger -filename v1.json" -api app/app.api -dir ./swaggers