go test -v -coverpkg ./... ./... -coverprofile=../../test/out/logic/coverage.out
go tool cover -html=../../test/out/logic/coverage.out -o ../../test/out/logic/coverage.html
open ../../test/out/logic/coverage.html