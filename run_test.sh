 cd test/out && rm -rf * && cd ../..

 # use gotestsum without installing
 # help: go run gotest.tools/gotestsum@latest --help
 go run gotest.tools/gotestsum@latest --format=testname >> ./test/out/report.txt
# open ./test/out/report.txt
 # use go test
# go test ./... -coverprofile=./test/out/coverage.out
# go tool cover -html=./test/out/coverage.out -o ./test/out/coverage.html
# open ./test/out/coverage.html