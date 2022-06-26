# ysbot

特定のTwitterアカウントの真似をするbot


## ysbot.ini
example
```
[oauth]
consKey     = xxxxxxxxxxxxxxxx
consSecret  = xxxxxxxxxxxxxxxx
accToken    = xxxxxxxxxxxxxxxx
accSecret   = xxxxxxxxxxxxxxxx

[target]
name        = twitter_name
sampleNum   = 100
```

- `[oauth]`  
    OAuthのそれぞれのキー/トークンを指定します
- `[target]`  
    - `name` 真似する対象のTwitterアカウント名を指定します
    - `sampleNum` どれくらいのtweetを参照して真似するかの数を指定します  

