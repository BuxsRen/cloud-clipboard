# 云内容 WebSocket案例

|   语言    |     框架      |   WebSocket   |
|:-------:|:-----------:|:-------------:|
| Go 1.22 | GoFrame 2.7 | gorilla 1.5.1 |

### DEMO 云剪贴板
```这是基于云内容实现的一个云剪贴板的案例，跨设备剪贴板同步```
> 获取剪贴板内容，将云内容置于剪贴板

#### 在线用例 ```网页版```
> https://fan.jx.cn:82/ct.html

### 用例 ```桌面版```

```shell
# 在命令行窗口运行，会运行在后台，启动后是无窗口模式
CloudContent.exe -mode=client -room=123456
```
- mode 运行模式
- room 房间ID

> 运行桌面版，手机上运行网页版，电脑上`Ctrl+C`手机上实时获取电脑剪贴板内容。手机上粘贴内容，电脑剪贴板上将获取到手机上粘贴的内容

### 开发

#### WebSocket 底层封装
- utility/client ``客户端``
- utility/websocket ``服务端``

#### WebSocket 业务层
- internal/packed/client ``客户端``
- internal/packed/websocket ``服务端``
