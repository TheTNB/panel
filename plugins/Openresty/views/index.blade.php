<!--
Name: Openresty管理器
Author: 耗子
Date: 2022-11-02
-->
<title>OpenResty</title>
<div class="layui-fluid" id="component-tabs">
    <div class="layui-row">
        <div class="layui-col-md12">
            <div class="layui-card">
                <div class="layui-card-header">OpenResty管理</div>
                <div class="layui-card-body">
                    <div class="layui-tab">
                        <ul class="layui-tab-title">
                            <li class="layui-this">运行状态</li>
                            <li>配置修改</li>
                            <li>负载状态</li>
                            <li>错误日志</li>
                        </ul>
                        <div class="layui-tab-content">
                            <div class="layui-tab-item layui-show">
                                <blockquote id="openresty-status" class="layui-elem-quote layui-quote-nm">当前状态：<span
                                        class="layui-badge layui-bg-black">获取中</span></blockquote>
                                <div class="layui-btn-container" style="padding-top: 30px;">
                                    <button id="openresty-start" class="layui-btn">启动</button>
                                    <button id="openresty-stop" class="layui-btn layui-btn-danger">停止</button>
                                    <button id="openresty-restart" class="layui-btn layui-btn-warm">重启</button>
                                    <button id="openresty-reload" class="layui-btn layui-btn-normal">重载</button>
                                </div>
                            </div>
                            <div class="layui-tab-item">
                                <blockquote class="layui-elem-quote">此处修改的是OpenResty主配置文件，如果你不了解各参数的含义，请不要随意修改！<br>
                                    提示：Ctrl+F 搜索关键字，Ctrl+S 保存，Ctrl+H 查找替换！
                                </blockquote>
                                <div id="openresty-config-editor"
                                     style="width: -webkit-fill-available; height: 600px;"></div>
                                <div class="layui-btn-container" style="padding-top: 30px;">
                                    <button id="openresty-config-save" class="layui-btn">保存</button>
                                </div>
                            </div>
                            <div class="layui-tab-item">
                                <table class="layui-hide" id="openresty-load-status"></table>
                            </div>
                            <div class="layui-tab-item">
                                <div class="layui-btn-container">
                                    <button id="openresty-clean-error-log" class="layui-btn">清空日志</button>
                                </div>
                                <pre id="openresty-error-log" class="layui-code">
                                    获取中...
                                </pre>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</div>

<script>
    let openresty_config_editor;// 定义openresty配置编辑器的全局变量
    layui.use(['index', 'code', 'table'], function () {
        let $ = layui.$
            , admin = layui.admin
            , element = layui.element
            , code = layui.code
            , table = layui.table;

        // 获取openresty运行状态并渲染
        admin.req({
            url: "/api/plugin/openresty/status"
            , method: 'get'
            , success: function (result) {
                if (result.code !== 0) {
                    console.log('耗子Linux面板：OpenResty运行状态获取失败，接口返回' + result);
                    return false;
                }
                if (result.data === "running") {
                    $('#openresty-status').html('当前状态：<span class="layui-badge layui-bg-green">运行中</span>');
                } else {
                    $('#openresty-status').html('当前状态：<span class="layui-badge layui-bg-red">已停止</span>');
                }

            }
            , error: function (xhr, status, error) {
                console.log('耗子Linux面板：ajax请求出错，错误' + error)
            }
        });

        // 获取openresty错误日志并渲染
        admin.req({
            url: "/api/plugin/openresty/errorLog"
            , method: 'get'
            , success: function (result) {
                if (result.code !== 0) {
                    console.log('耗子Linux面板：OpenResty错误日志获取失败，接口返回' + result);
                    $('#openresty-error-log').text('OpenResty错误日志获取失败，请刷新重试！');
                    code({
                        title: 'error.log'
                        , encode: true
                        , about: false

                    });
                    return false;
                }
                $('#openresty-error-log').text(result.data);
                code({
                    title: 'error.log'
                    , encode: true
                    , about: false

                });
            }
            , error: function (xhr, status, error) {
                console.log('耗子Linux面板：ajax请求出错，错误' + error)
            }
        });

        // 获取openresty配置并渲染
        admin.req({
            url: "/api/plugin/openresty/config"
            , method: 'get'
            , success: function (result) {
                if (result.code !== 0) {
                    console.log('耗子Linux面板：OpenResty主配置获取失败，接口返回' + result);
                    return false;
                }
                $('#openresty-config-editor').text(result.data);
                openresty_config_editor = ace.edit("openresty-config-editor", {
                    mode: "ace/mode/nginx",
                    selectionStyle: "text"
                });
            }
            , error: function (xhr, status, error) {
                console.log('耗子Linux面板：ajax请求出错，错误' + error)
            }
        });

        // 获取openresty负载状态并渲染
        table.render({
            elem: '#openresty-load-status'
            , url: '/api/plugin/openresty/load'
            , cols: [[
                {field: 'name', width: '80%', title: '属性',}
                , {field: 'value', width: '20%', title: '当前值'}
            ]]
        });
        element.render();

        // 事件监听
        $('#openresty-start').click(function () {
            admin.popup({
                title: '<span style="color: red;">警告</span>'
                ,
                shade: 0
                ,
                anim: -1
                ,
                area: ['300px', '200px']
                ,
                id: 'layadmin-layer-skin-openresty-start'
                ,
                skin: 'layui-anim layui-anim-upbit'
                ,
                content: '面板的正常访问依赖OpenResty，因此不支持在面板启动OpenResty，如您确需操作，请在SSH执行<span class="layui-badge-rim">systemctl start nginx</span>以启动OpenResty！'
            });
        });
        $('#openresty-stop').click(function () {
            admin.popup({
                title: '<span style="color: red;">警告</span>'
                ,
                shade: 0
                ,
                anim: -1
                ,
                area: ['300px', '200px']
                ,
                id: 'layadmin-layer-skin-openresty-stop'
                ,
                skin: 'layui-anim layui-anim-upbit'
                ,
                content: '面板的正常访问依赖OpenResty，因此不支持在面板停止OpenResty，如您确需操作，请在SSH执行<span class="layui-badge-rim">systemctl stop nginx</span>以停止OpenResty！'
            });
        });
        $('#openresty-restart').click(function () {
            layer.confirm('重启OpenResty有可能导致面板短时间无法访问，是否继续重启？', {
                btn: ['重启', '取消']
            }, function () {
                admin.req({
                    url: "/api/plugin/openresty/restart"
                    , method: 'get'
                    , beforeSend: function () {
                        layer.msg('已发送重启请求，请稍后刷新确认重启状态。');
                    }
                    , success: function (result) {
                        if (result.code !== 0) {
                            console.log('耗子Linux面板：OpenResty重启失败，接口返回' + result);
                            return false;
                        }
                        if (result.msg === 'error') {
                            layer.alert(result.data);
                            return false;
                        }
                        layer.alert('OpenResty重启成功！');
                        admin.events.refresh();
                    }
                    , error: function (xhr, status, error) {
                        console.log('耗子Linux面板：ajax请求出错，错误' + error)
                    }
                });
            }, function () {
                layer.msg('取消重启');
            });
        });
        $('#openresty-reload').click(function () {
            layer.msg('OpenResty重载中...');
            admin.req({
                url: "/api/plugin/openresty/reload"
                , method: 'get'
                , success: function (result) {
                    if (result.code !== 0) {
                        console.log('耗子Linux面板：OpenResty重载失败，接口返回' + result);
                        return false;
                    }
                    if (result.msg === 'error') {
                        layer.alert(result.data);
                        return false;
                    }
                    layer.alert('OpenResty重载成功！');
                }
                , error: function (xhr, status, error) {
                    console.log('耗子Linux面板：ajax请求出错，错误' + error)
                }
            });
        });
        $('#openresty-config-save').click(function () {
            layer.msg('OpenResty配置保存中...');
            admin.req({
                url: "/api/plugin/openresty/config"
                , method: 'post'
                , data: {
                    config: openresty_config_editor.getValue()
                }
                , success: function (result) {
                    if (result.code !== 0) {
                        console.log('耗子Linux面板：OpenResty配置保存失败，接口返回' + result);
                        return false;
                    }
                    if (result.msg === 'error') {
                        layer.alert(result.data);
                        return false;
                    }
                    layer.alert('OpenResty配置保存成功！');
                }
                , error: function (xhr, status, error) {
                    console.log('耗子Linux面板：ajax请求出错，错误' + error)
                }
            });
        });
        $('#openresty-clean-error-log').click(function () {
            layer.msg('错误日志清空中...');
            admin.req({
                url: "/api/plugin/openresty/cleanErrorLog"
                , method: 'get'
                , success: function (result) {
                    if (result.code !== 0) {
                        console.log('耗子Linux面板：OpenResty错误日志清空失败，接口返回' + result);
                        return false;
                    }
                    layer.msg('OpenResty错误日志已清空！');
                    setTimeout(function () {
                        admin.events.refresh();
                    }, 1000);
                }
                , error: function (xhr, status, error) {
                    console.log('耗子Linux面板：ajax请求出错，错误' + error)
                }
            });
        });
    });
</script>