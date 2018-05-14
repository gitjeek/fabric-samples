var express = require('express');
var app = express();
var query = require('../fabcar/fabric-api');
var test = require("./test");
var HttpApi = require("./HttpApi");
 
//  主页输出 "Hello World"
app.get('/', function (req, res) {
   console.log("主页 GET 请求");
   res.send("hello world");
})
 
app.get('/query', function (req, res) {
   console.log("主页 GET 请求 QUERY");
   var data = query.query();
   console.log(data);
   HttpApi.send(data,res);
})
 
//  POST 请求
app.post('/', function (req, res) {
   console.log("主页 POST 请求");
   var data = req.body();
   res.send(data);
})
 
//  /del_user 页面响应
app.get('/del_user', function (req, res) {
   console.log("/del_user 响应 DELETE 请求");
   res.send('删除页面');
})
 
//  /list_user 页面 GET 请求
app.get('/list_user', function (req, res) {
   console.log("/list_user GET 请求");
   res.send('用户列表页面');
})
 
// 对页面 abcd, abxcd, ab123cd, 等响应 GET 请求
app.get('/ab*cd', function(req, res) {   
   console.log("/ab*cd GET 请求");
   res.send('正则匹配');
})
 

 
var server = app.listen(8081, function () {
 
  var host = server.address().address
  var port = server.address().port
 
  console.log("应用实例，访问地址为 http://%s:%s", host, port)
 
})
