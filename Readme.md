# 需要使用修改后的kobe
github.com/rubinus/kobe

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o go-check-k8s

docker build -t rubinus/go-check-k8s:v1.0 .

docker run --rm --name check rubinus/go-check-k8s:v1.0

