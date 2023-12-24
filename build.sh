mkdir build
mkdir build/etc
go build -o ./build/dv.exe ./internel/serve/api
cp ./internel/serve/api/etc/task_serve.yaml ./build/etc/task_serve.yaml
