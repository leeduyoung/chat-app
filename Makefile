redis-up:
	kubectl port-forward -n redis svc/my-redis-master 6379:6379

test-redis:
	cd test && go test redis_test.go

test-all:
	cd test && go test

start:
	go run main.go