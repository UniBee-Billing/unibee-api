GOOS=linux GOARCH=amd64 go build -o temp/linux_amd64/main main.go
docker build -f manifest/docker/Dockerfile -t registry.cn-shenzhen.aliyuncs.com/heiku/heiku_gooverseapay:daily .
docker push registry.cn-shenzhen.aliyuncs.com/heiku/heiku_gooverseapay:daily
kubectl rollout restart deployment/yd-unib-server-daily