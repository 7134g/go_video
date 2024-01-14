import axios from "axios";

const request = axios.create({
    timeout: 5000
});

request.interceptors.request.use(
    config => {
        config.headers["Content-Type"] = "application/json";
        return config;
    },
    error => {
        return Promise.reject(error);
    }
);

request.interceptors.response.use(
    response => {
        let res = response.data;

        if (res.data.code === 200) {
            if (res.config.responseType === "blob") {
                return res;
            }

            if (typeof res == "string") {
                res = res ? JSON.parse(res) : res;
            }
            return res;
        } else {
            // 否则抛出一个异常，将错误信息传递给调用方处理
            throw new Error(res.data.message);
        }


    },
    (error) => {
        console.log("err" + error);
        throw new Error(error.message);
    }

);

export default request;
