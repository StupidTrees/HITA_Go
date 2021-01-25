# 接口表

## 用户系统

### 用户注册

- 请求方式：POST
- 请求地址：/user/sign_up
- 需要token：否
- 请求参数：

|  参数   | 含义  | 示例 |
|  ----  | ----  | ----|
|  username | 用户名 | larry|
| password  | 密码 | 123456|
|gender| 性别|MALE/FEMALE|

- 返回格式：Json

| 参数|含义|示例|
|----|----|----| 
|  code | 状态码 | 200=ok|
| msg| 返回消息|success!|
|data|返回数据||

- data字段详情：Json

| 字段|含义|示例|
|----|----|----| 
|  public_key | 用户公钥 | rsa公钥的pem格式字符串|
| token| 用户token||


### 用户登录

- 请求方式：POST
- 请求地址：/user/sign_up
- 需要token：否
- 请求参数：

|  参数   | 含义  | 示例 |
|  ----  | ----  | ----|
|  username | 用户名 | larry|
| password  | 密码 | 123456|


- 返回格式：Json

| 参数|含义|示例|
|----|----|----| 
|  code | 状态码 | 200=ok|
| msg| 返回消息|success!|
|data|返回数据||

- data字段详情：Json

| 字段|含义|示例|
|----|----|----| 
|  public_key | 用户公钥 | rsa公钥的pem格式字符串|
| token| 用户token||

