<!--
Name: 网站 - 列表
Author: 耗子
Date: 2022-11-28
-->
<title>网站</title>
<div class="layui-fluid">
    <div class="layui-row layui-col-space15">
        <div class="layui-col-md12">
            <div class="layui-card">
                <div class="layui-card-header">网站列表</div>
                <div class="layui-card-body">
                    <table class="layui-hide" id="website-list" lay-filter="website-list"></table>

                    <!-- 顶部工具栏 -->
                    <script type="text/html" id="website-list-bar">
                        <div class="layui-btn-container">
                            <button class="layui-btn layui-btn-sm" lay-event="website_add">添加网站</button>
                            <button class="layui-btn layui-btn-sm" lay-event="website_default_settings">全局设置</button>
                        </div>
                    </script>
                    <!-- 右侧网站设置和删除网站 -->
                    <script type="text/html" id="website-setting">
                        <a class="layui-btn layui-btn-xs" lay-event="edit">设置</a>
                        <a class="layui-btn layui-btn-danger layui-btn-xs" lay-event="del">删除</a>
                    </script>
                    <!-- 网站运行状态开关 -->
                    <script type="text/html" id="website-run">
                        <input type="checkbox" name="run" lay-skin="switch" lay-text="ON|OFF"
                               lay-filter="website-run-checkbox"
                               value="@{{ d.status }}" data-website-name="@{{ d.name }}"
                               @{{ d.status== 1 ? 'checked' : '' }} />
                    </script>
                    <!-- 网站SSL状态 -->
                    <script type="text/html" id="website-ssl">
                        @{{ d.ssl == 1 ? '已开启' : '未开启' }}
                    </script>
                </div>
            </div>
        </div>
    </div>
</div>

<script>
    let db_version, php_version;
    layui.use(['admin', 'table', 'form', 'view'], function () {
        var admin = layui.admin
            , table = layui.table
            , form = layui.form
            , view = layui.view;

        // 获取已安装的PHP和DB版本
        admin.req({
            url: "/api/panel/info/getInstalledDbAndPhp"
            , method: 'get'
            , success: function (result) {
                if (result.code !== 0) {
                    console.log('耗子Linux面板：已安装的PHP和DB版本获取失败，接口返回' + result);
                    layer.msg('已安装的PHP和DB版本获取失败，请刷新重试！')
                    return false;
                }
                db_version = result.data.db_version;
                php_version = result.data.php_version;
            }
            , error: function (xhr, status, error) {
                layer.open({
                    title: '错误'
                    , icon: 2
                    , content: '已安装的PHP和DB版本获取失败，接口返回' + xhr.status + ' ' + xhr.statusText
                });
                console.log('耗子Linux面板：ajax请求出错，错误' + error);
            }
        });
        table.render({
            elem: '#website-list'
            , url: '/api/panel/website/getList'
            , toolbar: '#website-list-bar'
            , title: '网站列表'
            , cols: [[
                {field: 'name', title: '网站名', width: 200, fixed: 'left', unresize: true, sort: true}
                , {field: 'run', title: '运行', width: 100, templet: '#website-run', unresize: true}
                , {field: 'path', title: '目录', width: 250}
                , {field: 'php', title: 'PHP', width: 60}
                , {field: 'ssl', title: 'SSL', width: 110, templet: '#website-ssl'}
                , {field: 'note', title: '备注', edit: 'textarea'}
                , {fixed: 'right', title: '操作', toolbar: '#website-setting', width: 150}
            ]]
            /**
             * TODO: 分页
             */
            //, page: true
        });

        // 头工具栏事件
        table.on('toolbar(website-list)', function (obj) {
            if (obj.event === 'website_add') {
                admin.popup({
                    title: '添加网站'
                    , area: ['70%', '60%']
                    , id: 'LAY-popup-website-add'
                    , success: function (layer, index) {
                        view(this.id).render('website/add', {
                            db_version: db_version,
                            php_version: php_version
                        }).done(function () {
                            form.render(null, 'LAY-popup-website-add');
                        });
                    }
                });
            } else if (obj.event === 'website_default_settings') {
                admin.popup({
                    title: '全局设置'
                    , area: ['70%', '60%']
                    , id: 'LAY-popup-website-add'
                    , success: function (layer, index) {
                        view(this.id).render('website/default_settings', {
                        }).done(function () {
                            form.render(null, 'LAY-popup-website-default-settings');
                        });
                    }
                });
            }
        });

        // 行工具事件
        table.on('tool(website-list)', function (obj) {
            let data = obj.data;
            if (obj.event === 'del') {
                layer.confirm('删除网站将一并删除站点目录（不包括数据库），是否继续？', function (index) {
                    admin.req({
                        url: "/api/panel/website/delete"
                        , method: 'post'
                        , data: data
                        , success: function (result) {
                            if (result.code !== 0) {
                                console.log('耗子Linux面板：网站删除失败，接口返回' + result);
                                layer.msg('网站删除失败，请刷新重试！')
                                return false;
                            }
                            obj.del();
                            layer.alert('网站' + data.name + '删除成功！');
                        }
                        , error: function (xhr, status, error) {
                            console.log('耗子Linux面板：ajax请求出错，错误' + error);
                        }
                    });
                    layer.close(index);
                });
            } else if (obj.event === 'edit') {
                let config;

                admin.req({
                    url: "/api/panel/website/getSiteSettings?name=" + data.name
                    , method: 'get'
                    , beforeSend: function (request) {
                        layer.load();
                    }
                    , success: function (result) {
                        if (result.code !== 0) {
                            console.log('耗子Linux面板：网站设置获取失败，接口返回' + result);
                            layer.alert('网站设置获取失败！');
                            return false;
                        }
                        config = result.data;
                        layer.closeAll('loading');
                        // 打开编辑网站页面
                        admin.popup({
                            title: '编辑网站 - ' + data.name
                            , area: ['70%', '80%']
                            , id: 'LAY-popup-website-edit'
                            , success: function (layero, index) {
                                view(this.id).render('website/edit', {
                                    db_version: db_version,
                                    php_version: php_version,
                                    data: data,
                                    config: config
                                }).done(function () {
                                    form.render(null, 'LAY-popup-website-edit');
                                });
                            }
                        });
                    }
                    , error: function (xhr, status, error) {
                        console.log('耗子Linux面板：ajax请求出错，错误' + error);
                    }
                });
            }
        });

        // 网站备注编辑
        table.on('edit(website-list)', function (obj) {
            var value = obj.value // 得到修改后的值
                , data = obj.data; // 得到行数据
            admin.req({
                url: "/api/panel/website/updateSiteNote"
                , method: 'post'
                , data: {
                    name: data.name,
                    note: value
                }
                , success: function (result) {
                    if (result.code !== 0) {
                        console.log('耗子Linux面板：网站备注更新失败，接口返回' + result);
                        layer.msg('网站备注更新失败，请刷新重试！')
                        return false;
                    }
                    layer.alert('网站 ' + data.name + ' 备注更新成功！');
                }
                , error: function (xhr, status, error) {
                    console.log('耗子Linux面板：ajax请求出错，错误' + error);
                }
            });
        });

        // 网站运行状态操作
        form.on('switch(website-run-checkbox)', function (obj) {
            let $ = layui.$;
            let website_name = $(this).data('website-name');
            let run = obj.elem.checked ? 1 : 0;

            //console.log(website_name); //当前行数据
            layer.msg('待开发功能！', {icon: 2});
        });

    });
</script>