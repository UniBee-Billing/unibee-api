 cd test/out && rm -rf * && cd ../..
 go test ./... -coverprofile=./test/out/coverage.out
 go tool cover -html=./test/out/coverage.out -o ./test/out/coverage.html
 open ./test/out/coverage.html