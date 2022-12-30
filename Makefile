build:	
		@go build -o bin/dgpoker
run: build
		@./bin/dgpoker
test:
		go test	-v	./...
