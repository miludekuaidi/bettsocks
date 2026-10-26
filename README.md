[✈️ 一分机场 便宜量大 三网均衡](https://xn--4gqx1hgtfdmt.com/#/register?code=Ife7C7AX)

[🛩️ 吹雪云 比良心良心的机场](https://xn--9kqs1lo79d.com/#/register?code=84TM9zSm)

月付通通低至两月🎉好吃😋不贵量大管饱

# Proxy Core

给 Telegram Android 客户端（NagramX / exteraGram 一类支持 Python 插件的客户端）用的代理插件。
插件自带 **mihomo** 和 **sing-box** 两个内核，在客户端里一键安装、连接、选节点，不需要另外装代理软件。

> **状态说明**：代码在电脑上做过语法检查和解析逻辑测试，但没有完整的实机验证。
> 不同客户端的插件接口可能有差异，遇到问题请自己fork修复或者修改UI(本人小学毕业 不会写代码）

---

## 功能

- **双内核**：mihomo 和 sing-box 可切换（同一进程只能加载一种，切换后可能需要重启客户端）
- **一键安装核心**：不用填下载链接，插件自动从本项目的 Releases 获取最新核心并校验 SHA-256
- **卡片式控制面板**：状态卡片、节点卡片、延迟检测，跟随深色 / 浅色主题
- **订阅管理**：可添加多个订阅或手动节点链接，支持编辑、删除、更新
- **节点选择**：测速、手动选节点，或使用“自动选择最快节点”
- **手动停止**：随时停止内核，并恢复 Telegram 原来的代理设置
- **避免端口冲突**：本地端口被占用时自动顺延到空闲端口
- **不创建 VPN / TUN**：内核只监听 `127.0.0.1`，只为 Telegram 提供代理，不影响其它应用的网络

## 项目结构

```
.
├── .github/workflows/build-cores.yml   # 构建两个内核并发布到同一个 Release
├── Doublecore/
│   ├── mihomo.go                       # mihomo 的 c-shared 封装
│   └── singbox.go                      # sing-box 的 c-shared 封装
├── proxy_core.plugin                   # Telegram 插件本体
└── README.md
```

## 快速开始

### 1. 构建内核（只需做一次）

1. 把 `Doublecore/`、`proxy_core.plugin`、`README.md` 放进你的 GitHub 仓库，
   把 `build-cores.yml` 放到 `.github/workflows/` 下。
2. 打开仓库的 **Actions → build-cores → Run workflow**。
3. 等待完成。会生成**一个** Release，标签为 `doublecore-N`，包含：

   | 文件 | 说明 |
   | --- | --- |
   | `libmihomo.so.gz` | mihomo 内核（android / arm64） |
   | `libsingbox.so.gz` | sing-box 内核（android / arm64） |
   | `sha256-mihomo.txt` | mihomo 解压后 `.so` 的 SHA-256 |
   | `sha256-singbox.txt` | sing-box 解压后 `.so` 的 SHA-256 |

两个内核并行构建，其中一个失败时，另一个仍会发布。

### 2. 安装插件

把 `proxy_core.plugin` 装进客户端的插件管理。

> 插件默认从 `miludekuaidi/bettsocks` 获取核心。如果你用的是自己的仓库，
> 修改 `proxy_core.plugin` 顶部的 `CORE_REPO = "用户名/仓库名"` 即可。

### 3. 使用

1. 打开 **插件设置 → 打开控制面板**（或侧边栏的“代理控制面板”）。
2. 点 **安装组件**，等待核心下载完成。
3. 点 **订阅管理 → 添加订阅 / 节点链接**，填入订阅地址或分享链接。
4. 点 **连接**。
5. 在节点卡片里点 **检测延迟**，点卡片本身即可选用该节点；也可以使用“自动选择最快节点”。

不用时点 **断开连接**，Telegram 的代理设置会恢复成连接前的样子。

## 控制面板

| 区域 | 说明 |
| --- | --- |
| 状态卡片 | 显示连接状态；按钮会在“安装组件 / 连接 / 断开连接”之间变化 |
| 内核：xxx | 在 mihomo 和 sing-box 之间切换 |
| 订阅管理 | 添加、编辑、删除订阅或手动节点（可一次填多条链接） |
| 更新订阅 | 按当前所有订阅重新生成配置并重启内核，逐个汇报节点数和错误 |
| 更新核心 | 检查并安装项目 Releases 里的最新核心 |
| 节点卡片 | 显示 `协议 · 传输 · 加密`、延迟；点击选用，按钮检测延迟 |

插件设置页还有：端口（默认 `17891`，控制端口为端口 + 1）、随客户端启动时自动连接。

## 两个内核的区别

| | mihomo | sing-box |
| --- | --- | --- |
| 订阅格式 | 通用（base64）订阅、Clash YAML 订阅都支持 | 仅通用（base64 / 链接列表），**不支持 Clash YAML** |
| 订阅下载 | 由内核自己拉取 | 由插件下载并转换为 outbound |
| 协议 | 内核支持的协议（含 Reality 等） | vmess、vless（含 Reality）、trojan、ss（不含插件）、hysteria2、tuic |
| 传输 | 内核支持的传输 | tcp、ws、grpc、httpupgrade、h2 |
| 节点选择 | 分组 `PROXY`（手动）+ `AUTO`（自动测速） | 同左 |

建议：订阅是 Clash YAML 就用 mihomo；订阅是通用链接列表两者都可以。

## 工作原理

- 构建时把内核编译成 Android arm64 的 `.so`（Go `c-shared`），只导出三个函数：
  `Start`、`Stop`、`LastError`。
- 插件用 `ctypes` 加载 `.so`，生成配置文件后调用 `Start`，内核在本机开一个混合代理端口。
- 插件把 Telegram 的代理设置为 `127.0.0.1:端口`，并通过内核的 Clash API（本机，随机密码）读取节点、测速、切换。
- 停止时调用 `Stop`，并把 Telegram 原来的代理设置恢复回去。

## 常见问题

**安装插件提示应用版本不够？**
插件的版本要求是 `>=12.0.0`。如果你的客户端更旧，需要修改 `proxy_core.plugin` 顶部的 `__app_version__`。

**点“安装组件”失败？**
先确认仓库的 Releases 里确实有 `libmihomo.so.gz` / `libsingbox.so.gz`。
直连 GitHub 失败时，插件会在内核已运行的情况下改走本地代理重试；如果两者都不行，换个网络再试。

**更新订阅后显示 0 个节点？**
结果弹窗会列出每个订阅的错误原因。常见情况：订阅需要代理才能下载、订阅内容为空或格式不被识别；
sing-box 遇到 Clash YAML 订阅会提示不支持，请切换到 mihomo。

**切换内核后无法启动？**
同一个进程里只能加载一种 Go 内核（混用不稳定）。已经加载过一种后再切换，需要重启客户端，
插件会弹出“立即重启客户端”的提示；如果没有自动拉起，手动重新打开即可。

**会影响我手机上其它代理软件吗？**
内核只监听 `127.0.0.1`，不创建 VPN / TUN，也不改系统路由。端口被占用时会自动换端口。
如果你仍然遇到冲突，可以随时点“停止内核”，Telegram 的代理设置会恢复。

**语音通话可以走代理吗？**
目前没有对语音通话做专门适配，不保证可用。

## 已知限制

- 仅支持 **Android arm64-v8a**，需要 Android 7.0 及以上。
- 同一进程只能加载一种内核。
- sing-box 模式不支持 Clash YAML 订阅、带插件的 ss 节点、tcp 的 http 伪装。
- 节点卡片里的协议 / 传输 / 加密信息是从订阅缓存中扫描出来的，Clash YAML 使用简化解析，特殊写法可能显示为“—”。
- 控制面板最多显示前 300 个节点。
- 内核封装代码按 mihomo（`Alpha` 分支）和 sing-box 1.12.x 的接口编写，上游接口变化时可能需要调整。

## 安全说明

- 插件只从你配置的仓库（`CORE_REPO`）下载核心，并校验 Release 中提供的 SHA-256。
- 核心是原生代码，会被加载进 Telegram 进程，请只使用你自己构建、自己信任的仓库。
- 内核的控制接口只监听本机，并使用每次启动随机生成的密码。

## 致谢

- [mihomo](https://github.com/MetaCubeX/mihomo)
- [sing-box](https://github.com/SagerNet/sing-box)

## 许可

请在此处填写你的许可证。注意 mihomo 和 sing-box 各自有自己的许可证，分发构建产物时请遵守。
