for loop in "build" "build/etc"
do
    if [ -d "$loop" ]; then
        echo "Folder already exists" "$loop"
    else
        mkdir "$loop"
        echo "Folder created successfully"
    fi
done

go build -o ./build/serve ./internel/serve/api/task_serve.go
go build -o ./build/dv ./internel/serve/api/tool/main.go
cp ./internel/serve/api/etc/task_serve.yaml ./build/etc/task_serve.yaml
