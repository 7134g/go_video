import request from "@/request/init";


var baseSiteUrl = "http://127.0.0.1:10086"

function DeleteColumn(tableName , id) {
    var params = {
        table_name: tableName,
        id: id,

    };
    request.delete(baseSiteUrl+'/model/delete', { params: params }).then(
        res => {
            console.log(res);
        }
    ).catch(
        error => {
            console.error(error);
        }
    )
}

function InsertColumn(tableName, body) {
    var insertURL = baseSiteUrl+'/model/insert?table_name=' + tableName
    request.post(insertURL, body).then(
        res => {
            console.log(res);
        }
    ).catch(
        error => {
            console.error(error);
        }
    )

}

function UpdateColumn(tableName, id, body) {
    var updateURL = baseSiteUrl+'/model/update?table_name=' + tableName + '&id=' + id
    request.put(insertURL, body).then(
        res => {
            console.log(res);
        }
    ).catch(
        error => {
            console.error(error);
        }
    )

}

function GetDataList(dataPage) {
    return new Promise((resolve, reject) => {
        var result = {
            list: [],
            total: 0,
        }

        let listURL = baseSiteUrl + '/model/list';
        let body = {
            table_name: dataPage.table_name,
            db_type: dataPage.db_type,
            page: dataPage.page,
            size: dataPage.size
        };
        // console.log(listURL)
        // console.log(JSON.stringify(dataPage))
        // console.log(JSON.stringify(body))
        request.post(listURL, body).then(
            res => {
                if (res === null || !res.hasOwnProperty("code") || res.code !== 200) {
                    console.log("err: ", res);
                    resolve(result);
                    return
                }

                // console.log(JSON.stringify(res.data))
                for (const index in res.data.list) {
                    // console.log(res.data[index])
                    result.list.push(res.data.list[index])
                }
                result.total = res.data.total
                result.message = "success"
                resolve(result);
            }
        ).catch(
            error => {
                result.message = error
                reject(result);
            }
        )
    });
}

function GetTables(db_type) {
    return new Promise((resolve, reject) => {
        let result = {
            tableStruct: [],
            message: ""
        };

        const url = 'http://127.0.0.1:10086/model/tables?db_type=' + db_type;
        // console.log(url)
        request.get(url).then(
            res => {
                if (res === null || !res.hasOwnProperty("code") || res.code !== 200) {
                    console.log("err: ", res);
                    resolve(result);
                    return
                }


                let i = 0;
                for (const key in res.data) {
                    i++;
                    if (!res.data.hasOwnProperty(key)) {
                        continue
                    }
                    const table = {
                        number: i,
                        name: key,
                        fields: res.data[key],
                    }
                    result.tableStruct.push(table)
                }
                // console.log(JSON.stringify(result))
                result.message = "success"
                resolve(result);
            }
        ).catch(
            error => {
                result.message = error
                reject(result);
            }
        )
    });
}


export default {
    GetTables,
    GetDataList,
};
