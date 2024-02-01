import request from "@/request/init";

const baseSiteUrl = "http://127.0.0.1:8888";

function DeleteTask(id) {
    return new Promise((resolve, reject) => {
        let body = {
            id: id,
        };
        request.post(baseSiteUrl+'/task/delete', body).then(
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

function InsertTask(body) {
    return new Promise((resolve, reject) => {
        let insertURL = baseSiteUrl+'/task/create';
        request.post(insertURL, body).then(
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

function UpdateTask(body) {

    return new Promise((resolve, reject) => {
        let updateURL = baseSiteUrl+'/task/update';
        request.post(updateURL, body).then(
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
                    let task_status;
                    switch (res.data.list[index].status) {
                        case 0:
                            task_status="未开始"
                            break
                        case 1:
                            task_status="运行中"
                            break
                        case 2:
                            task_status="执行失败"
                            break
                        case 3:
                            task_status="完成"
                            break
                    }
                    result.list.push({
                         id:         res.data.list[index].id,
                         name:       res.data.list[index].name,
                         video_type: res.data.list[index].video_type,
                         type:       res.data.list[index].type,
                         data:       res.data.list[index].data,
                         status:     task_status,
                         score:      res.data.list[index].score / 100,
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

function GetProgramStatus(body) {
    return new Promise((resolve, reject) => {
        let runURL = baseSiteUrl+'/task/status';
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
    DeleteTask,
    InsertTask,
    UpdateTask,
    GetProgramStatus,
};
