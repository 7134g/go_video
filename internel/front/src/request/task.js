import request from "@/request/init";

const baseSiteUrl = "http://127.0.0.1:8888";

function DeleteTask(id) {
    let body = {
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

function InsertTask(body) {
    let insertURL = baseSiteUrl+'/task/create';
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

function UpdateTask(body) {
    let updateURL = baseSiteUrl+'/task/update';
    request.post(updateURL, body).then(
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
        let result = {
            list: [],
            total: 0,
        }

        let listURL = baseSiteUrl + '/task/list';
        let body = {
            page: dataPage.page,
            size: dataPage.size,
            where: {
                type: dataPage.where.type,
                video_type: dataPage.where.video_type,
            },
        };

        // console.log("body",dataPage , body)
        request.post(listURL, body).then(
            res => {
                for (const index in res.data.list) {
                    result.list.push({
                         id:         res.data.list[index].id,
                         name:       res.data.list[index].name,
                         video_type: res.data.list[index].video_type,
                         type:       res.data.list[index].type,
                         data:       res.data.list[index].data,
                         status:     res.data.list[index].status,
                    })
                }
                result.total = res.data.total
                result.message = "success"
                // console.log(result)
                resolve(result);
            }
        ).catch(
            error => {
                // console.log("qqqqqqq", error)
                result.message = error
                reject(result);
            }
        )
    });
}

function RunTask(body) {
    return new Promise((resolve, reject) => {
        let runURL = baseSiteUrl+'/task/run';
        request.get(runURL, body).then(
            res => {
                // console.log(res);
                resolve(res);
            }
        ).catch(
            error => {
                console.error(error);
                reject(error);
            }
        )
    })


}

export default {
    GetTaskList,
    RunTask,
};
