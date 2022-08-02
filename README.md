# SAST-2022-Go-LeaderBoard

> pyz_creeper 2022.7.25
>
> 本文档有~~一些~~很多内容借用了[sast2022-django-training/README.md at master · Btlmd/sast2022-django-training (github.com)](https://github.com/Btlmd/sast2022-django-training/blob/master/README.md)
>
> 十分感谢lambda QwQ

​		本次后端作业为重构leader board的后端，目的是让大家熟悉前后端分离的写法。

## 1.功能简介

与我们之前使用的 Leaderboard 一样，我们需要完成的基本功能包括

- 用户提交内容，对内容进行评定后在排行榜中按一定规则显示用户的排名
- 用户可以提交一个 Avatar，作为用户在排行榜中的标识

此外我们还希望包括

- 投票。用户可以给排行榜中的指定用户点赞。



## 2.已有代码介绍

作业文件结构如下所示

```
.
├── config
│   ├── config.json
│   └── parseconfig.go
├── go.mod
├── go.sum
├── main.go
├── model
│   ├── global.go
│   ├── submission.go
│   └── user.go
└── route
    ├── board_handle.go
    ├── init.go
    ├── user_handle.go
    └── vote_handle.go
```

1. config文件夹：修改config.json，输入对应的ip地址和端口号，MySQL用户名和密码以及数据库名称，无需改动parseconfig.go
2. go.mod和go.sum包含了使用的gorm和gin的依赖
3. 任务：
   1. 完成依据`groundtruth.txt`的计分函数。可以参考django作业中给出的python实现逻辑，`groundtruth.txt`放在了model/文件夹下，便于`submission.go`中的函数使用它。~~当然你要是实在不想写给四五个随机数也是可以的。~~
   2. 完成model文件夹下`submission.go`和`user.go`中标记TODO的函数。这两个文件中完成的结构体`submission`和`user`并不强制使用，你可以修改它，完成对应的功能即可。
   3. 完成route文件夹下，`*_handle.go`文件中标记了TODO的handle function，利用这些handle function在init.go中完成路由的创建。其中`user_handle.go`中并没有需要实现的代码，它的实现可以作为实现其他handle function的参考。
   4. 配置MySQL服务器，填写config.json。注意在创建数据库时，可能需要指定字符集为utf8mb4。
   5. 配置`main.go`中运行的端口号。这一步涉及到部署，不做为必做项目。
4. 样例前端 http://front.sast2022.lmd.red/

### 推荐完成顺序

​		首先完成任务4，其次完成1、2，最后完成任务3。如有部署的打算，可以使用tmux直接在服务器上运行编译后的文件，此时要注意修改`config.json`，与服务器上的mysql服务保持一致。



## 3.API接口要求

### 排行榜

```
[GET] /leaderboard
```

该接口给出全部排行榜信息，先按照按照 `score` 降序排列，`score` 一样的，按照提交时间 `time` 升序排列。

对于用户的多次提交，无论分数高低，只返回【最后一次】提交。

#### 响应

```json
[
    {
        "user": "lambda_x",
        "score": 33,
        "subs": [22, 45, 32],
        "time": 1658419888,
        "avatar": "XXXXX",
        "votes":0,
    },
    {
        "user": "lambda_y",
        "score": 0,
        "subs": [1, 0, 0],
        "time": 1658419999,
        "avatar": "XXXXX",
        "votes":1,
    },
    ...
]
```

#### TODO

​		在`board_handle.go`中完成`HandleGetBoard`这一函数，并在`init.go`中加入对应的路由和handle function。注意本接口限制方法为GET。



### 提交历史

```
[GET] /history/<user>
```

该接口提供指定用户的提交历史，按照提交时间 `time` 升序排列

#### 请求参数

- 用户名称，从请求 URL 中获得。

#### 响应

该用户的全部历史提交信息。

```json
[
    {
        "score": 0.37,
        "subs": [1, 0, 0],
        "time": 1658419999
    },
    {
        "score": 99,
        "subs": [99, 99, 99],
        "time": 1658420008
    },
    ...
]
```

#### TODO

- 在`board_handle.go`内完成`HandleUserHistory`函数，并在`init.go`中加入对应的路由，加入的方法可以参考课程讲义的内容。

  - 注意处理用户不存在的情形，例如不存在时返回一个 

    ```json
    {
        "code": -1
    }
    ```

    同时返回4xx的http返回码，说明这一请求不合法（返回400）





### 提交

```
[POST] /submit
```

该接口用于接受用户提交的内容，进行评判，然后更新 Leaderboard。

#### 请求体样例

接收到的请求形如

```json
{
    "user": "lambda_x",
    "avatar": "...",
    "content": "..."
}
```

| 字段    | 说明                                         |
| ------- | -------------------------------------------- |
| user    | 用户名                                       |
| avatar  | 用 base64 编码的用户头像，可以直接视为字符串 |
| content | 用户提交的内容，一个字符串                   |

#### 响应

响应主要包括两部分

- `code` 表示请求的状态，0 为成功，其他表示失败
- `msg` 表示请求的说明文字，可用于前端给用户提示
- `data` 前端可能用到的数据

当用户提交成功的内容合法时，返回以下内容

```json
{
    "code": 0,
    "msg": "提交成功",
    "data": {
        "leaderboard": [
            ... // 与[排行榜]这一接口返回内容相同，为更新后的排行榜
        ]
    }
}
```

当请求参数不全时，返回

```json
{
    "code": 1,
    "msg": "参数不全啊"
}
```

当用户名长于 255 字符，返回

```json
{
    "code": -1,
    "msg": "用户名太长了"
}
```

当请求的 `avatar` 超过 100K 字符时，返回

```json
{
    "code": -2,
    "msg": "图像太大了",
}
```

当检测到用户提交的内容不合法时，返回

```json
{
    "code": -3,
    "msg": "提交内容非法呜呜"
}
```

#### TODO

- 约定：当用户不存在时创建该用户
- 实现该handle function的参考逻辑
  - 检查参数格式是否符合约定（使用`shouldbindJson`函数）
  - 检查参数内容是否符合约定
  - 如不符合约定返回对应的code和msg，同时返回4xx的http状态码
  - 如符合，视用户存在与否创建用户，并调用自己写的评判得分函数，创建submission
  - 依据格式返回成功后的结果。



### 投票

```
[POST] /vote
```

#### 请求体样例

```json
{
	"user": "lambda_x"
}
```

| 字段 | 说明                         |
| ---- | ---------------------------- |
| user | 接收投票的用户，这里的用户名 |

此时该用户的 `vote` 数加 1。

为了防止刷票，我们象征性地用中间件查验请求的User-Agent，代码中给出了参考实现，拒绝了User-Agent为空。你可以完成自己想实现的判断逻辑，但注意，这一中间件有可能会影响postman测试，要自己手动加User-Agent。

#### 响应

对于不符合要求的请求，返回

```json
{
    "code": -1
}
```

否则返回

```json
{
    "code": 0,
    "data": {
        "leaderboard": [
            ... // 与[排行榜]这一接口返回内容相同，为更新后的排行榜
        ]
    }
}
```

#### TODO

- 利用`model/user.go`中的代码可以实现本功能的加票数部分，当然推荐你自己实现，注意其中可能存在的并发问题！



## 4.部署

​		本作业不对部署有要求，可以使用tmux完成，也可以在听完docker课程后，使用docker部署，后者是更被推荐的。

​		



## 5.提交

在代码仓库中提出 Issue，于 Issue 中注明部署地址和代码仓库地址。

如果你没有部署的条件或者不想部署，也可以只提交代码仓库。

如果提交的部署地址不是 `59.66.131.240:XXXXX` 的形式（如自己的域名），可以在 Issue 中进行说明或者将 Issue 编号和学号发至 `liu-md20@mails.tinghua.edu.cn` ，便于我们统计大家的提交情况。







