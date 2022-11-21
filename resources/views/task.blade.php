<!--
Name: 任务中心
Author: 耗子
Date: 2022-10-09
-->

<title>任务中心</title>
<div class="layui-fluid" id="component-tabs">
    <div class="layui-row">
        <div class="layui-col-md12">
            <div class="layui-card">
                <div class="layui-card-header">任务列表</div>
                <div class="layui-card-body">
                    <div class="layui-tab">
                        <ul class="layui-tab-title">
                            <li class="layui-this">进行中</li>
                            <li>等待中</li>
                            <li>已完成</li>
                        </ul>
                        <div class="layui-tab-content">
                            <div class="layui-tab-item layui-show">
                                <script type="text/html" template lay-url="/api/panel/task/getListRunning"
                                        lay-done="layui.data.done(d);">
                                    @{{# if(d.data != ""){ }}
                                    <blockquote class="layui-elem-quote">安装 @{{ d.data.name }}</blockquote>
                                    <pre id="plugin-install-log" class="layui-code">
                                        日志获取中...
                                    </pre>
                                    @{{# } else { }}
                                    <blockquote class="layui-elem-quote">暂无任务</blockquote>
                                    @{{# } }}
                                </script>
                            </div>
                            <div class="layui-tab-item">
                                <table id="panel-task-waiting" lay-filter="panel-task-waiting"></table>
                            </div>
                            <div class="layui-tab-item">
                                <table id="panel-task-finished" lay-filter="panel-task-finished"></table>
                                <script type="text/html" id="panel-task-finished-control-tpl">
                                    <a class="layui-btn layui-btn-xs" lay-event="remove">移除</a>
                                </script>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</div>

<script>
    var install_plugin_timer;
    function render_plugin_install_log(d) {
        let admin = layui.admin
            ,$ = layui.$
            ,code = layui.code;
        admin.req({
            url: "/api/panel/task/getTaskLog?name=" + d.data.name
            , method: 'get'
            , success: function (result) {
                if (result.code !== 0) {
                    console.log('耗子Linux面板：实时安装日志获取失败，接口返回' + result);
                    $('#plugin-install-log').html('实时安装日志获取失败，请刷新重试！');
                    code({
                        title: 'install.log'
                        , encode: true
                        , about: false

                    });
                    return false;
                }
                $('#plugin-install-log').html(result.data);
                code({
                    title: 'install.log'
                    , encode: true
                    , about: false

                });
            }
            , error: function (xhr, status, error) {
                console.log('耗子Linux面板：ajax请求出错，错误' + error);
            }
        })
    }
    layui.data.done = function (d) {
        // 判断是否有任务
        if (d.data.name) {
            render_plugin_install_log(d);
            // 清除定时器
            clearInterval(install_plugin_timer);
            // 设置定时器，3s一次查询实时安装日志
            install_plugin_timer = setInterval(function () {
                render_plugin_install_log(d);
            }, 3000);
        }

    };

    layui.use(['admin', 'table', 'jquery'], function () {
        var $ = layui.$
            , form = layui.form
            , table = layui.table
            , admin = layui.admin;

        table.render({
            elem: '#panel-task-waiting'
            , url: '/api/panel/task/getListWaiting'
            , cols: [[
                {field: 'slug', hide: true, title: 'Slug', sort: true}
                , {field: 'name', width: '100%', title: '任务名'}
            ]]
            , page: false
            , text: '耗子Linux面板：数据加载出现异常！'
            , done: function () {
                //element.render('progress');
            }
        });

        table.render({
            elem: '#panel-task-finished'
            ,id: 'panel-task-finished-table'
            , url: '/api/panel/task/getListFinished'
            , cols: [[
                {field: 'slug', hide: true, title: 'Slug', sort: true}
                , {field: 'name', width: '80%', title: '任务名'}
                , {field: 'control', title: '操作', templet: '#panel-task-finished-control-tpl', fixed: 'right', align: 'center'}
            ]]
            , page: false
            , text: '耗子Linux面板：数据加载出现异常！'
            , done: function () {
                //element.render('progress');
            }
        });

        //工具条
        table.on('tool(panel-task-finished)', function (obj) {
            let data = obj.data;
            if (obj.event === 'remove') {
                layer.confirm('确定移除该记录吗？', function (index) {
                    layer.close(index);
                    admin.req({
                        url: '/api/panel/task/deleteTask',
                        type: 'post',
                        data: {
                            name: data.name
                        }
                        , success: function (res) {
                            if (res.code == 0) {
                                layer.msg('移除任务：' + data.name + ' 成功！', {icon: 1, time: 1000}, function () {
                                    // 重载表格
                                    table.reload('panel-task-finished-table');
                                });
                            } else {
                                layer.msg(res.msg, {icon: 2, time: 1000});
                            }
                        }
                        , error: function (xhr, status, error) {
                            console.log('耗子Linux面板：ajax请求出错，错误' + error);
                        }
                    });
                });
            }
        });
    });
</script>