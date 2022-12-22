<!--
Name: 文件管理
Author: 耗子
Date: 2022-12-22
-->
<title>文件管理</title>
<div class="layui-fluid">
    <div class="layui-card">
        <div class="layui-card-header">该功能尚未完成开发，请等待！（推荐使用sftp临时代替）</div>
        <div class="layui-card-body">
            <div class="layui-fluid">
                <div id="files" lay-filter="files"></div>
            </div>
        </div>
    </div>
</div>
<script>
    layui.use(['layer', 'admin', 'view', 'form', 'file'], function () {
        let $ = layui.jquery;
        let admin = layui.admin;
        let view = layui.view;
        let form = layui.form;
        let file = layui.file;

        let path = '/www/wwwroot';// 当前路径

        let router = layui.router();
        if (router.search.path) {
            path = window.atob(router.search.path);
        }

        file.render({
            elem: '#files'
            , method: 'get'
            , id: 'files-list'
            , btn_upload: true
            , btn_create: true
            , url: '/api/panel/file/getList'
            , thumb: {'nopic': '/panel/style/ico/null.jpg', 'width': 100, 'height': 100}
            , icon_url: '/panel/style/ico/'
            , done: function (res, curr, count) {
                // console.log(res,curr,count)
            }
            , page: {limit: 10}
            , where: {path: '/www/wwwroot'}
        });
    });
</script>
