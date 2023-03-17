# mini-push



## 如何使用

### Get ``/push/{token}/{text}``

### Get ``/push/{token}/{title}/{text}``

### Get ``/push/{token}/copy/{text}``

### Get ``/push/{token}/copy/{title}/{text}``

### Post ``/push/{token}``

```json
{
  "title": "title",
  "text": "text",
  "copy": true
}
```

### 兼容PushPeer
 
#### Post ``message/push``

```json
{
  "pushkey": "token",
  "text": "title",
  "desp": "text",
  "type": "type"
}
```

### 兼容telegram格式

#### Post ``/tg-format/{pushId}``

```json
{
  "text": "text",
  "chat_id": 11111111,
  "parse_mode": "type"
}
```


# 感谢

[Telegram Bot](https://t.me/fpush_bot)

[PushLite](https://github.com/xlvecle/PushLite)

[Bark-Chrome-Extension](https://github.com/xlvecle/Bark-Chrome-Extension)
