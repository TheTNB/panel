<title>文件</title>

<div class="layui-fluid">
    <div class="layui-row layui-col-space15">
        <div class="layui-col-md12">
            <div class="layui-card">
                <div class="layui-card-header">文件（Beta版）</div>
                <div class="layui-card-body">
                    <iframe id="panel_fm" src="{{ asset('../../api/fm') }}" style="width: 100%; height: 800px; border: none;"></iframe>
                </div>
            </div>
        </div>
    </div>
</div>
<script>
layui.use(['jquery'], function () {
    var $ = layui.jquery;
    // 获取iframe的src
    var src = $('#panel_fm').attr('src');
    // src后面加上access_token
    if (layui.data('haozi_panel').access_token !== undefined) {
        $('#panel_fm').attr('src', src + '?access_token=' + layui.data('haozi_panel').access_token);
    }
});
</script>