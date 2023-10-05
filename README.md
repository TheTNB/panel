<h1 align="center">耗子 Linux 面板</h1>

<p align="center">
  <a href="https://github.com/haozi-team/panel/releases"><img src="https://img.shields.io/github/release/haozi-team/panel.svg"></a>
  <a href="https://github.com/haozi-team/panel/actions"><img src="https://github.com/haozi-team/panel/actions/workflows/test.yml/badge.svg"></a>
  <a href="https://goreportcard.com/report/github.com/haozi-team/panel"><img src="https://goreportcard.com/badge/github.com/haozi-team/panel"></a>
  <a href="https://codecov.io/gh/haozi-team/panel"><img src="https://codecov.io/gh/haozi-team/panel/branch/main/graph/badge.svg?token=XFT5NGNSRG"></a>
  <a href="https://img.shields.io/github/license/haozi-team/panel"><img src="https://img.shields.io/github/license/haozi-team/panel"></a>
</p>

<p align="center">
[简体中文] | [<a href="README_EN.md">English</a>]
</p>

耗子 Linux 面板是针对我们自身业务需要使用 Golang 开发的轻量 Linux 服务器运维管理面板，以 Apache 2.0 协议开源。

免责声明：严禁使用耗子 Linux 面板从事任何非法活动，非法站点请勿向我们请求任何形式的技术支持，如果在技术支持过程中发现非法内容，我们将立即停止技术支持并留存相关证据。

交流QQ群：[12370907](https://jq.qq.com/?_wv=1027&k=I1oJKSTH) | QQ频道：[pd.qq.com/s/fyol46wfy](https://pd.qq.com/s/fyol46wfy)

## UI 展示

![UI展示](ui.png)

## 运行环境

耗子Linux面板仅支持 `amd64` | `arm64` 架构下的主流系统的最新版本，不支持 `Ubuntu`，因为其发版太过频繁，难以维护。

低配机器建议使用 `Debian`，资源占用较 `RHEL` 系更低。其他机器建议使用 `AlmaLinux` | `RockyLinux`，维护周期更长也更稳定。

不在下表中的其他系统（OpenCloudOS 8、Anolis 8、CentOS Stream 8/9、Debian 11等），可自行尝试安装，但不保证能够正常运行，且不提供无偿技术支持（理论上不会有大问题）。

CentOS Stream 可使用迁移脚本迁移至支持的系统: [CentOS 8/9 迁移脚本](https://github.com/haozi-team/byecentos)

| 系统         | 版本 |
|------------|----|
| RHEL       | 9  |
| AlmaLinux  | 9  |
| RockyLinux | 9  |
| Debian     | 12 |

随着系统版本的不断更新，我们亦可能会终止部分过于老旧的系统的支持，以保证面板的稳定性。

## 安装面板

安装面板前，你需要了解LNMP环境的基本知识，以及如何处理常见的LNMP环境问题，不建议0基础的用户安装和使用耗子Linux面板（[推荐: 宝塔 - 简单好用服务器运维面板](https://www.bt.cn/?invite_code=M190eXRpZWE=)）。

如果你决定继续，请以`root`用户登录服务器，执行以下命令安装面板：

```shell
bash <(curl -sSL https://dl.cdn.haozi.net/panel/install_panel.sh)
```

## 日常维护

使用`panel`命令进行日常维护：

```shell
panel
```

在 [Wiki](https://github.com/haozi-team/panel/wiki) 中查看更多使用方法和技巧。

## 问题反馈

使用类问题，可在 [WePublish 社区论坛](https://wepublish.cn/forums) 提问或QQ群`@坤坤`寻求 AI 帮助，亦可在群里寻求付费支持。

面板自身问题，可在 GitHub 的`Issues`页面提交问题反馈，注意[提问的智慧](https://github.com/ryanhanwu/How-To-Ask-Questions-The-Smart-Way/blob/main/README-zh_CN.md)。

## 贡献代码

### 寻找/创建 Issue

您可以在 [Issue 列表](https://github.com/haozi-team/panel/issues) 中寻找或创建一个 Issue，留言表达想要处理该 Issue 的意愿，得到维护者的确认后，即可开始处理。

### 创建 PR

- 在开发过程中，如果遇到问题可以随时在 Issue 中详尽描述该问题，以进一步沟通，但在此之前请确保自己已通过 Google 等方式尽可能的尝试解决问题；
- PR 须提交至我们的极狐 GitLab 仓库[https://jihulab.com/haozi-team/panel](https://jihulab.com/haozi-team/panel)，勿在 GitHub 上提交；
- 当 PR 开发完毕后，请为其添加 `🚀 Review Ready` 标签，维护者将及时进行评审；
- 我们非常欢迎您的贡献，将在下次发版时将您添加到首页贡献者中；❤️

## 赞助商

### 服务器

- [盾云](https://www.ddunyun.com/aff/PNYAXMKI)

### CDN

- [无畏云加速](https://su.sctes.com/register?code=8st689ujpmm2p)

- [又拍云](https://www.upyun.com/?utm_source=lianmeng&utm_medium=referral)

- [AnyCast.Ai](https://www.anycast.ai/)

- [盾云CDN](http://cdn.ddunyun.com/)

### 对象存储

- [又拍云](https://www.upyun.com/?utm_source=lianmeng&utm_medium=referral)

### DevOps

- [极狐 GitLab](https://www.jihulab.com/)

**接受云资源和资金赞助，可通过QQ群咨询联系**

## 贡献者

这个项目的存在要归功于所有做出贡献的人，参与贡献请先查看贡献代码部分。

<a href="https://github.com/DevHaoZi" target="_blank"><img src="https://avatars.githubusercontent.com/u/115467771?v=4" width="48" height="48"></a>

## Star 历史

<a href="https://star-history.com/#haozi-team/panel&Date">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/svg?repos=haozi-team/panel&type=Date&theme=dark" />
    <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/svg?repos=haozi-team/panel&type=Date" />
    <img alt="Star History Chart" src="https://api.star-history.com/svg?repos=haozi-team/panel&type=Date" />
  </picture>
</a>
