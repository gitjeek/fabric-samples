var sleep = function (time) {
    return new Promise(function (resolve, reject) {
        setTimeout(function () {
            // 返回 ‘ok’
            resolve('ok');
        }, time);
    })
};
var data = null;
var start = async function () {
    let result = await sleep(3000);
    //console.log(result); // 收到 ‘ok’
    data = result;
    console.log(data);
    return data;
};

exports.test = start;

