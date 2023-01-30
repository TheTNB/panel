<!--
Name: 主页模板
Author: 耗子
Date: 2022-11-30
-->
<div class="layui-fluid">
    <div class="layui-row layui-col-space15">
        <div id="address1" class="layui-col-md12">
            <div class="layui-collapse">
                <div style="background: #fff;" class="layui-colla-content layui-show">
                    <div class="text" style="overflow: hidden;height: 22px;">
                        <div class="layui-carousel" id="home_ad" lay-filter="home_ad">
                            <div carousel-item="">
                                <a style="background: #fff;" href="https://hzbk.net"
                                   title="耗子博客" target="_blank"><i class="layui-icon layui-icon-release"></i> 耗子博客</a>
                                <a style="background: #fff;" href="https://weavatar.com"
                                   title="WeAvatar" target="_blank"><i class="layui-icon layui-icon-release"></i>
                                    WeAvatar - 互联网公共头像服务</a>
                                <a style="background: #fff;" href="https://wepublish.cn"
                                   title="WePublish" target="_blank"><i class="layui-icon layui-icon-release"></i>
                                    WePublish - WordPress的本土化版本</a>
                            </div>
                        </div>
                    </div>

                </div>
            </div>

        </div>
        <div id="monitor1" class="layui-col-md6">

            <div class="layui-card">
                <div class="layui-card-header">资源使用</div>
                <div class="layui-card-body layadmin-takerates">
                    <div class="layui-progress" lay-showPercent="yes" lay-filter="home_cpu">
                        <h3 id="home_cpu">CPU信息加载中</h3>
                        <div class="layui-progress-bar" lay-percent="0%"></div>
                    </div>
                    <div class="layui-progress" lay-showPercent="yes" lay-filter="home_mem">
                        <h3 id="home_mem">内存信息加载中</h3>
                        <div class="layui-progress-bar layui-bg-red" lay-percent="0%"></div>
                    </div>
                </div>
            </div>
        </div>
        <div id="monitor2" class="layui-col-md3">

            <div class="layui-card">
                <div class="layui-card-header">系统负载</div>
                <div class="layui-card-body layadmin-takerates">
                    <div class="layui-progress" lay-showPercent="yes" lay-filter="uptime_1">
                        <h3>近1分钟</h3>
                        <div class="layui-progress-bar" lay-percent="0%"></div>
                    </div>
                    <div class="layui-progress" lay-showPercent="yes" lay-filter="uptime_5">
                        <h3>近5分钟</h3>
                        <div class="layui-progress-bar layui-bg-red" lay-percent="0%"></div>
                    </div>
                </div>
            </div>
        </div>
        <div class="layui-col-md3">
            <div class="layui-card">
                <div class="layui-card-header">
                    实时流量
                    <span class="layui-badge layui-bg-blue layuiadmin-badge">发送 / 接收</span>
                </div>
                <div class="layui-card-body layuiadmin-card-list">
                    <p style="text-align: center;"><b id="home_net_now">获取中...</b></p>
                </div>
                <div class="layui-card-header">
                    累计流量
                    <span class="layui-badge layui-bg-blue layuiadmin-badge">发送 / 接收</span>
                </div>
                <div class="layui-card-body layuiadmin-card-list">
                    <p style="text-align: center;"><b id="home_net_total">获取中...</b></p>
                </div>
            </div>
        </div>
        <div class="layui-col-md8">
            <div class="layui-card">
                <div class="layui-card-header">应用</div>
                <div class="layui-card-body">
                    <div class="layui-carousel layadmin-carousel layadmin-shortcut" lay-anim="">
                        <div carousel-item="">
                            <ul class="layui-row layui-col-space10 layui-this">
                                <script type="text/html" template lay-url="/api/panel/info/getHomePlugins">
                                    @{{#  layui.each(d.data, function(index, item){ }}
                                    <li class="layui-col-xs4 layui-col-md2 layui-col-sm4">
                                        <a lay-href="/plugin/@{{ item.slug }}">
                                            <i class="layui-icon layui-icon-engine"></i>
                                            <cite>@{{ item.name }}</cite>
                                        </a>
                                    </li>
                                    @{{#  }); }}
                                    @{{#  if(d.data.length === 0){ }}
                                    这里好像啥也没有...
                                    @{{#  } }}
                                </script>

                            </ul>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <div class="layui-col-md4">
            <div class="layui-card">
                <div class="layui-card-header">基本信息</div>
                <div class="layui-card-body layui-text layadmin-version">
                    <table class="layui-table">
                        <tbody>
                        <tr>
                            <td>系统信息</td>
                            <td id="home_os_name">
                                获取中...
                            </td>
                        </tr>
                        <tr>
                            <td>面板版本</td>
                            <td id="home_panel_version">
                                获取中...
                            </td>
                        </tr>
                        <tr>
                            <td>运行时间</td>
                            <td id="home_uptime">
                                获取中...
                            </td>
                        </tr>
                        <tr>
                            <td>操作</td>
                            <td>
                                <button id="update_panel" class="layui-btn layui-btn-xs">更新</button>
                            </td>
                        </tr>
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
        <div class="layui-col-md12">
            <div class="layui-card">
                <div class="layui-card-header">
                    关于面板
                    <i class="layui-icon layui-icon-tips" lay-tips="你干嘛，哎哟！" lay-offset="5"></i>
                </div>
                <div class="layui-card-body layui-text layadmin-text">
                    <blockquote class="layui-elem-quote">
                        <p style="color: red;">开发组祝大家2023新年快乐！新的一年永无Bug，永不宕机！</p>
                    </blockquote>
                    <blockquote class="layui-elem-quote">
                        <p>欢迎您使用耗子Linux面板。如遇到问题/Bug，可通过 <a
                                    href="https://jq.qq.com/?_wv=1027&amp;k=I1oJKSTH">Q群12370907</a> / <a
                                    target="_blank" href="https://pd.qq.com/s/fyol46wfy">QQ频道</a> 寻求帮助</p>
                    </blockquote>
                </div>
            </div>
        </div>
    </div>
</div>

<script>
    var home_timer;
    var uptime_1 = '获取中', uptime_5 = '获取中', uptime_15 = '获取中';

    function refresh_home_info() {
        layui.use(['index', 'jquery', 'admin'], function () {
            let $ = layui.jquery
                , admin = layui.admin
                , element = layui.element;

            let device = layui.device();
            let cpu_info;
            admin.req({
                url: "/api/panel/info/getNowMonitor"
                , method: 'get'
                , success: function (result) {
                    if (result.code !== 0) {
                        console.log('耗子Linux面板：系统资源信息获取失败，接口返回' + result);
                        layer.msg('系统资源获取失败，请刷新重试！')
                        return false;
                    }
                    element.progress('home_cpu', result.data.cpu_use);
                    element.progress('home_mem', result.data.mem_use_p);
                    element.progress('uptime_1', result.data.uptime_1_p);
                    element.progress('uptime_5', result.data.uptime_5_p);
                    uptime_1 = result.data.uptime_1;
                    uptime_5 = result.data.uptime_5;
                    uptime_15 = result.data.uptime_15;
                    // 判断一下移动设备不显示CPU型号，放不下。。。
                    if (device.mobile) {
                        cpu_info = result.data.cpu_info.physical + 'CPU ' + result.data.cpu_info.cores + '核心 ' + result.data.cpu_info.siblings + '线程';
                    } else {
                        cpu_info = result.data.cpu_info.name + ' ' + result.data.cpu_info.physical + 'CPU ' + result.data.cpu_info.cores + '核心 ' + result.data.cpu_info.siblings + '线程';
                    }
                    $('#home_net_total').html(result.data.tx_total + ' / ' + result.data.rx_total);
                    $('#home_net_now').html(result.data.tx_now + '/s / ' + result.data.rx_now + '/s');
                    $('#home_cpu').text(cpu_info);
                    $('#home_mem').text('使用' + result.data.mem_use + 'MB / ' + '总计' + result.data.mem_total + 'MB');
                    element.render('progress');
                }
                , error: function (xhr, status, error) {
                    console.log('耗子Linux面板：ajax请求出错，错误' + error);
                }
            });
        });
    }

    // 先执行一次
    refresh_home_info();
    // 然后设置个定时器3s一次刷新
    clearInterval(home_timer);
    home_timer = setInterval(refresh_home_info, 3000);
    // 获取系统信息，这部分信息无需更新。
    layui.use(['index', 'jquery', 'admin', 'carousel'], function () {
        let $ = layui.jquery
            , admin = layui.admin
            , element = layui.element
            , carousel = layui.carousel;
        carousel.render({
            elem: '#home_ad'
            , width: '100%'
            , height: '200px'
            , anim: 'fade'
            , arrow: 'none'
            , indicator: 'none'
        });
        admin.req({
            url: "/api/panel/info/getSystemInfo"
            , method: 'get'
            , success: function (result) {
                if (result.code !== 0) {
                    console.log('耗子Linux面板：系统信息获取失败，接口返回' + result);
                    layer.msg('系统信息获取失败，请刷新重试！')
                    return false;
                }
                $('#home_os_name').text(result.data.os_name);
                $('#home_panel_version').text(result.data.panel_version);
                $('#home_uptime').text('已不间断运行 ' + result.data.uptime + ' 天');
            }
            , error: function (xhr, status, error) {
                console.log('耗子Linux面板：ajax请求出错，错误' + error);
            }
        });
    });

    // 监听鼠标悬停到uptime上的事件
    // 用于显示1分钟、5分钟、15分钟的负载
    layui.use(['jquery', 'layer'], function () {
        let $ = layui.jquery
            , layer = layui.layer
            , admin = layui.admin;
        $('#monitor2').hover(function () {
            layer.tips('1分钟负载：' + uptime_1 + '<br>5分钟负载：' + uptime_5 + '<br>15分钟负载：' + uptime_15, '#monitor2', {
                tips: 1,
                time: 0
            });
        }, function () {
            layer.closeAll('tips');
        });
        // 监听更新按钮点击事件
        $('#update_panel').click(function () {
            index = layer.msg('正在获取版本信息...', {
                icon: 16
                , time: 0
            });
            admin.req(
                {
                    url: '/api/panel/info/checkUpdate'
                    , method: 'get'
                    , success: function (result) {
                        layer.close(index);
                        if (result.code !== 0) {
                            layer.msg('获取版本信息失败，请刷新重试！')
                            return false;
                        }
                        if (result.data.version) {
                            admin.popup({
                                title: '提示'
                                ,
                                shade: 0
                                ,
                                anim: -1
                                ,
                                area: ['400px', '200px']
                                ,
                                id: 'layadmin-layer-skin-update-panel'
                                ,
                                skin: 'layui-anim layui-anim-upbit'
                                ,
                                content: '最新版本：' + result.data.version + '<br><br>更新日志：' + result.data.describe + '<br><br>请在SSH执行<span class="layui-badge-rim">panel update</span>以更新面板！'
                            });
                        } else {
                            layer.msg('当前已是最新版本！')
                        }
                    }
                    , error: function (xhr, status, error) {
                        layer.close(index);
                        layer.msg('获取版本信息失败，请刷新重试！')
                    }
                }
            );
        });
    });
</script>

