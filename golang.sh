# below code can be upgrade
if [ ! -d "/go/src/github.com" ]; then
  echo 'go get'
  go get -u github.com/gin-gonic/gin
  go get github.com/lib/pq
  go get github.com/itsjamie/gin-cors
  echo 'finish get'
fi
go get -u google.golang.org/protobuf/cmd/protoc-gen-go
go get github.com/pilu/fresh

cd /go/admin-api/ && fresh
