import request from "@/request/init";


var baseSiteUrl = "http://127.0.0.1:8888"

function DeleteColumn(id) {
    var body = {
        id: id,
    };
    request.post(baseSiteUrl+'/task/delete', body).then(
        res => {
            console.log(res);
        }
    ).catch(
        error => {
            console.error(error);
        }
    )
}

function InsertColumn(body) {
    var insertURL = baseSiteUrl+'/task/create'
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

function UpdateColumn(body) {
    var updateURL = baseSiteUrl+'/task/update'
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

function GetTaskList(dataPage) {
    return new Promise((resolve, reject) => {
        var result = {
            list: [],
            total: 0,
        }

        let listURL = baseSiteUrl + '/task/list';
        let body = {
            page: dataPage.page,
            size: dataPage.size,
            type: dataPage.type,
            video_type: dataPage.video_type,
        };

        request.post(listURL, body).then(
            res => {

                for (const index in res.data.list) {
                    // console.log(res.data[index])
                    result.list.push(new TaskData(
                        res.data.list[index].id,
                        res.data.list[index].name,
                        res.data.list[index].video_type,
                        res.data.list[index].type,
                        res.data.list[index].data,
                        res.data.list[index].status,
                    ))
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

export default {
    GetTaskList,
};
