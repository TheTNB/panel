<h1 align="center">HaoZi Linux Panel</h1>

<p align="center">
  <a href="https://github.com/haozi-team/panel/releases"><img src="https://img.shields.io/github/release/haozi-team/panel.svg"></a>
  <a href="https://github.com/haozi-team/panel/actions"><img src="https://github.com/haozi-team/panel/actions/workflows/test.yml/badge.svg"></a>
  <a href="https://goreportcard.com/report/github.com/haozi-team/panel"><img src="https://goreportcard.com/badge/github.com/haozi-team/panel"></a>
  <a href="https://codecov.io/gh/haozi-team/panel"><img src="https://codecov.io/gh/haozi-team/panel/branch/main/graph/badge.svg?token=XFT5NGNSRG"></a>
  <a href="https://img.shields.io/github/license/haozi-team/panel"><img src="https://img.shields.io/github/license/haozi-team/panel"></a>
  <a href="https://app.fossa.com/projects/git%2Bgithub.com%2Fhaozi-team%2Fpanel?ref=badge_shield"><img src="https://app.fossa.com/api/projects/git%2Bgithub.com%2Fhaozi-team%2Fpanel.svg?type=shield" alt="FOSSA Status"></a>
</p>

<p align="center">
[<a href="README.md">简体中文</a>] | [English]
</p>

The HaoZi Linux Panel is a lightweight Linux server operation and maintenance management panel developed using Golang and Vue. It is open source under the GNU Affero General Public License v3.0 protocol.

Communication QQ group: [12370907](https://jq.qq.com/?_wv=1027&k=I1oJKSTH) | QQ channel: [pd.qq.com/s/fyol46wfy](https://pd.qq.com/s/fyol46wfy)

## Advantages

1. Extremely low resource usage: Deploying the panel + LNMP environment under Debian 12, the memory usage is less than 400 MB, far ahead of other panels using containerization.
2. Efficient operation and maintenance: The panel UI interface is simple and easy to operate, and you can quickly deploy various environments and adjust application settings without complicated configuration.
3. Support offline operation: The panel can run without relying on any external services. You can even stop the panel process after deployment is complete, and it will not affect the deployed services.
4. Tested by time: We have been using it in production environment since 2022, and it has been running stably for 1 year without any accidents.
5. Open source and open: The panel is open source, you can freely modify and audit the panel source code, and the security is guaranteed.

## UI Screenshots

![UI Screenshots](ui.png)

## Operating Environment

HaoZi Linux Panel only supports the latest version of mainstream systems under the `amd64` | `arm64` architecture. It does not support `Ubuntu` because its releases are too frequent and difficult to maintain.

Recommended to use `Debian` for low-configuration machines, as its resource usage is lower than that of the `RHEL` system. For other machines, recommended to use `RockyLinux` | `AlmaLinux`, which has a longer maintenance cycle and is more stable.

For other systems not in the table below (OpenCloudOS 8, Anolis 8, CentOS Stream 8/9, Debian 11, etc.), you can try to install it yourself, but normal operation is not guaranteed, and free technical support is not provided (theoretically there will be no major question).

CentOS Stream can be migrated to a supported system using the migration script: [CentOS 8/9 Migration Script](https://github.com/haozi-team/byecentos)

| OS         | Version |
|------------|---------|
| RHEL       | 9       |
| RockyLinux | 9       |
| AlmaLinux  | 9       |
| Debian     | 12      |

As system versions are constantly updated, we may also terminate support for some older systems to ensure the stability of the panel.

## Install Panel

Before installing the panel, you need to understand the basic knowledge of the LNMP environment and how to deal with common LNMP environment problems. It is not recommended for users with zero basic knowledge to install and use HaoZi Linux Panel ([Recommended: aaPanel - easy-to-use server operation and maintenance panel](https://aapanel.com)).

If you decide to continue, please log in to the server as `root` user and execute the following command to install the panel:

```shell
HAOZI_DL_URL="https://jihulab.com/haozi-team/download/-/raw/main/panel"; curl -sSL -O ${HAOZI_DL_URL}/install_panel.sh && curl -sSL -O ${HAOZI_DL_URL}/install_panel.sh.checksum.txt && sha256sum -c install_panel.sh.checksum.txt && bash install_panel.sh || echo "Checksum Verification Failed, File May Have Been Tampered With, Operation Terminated"
```

## Uninstall Panel

Recommended to back up data and reinstall the system first, so that the system can be kept clean.

If you are unable to reinstall the system, log in to the server as the `root` user and execute the following command to uninstall the panel:

```shell
HAOZI_DL_URL="https://jihulab.com/haozi-team/download/-/raw/main/panel"; curl -sSL -O ${HAOZI_DL_URL}/uninstall_panel.sh && curl -sSL -O ${HAOZI_DL_URL}/uninstall_panel.sh.checksum.txt && sha256sum -c uninstall_panel.sh.checksum.txt && bash uninstall_panel.sh || echo "Checksum Verification Failed, File May Have Been Tampered With, Operation Terminated"
```

Before uninstalling the panel, please be sure to back up all data and uninstall all panel plugins in advance. The data will **not be recoverable** after uninstallation!

## Daily Maintenance

Use `panel` command for daily maintenance:

```shell
panel
```

See more usage methods and tips in [Wiki](https://github.com/haozi-team/panel/wiki).

## Feedback

For usage issues, you can ask questions in the [WePublish Community Forum](https://wepublish.cn/forums) or seek AI help in the QQ group `@坤坤`. You can also seek paid support in the group.

For issues with the panel itself, you can submit feedback on the `Issues` page of GitHub. Please note [the wisdom of asking questions](http://www.catb.org/~esr/faqs/smart-questions.html).

## Contribute code

### Find/Create Issue

You can find or create an Issue in [Issue List](https://github.com/haozi-team/panel/issues), leave a message to express your willingness to deal with the Issue, and get confirmation from the maintainer. Start processing.

### Create PR

- During the development process, if you encounter a problem, you can describe the problem in detail in the Issue at any time for further communication, but before doing so, please make sure that you have tried your best to solve the problem through Google and other methods;
- PRs must be submitted to our JiHu GitLab repository [https://jihulab.com/haozi-team/panel](https://jihulab.com/haozi-team/panel). Do not submit on GitHub;
- When the PR is completed, please add the `🚀 Review Ready` tag to it, and the maintainer will review it in time;
- We very much welcome your contributions and will add you to the homepage contributors in the next release; ❤️

## Sponsor

### CDN

- [无畏云加速](https://su.sctes.com/register?code=8st689ujpmm2p)

- [盾云CDN](http://cdn.ddunyun.com/)

### Object Storage

- [又拍云](https://www.upyun.com/?utm_source=lianmeng&utm_medium=referral)

### DevOps

- [JiHu GitLab](https://www.jihulab.com/)

**Accept cloud resources and financial sponsorship, you can contact us through QQ group**

## Contributor

This project owes its existence to all those who have contributed. To contribute, please check the contributed code section first.

<a href="https://github.com/DevHaoZi" target="_blank"><img src="https://avatars.githubusercontent.com/u/115467771?v=4" width="48" height="48"></a>

## Star History

<a href="https://star-history.com/#haozi-team/panel&Date">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/svg?repos=haozi-team/panel&type=Date&theme=dark" />
    <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/svg?repos=haozi-team/panel&type=Date" />
    <img alt="Star History Chart" src="https://api.star-history.com/svg?repos=haozi-team/panel&type=Date" />
  </picture>
</a>

## Disclaimer

It is strictly prohibited to use the HaoZi Linux Panel to engage in any illegal activities. Please do not request any form of technical support from us for illegal sites. If illegal content is discovered during the technical support process, we will immediately stop technical support and retain relevant evidence.
