# 标准返回格式

JSON
- code 返回码 200表示成功，其余都是失败
- message 返回消息 success表示成功
- data 返回数据

以下api的"返回数据"指的都是data字段！

# 用户系统

## 注册

- 地址 /user/sign_up
- 请求方式 POST
- 请求表单
```
    username 字符串，用户名
    password 字符串，密码
    nickname 字符串，昵称
    gender 字符串，性别，取MALE/FEMALE/OTHER
```
- 返回数据（如果注册成功）
```
    token 字符串，登录令牌（这个需要在登录后保存起来）
    id 数字，用户id，是用户的唯一标识
    username 用户名
    nickname 字符串，昵称
    gender 字符串，性别
    avatar 数字，用户头像id
    signature 字符串，用户个性签名
    studentId 字符串，学号
    school 字符串，学院
    publicKey 字符串，用户公钥
  ```
    
## 登录

- 地址 /user/log_in
- 请求方式 POST
- 请求表单
```
    username 字符串，用户名
    password 字符串，密码
```
- 返回数据（如果登录成功）
```
    token 字符串，登录令牌（这个需要在登录后保存起来）
    id 数字，用户id，是用户的唯一标识
    username 用户名
    nickname 字符串，昵称
    gender 字符串，性别
    avatar 数字，用户头像id
    signature 字符串，用户个性签名
    studentId 字符串，学号
    school 字符串，学院
    publicKey 字符串，用户公钥
```

## 加载头像

直接用url访问/profile/avatar?imageId=用户id 就可以得到头像

## 获取用户资料
- 地址 /profile/get
- 请求方式 GET
- 请求Header Authorization字段需要带"Bearer XXXX"，其中XXX为用户登录后保存的token
- 请求参数 
```
    userId 数字，用户id
```
- 返回数据
```
    id 数字，用户id，是用户的唯一标识
    username 用户名
    nickname 字符串，昵称
    gender 字符串，性别
    avatar 数字，用户头像id
    signature 字符串，用户个性签名
    studentId 字符串，学号
    school 字符串，学院
    followed 布尔，你是否关注了TA
    followingNum 数字，TA关注了多少人
    fansNum 数字 TA有多少个粉丝
```


## 多模式获取用户列表
####需要分页加载
- 地址 /profile/gets
- 请求方式 GET
- 请求Header Authorization字段需要带"Bearer XXXX"，其中XXX为用户登录后保存的token
- 请求参数
```
    mode 模式
    pageSize 分页大小
    pageNum 分页标号
    extra 额外参数
    
    其中，mode可以取：
    (1) search 搜索模式
        此时extra填搜索关键字
    (2) following 查找某人关注的所有人
        此时extra填这个人的id（字符串形式）
    (3) fans 查找某人所有的粉丝
        此时extra填这个人的id（字符串形式）
    (4) liked 查找点赞了某个帖子的所有用户
        此时extra填这个帖子的id（字符串形式）
```

- 返回数据
```
    是个JSONARRAY格式，其中每一项都是：
    {
    id, 数字，用户id，是用户的唯一标识
    username, 用户名
    nickname, 字符串，昵称
    gender, 字符串，性别
    avatar, 数字，用户头像id
    signature, 字符串，用户个性签名
    studentId, 字符串，学号
    school, 字符串，学院
    followed, 布尔，你是否关注了TA
    followingNum, 数字，TA关注了多少人
    fansNum, 数字 TA有多少个粉丝
    }
```


# 社区

## 发帖
- 地址 /article/create
- 请求方式 POST
- 请求Header Authorization字段需要带"Bearer XXXX"，其中XXX为用户登录后保存的token
- 请求表单
```
    content 字符串，帖子正文
    topicId 字符串，话题id
    repostId 字符串，转发自的帖子id
    asAttitude 布尔，是否发表为态度投票
    anonymous 布尔，是否匿名
```
- 返回数据 无

## 获取某个帖子的信息

- 地址 /article/get
- 请求方式 GET
- 请求Header Authorization字段需要带"Bearer XXXX"，其中XXX为用户登录后保存的token
- 请求表单
```
    articleId 字符串，帖子id
    digOrigin 布尔，为true时返回的是这个id所转发的原帖的信息
```
- 返回数据
```
{
	id                 帖子id
	authorId           作者的id
	authorName         作者昵称
	authorAvatar       作者的头像id
	repostId           转发自原帖的id（没转发就是0）
	repostAuthorId     转发自原帖的作者id
	repostAuthorAvatar 转发自原帖的作者头像id
	repostAuthorName   转发自原帖的作者昵称
	repostContent      转发自原帖的正文
	repostImages       转发自原帖的图片id列表
	repostTime         转发自原帖的发布时间
	topicId            所属话题id
	topicName          所属话题名称
	content            正文
	images             图片id列表
	type               字符串，NORMAL普通/VOTE态度投票
	anonymous          是否匿名
	likeNum            点赞数量
	isMine             布尔，是否是我发的
	liked              布尔，我（查询者）是否点赞了这帖
	votedUp            布尔，我（查询者）是否赞同了这贴（如果这贴是态度投票）
	upNum              赞同人数（如果这贴是态度投票）
	downNum            不赞同人数（如果这贴是态度投票）
	starred            布尔，我是否收藏了
	commentNum         数字，评论数
	createTime         时间戳（long），发帖时间
}

```

## 多模式获取帖子列表
####需要分页加载
- 地址 /article/gets
- 请求方式 GET
- 请求Header Authorization字段需要带"Bearer XXXX"，其中XXX为用户登录后保存的token
- 请求表单
```
    afterTime：long，时间戳，表示查询这之后的帖子
    beforeTime：long，时间戳，表示查询这之前的帖子
    pageSize：数字，分页大小
    mode：字符串，模式
    extra：额外参数
    
    其中，mode可以取：
    (1) all 获取全部的帖子
        此时extra不填
    (2) following 获取我（查询者）关注的所有人的发帖
        此时extra不填
    (3) search 按关键字搜索帖子
        此时extra填搜索关键字
    (4) user 获取某个用户的所有帖子
        此时extra填这个用户的id（字符串形式）
    (5) repost 获取转发某个帖子的所有帖子
        此时extra填这个帖子的id（字符串形式）
    (6) star 获取某个用户的收藏所有帖子
        此时extra填这个用户的id（字符串形式）
    (7) topic 获取某个话题的所有帖子
        此时extra填这个话题的id（字符串形式）
```
- 返回数据
    JSONARRAY格式，每一项为（就是上面那个接口返回的内容）

```
{
	id                 帖子id
	authorId           作者的id
	authorName         作者昵称
	authorAvatar       作者的头像id
	repostId           转发自原帖的id（没转发就是0）
	repostAuthorId     转发自原帖的作者id
	repostAuthorAvatar 转发自原帖的作者头像id
	repostAuthorName   转发自原帖的作者昵称
	repostContent      转发自原帖的正文
	repostImages       转发自原帖的图片id列表
	repostTime         转发自原帖的发布时间
	topicId            所属话题id
	topicName          所属话题名称
	content            正文
	images             图片id列表
	type               字符串，NORMAL普通/VOTE态度投票
	anonymous          是否匿名
	likeNum            点赞数量
	isMine             布尔，是否是我发的
	liked              布尔，我（查询者）是否点赞了这帖
	votedUp            布尔，我（查询者）是否赞同了这贴（如果这贴是态度投票）
	upNum              赞同人数（如果这贴是态度投票）
	downNum            不赞同人数（如果这贴是态度投票）
	starred            布尔，我是否收藏了
	commentNum         数字，评论数
	createTime         时间戳（long），发帖时间
}

```