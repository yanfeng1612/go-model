import axios from 'axios';
import {Message} from 'element-ui';

const service = axios.create({
    // process.env.NODE_ENV === 'development' 来判断是否开发环境
    // easy-mock服务挂了，暂时不使用了
    // baseURL: 'https://www.easy-mock.com/mock/592501a391470c0ac1fab128',
    timeout: 5000
});

service.interceptors.request.use(
    config => {
        return config;
    },
    error => {
        console.log(error);
        return Promise.reject();
    }
);

service.interceptors.response.use(
    response => {
        if (response.status === 200) {
            return response.data;
        } else {
            Promise.reject();
        }
    },
    error => {
        console.log(error);
        return Promise.reject();
    }
);

export default service;

export const requestData = (config, callback) => {
    if (config != null) {
        axios({
            method: "post",
            url: config.url,
            data: config.data,
            // headers: {'Content-Type': 'application/x-www-form-urlencoded'},
        }).then((res) => {
            if (res.data.code === "200" || res.data.code === "401") {
                config.callback(res.data.data);
            } else if (res.data.code === "300") {
                let message = res.data.data == null ? res.data.message : res.data.data.message;
                Message({
                    type: 'error',
                    message: message
                });
            } else if (res.data.code === "409") {
                Message({
                    type: 'error',
                    message: config.duplicateMsg + '重复，请修改后再次点击保存'
                });
            } else {
                Message({
                    type: 'error',
                    message: '网络错误，请联系管理员'
                });
            }
        }).catch(function (err) {
            console.log("err:", err);
            Message({
                type: 'error',
                message: "网络异常，请稍后重试",
            });
        });
    }
};

export const exportData = (config, callback) => {
    if (config != null) //已登录
    {
        axios({
            method: "post",
            url: config.url,
            responseType: "arraybuffer",
            data: config.data,
            headers: {'Content-Type': 'application/x-www-form-urlencoded'},
        }).then((res) => {
            let headers = res.headers;
            let blob = new Blob([res.data], {
                type: headers['content-type']
            });
            let link = document.createElement('a');
            link.target = "_blank";
            link.href = window.URL.createObjectURL(blob);
            var title = "";
            if (!title) {
                const fileName = headers['content-disposition'];
                title = fileName.includes('filename=') ? fileName.split('=')[1] : '下载的表单文件';
            }
            document.body.appendChild(link); // 修复firefox中无法触发click
            if (config.title != null) {
                title = config.title;
            }
            link.download = title;
            link.click();
        }).catch(function (err) {
            console.log("export : " + err);
            Message({
                type: 'error',
                message: '导出数据错误，请联系管理员'
            });
        });
    }
};
