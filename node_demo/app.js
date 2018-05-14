var express = require('express');
var app = express();
var fabric_api = require('../fabcar/fabric-api');
var test = require("./test");
var HttpApi = require("./HttpApi");
var bodyParser = require('body-parser');
//mysql模块
var mysql = require('mysql');
var dbConfig = require('./db/DBConfig');

// 使用DBConfig.js的配置信息创建一个MySQL连接池
var pool = mysql.createPool( dbConfig.mysql );

app.use(bodyParser.urlencoded({extended:false}));
app.use("*", function (req, res, next) {
  res.header('Access-Control-Allow-Origin', '*');
  res.header("Access-Control-Allow-Headers", "Content-Type,Content-Length, Authorization, Accept,X-Requested-With");
  res.header("Access-Control-Allow-Methods","PUT,POST,GET,DELETE,OPTIONS");
  if (req.method === 'OPTIONS') {
    res.send(200)
  } else {
    next()
  }
});

// 响应一个JSON数据
var responseJSON = function (res, ret) {
   if(typeof ret === 'undefined') { 
        res.json({     code:'-200',     msg: '操作失败'   
      }); 
  } else { 
    res.json(ret); 
}};

//获得学生数据(小程序端接口)
app.get('/getStdInfo', function(req, res){

  // 从连接池获取连接 
  pool.getConnection(function(err, connection) { 
    // 获取前台页面传过来的参数  
    var param = req.query || req.params;   
    var sql = 'SELECT * FROM student_information WHERE wechat_number = '+param.stu_weixin;
    console.log(sql);
    // 建立连接 增加一个用户信息 
    connection.query(sql, function(err, result) {
      
      if(result) {      
        var std_key = result[0].fabric_key;
        var response = fabric_api.queryStd(std_key);
        HttpApi.send(response,res);
      }     
      // 释放连接  
      connection.release();  

    });
  });
});

//增加学生(web端接口)
app.post('/addStd', function (req, res) {
   console.log("addStd");
   if (req.body.Number && req.body.Name && req.body.Major && req.body.School) {
     var data = fabric_api.addStd(req.body);
     HttpApi.send(data,res);
   } else {
     res.send('error');    
   }
})

//修改学生记录(web端接口)
app.post('/modifyStdInfo', function (req, res) {
   console.log("modifyStdInfo");
   console.log(req.body);
   if (req.body.ID) {
      var response = fabric_api.modifyStdInfo(req.body);
      HttpApi.send(response,res);
   } else {
     res.send('error');    
   }
})

//删除学生(web端接口)
app.post('/deleteStd', function (req, res) {
   console.log("deleteStd");
   console.log(req.body);
   if (req.body.ID) {
      var response = fabric_api.deleteStd(req.body);
      HttpApi.send(response,res);
   } else {
     res.send('error');    
   }
})

//查询所有学生(web端接口)
app.get('/queryAllStds', function (req, res) {
   console.log("queryAllStds");
   var data = fabric_api.queryAllStds();
   console.log(data);
   HttpApi.send(data,res);
})

//查询单个学生(web端接口)
app.post('/queryStd', function (req, res) {
   console.log("queryStd");
   if (req.body.ID) {
      var response = fabric_api.queryStd(req.body.ID);
      HttpApi.send(response,res);
   } else {
     res.send('error');    
   }
})

 
var server = app.listen(8081, function () {
 
  var host = server.address().address
  var port = server.address().port
 
  console.log("应用实例，访问地址为 http://%s:%s", host, port)
 
})
