proto:
	protoc -I=./definitions --go_out=./go --ruby_out=./ruby/lib --go-grpc_out=./go ./definitions/*.proto
	grpc_tools_ruby_protoc -I=./definitions --ruby_out=./ruby/lib --grpc_out=./ruby/lib ./definitions/*.proto


ruby-app:
	cd ./ruby && bundle exec ruby ./app.rb


user-service:
	cd ./go && go run ./cmd/user-service/main.go

trade-service:
	cd ./go && go run ./cmd/trade-service/main.go

topics:
	cd ./kafka-tools && ./bin/kafka-topics.sh --create --topic trades --bootstrap-server localhost:9092