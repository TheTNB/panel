<title>资源监控</title>

<div class="layui-fluid">
    <div class="layui-card">
        <div class="layui-form layui-card-header layuiadmin-card-header-auto">
            <div class="layui-inline">
                <span style="margin-right: 10px;">开启监控</span><input type="checkbox" id="monitor-switch"
                                                                        lay-filter="monitor" lay-skin="switch"
                                                                        lay-text="ON|OFF">
                <span style="margin-left: 40px; margin-right: 10px;">保存天数</span>
                <div class="layui-input-inline"><input type="number" name="monitor-save-days" class="layui-input"
                                                       style="height: 30px; margin-top: 5px;" min=0 max=30 disabled>
                </div>
                <div class="layui-input-inline">
                    <button id="save_monitor_date" class="layui-btn layui-btn-sm" style="margin-left: 10px;">确定
                    </button>
                </div>
            </div>
            <div style="float: right;">
                <button id="clear_monitor_record" class="layui-btn layui-btn-sm layui-btn-danger">清空监控记录
                </button>
            </div>
        </div>
    </div>

    <div class="layui-row layui-col-space10">
        <div class="layui-col-xs12 layui-col-md6">
            <div class="layui-card">
                <div class="layui-card-header">
                    <span>负载</span>
                    <div style="float: right;">
                        {{--<button class="layui-btn layui-btn-xs layui-btn-primary">昨天
                        </button>--}}
                        <button class="layui-btn layui-btn-xs">今天</button>
                        {{--<button class="layui-btn layui-btn-xs layui-btn-primary">
                            最近七天
                        </button>
                        <button class="layui-btn layui-btn-xs layui-btn-primary">
                            最近30天
                        </button>
                        <button id="test" class="layui-btn layui-btn-xs layui-btn-primary">自定义时间
                        </button>--}}
                    </div>
                </div>
                <div class="layui-card-body">
                    <div id="load_monitor" style="width: 100%;height: 400px;"></div>
                </div>
            </div>
        </div>

        <div class="layui-col-xs12 layui-col-md6">
            <div class="layui-card">
                <div class="layui-card-header">
                    <span>CPU</span>
                    <div style="float: right;">
                        {{--<button class="layui-btn layui-btn-xs layui-btn-primary">昨天
                        </button>--}}
                        <button class="layui-btn layui-btn-xs">今天</button>
                        {{--<button class="layui-btn layui-btn-xs layui-btn-primary">
                            最近七天
                        </button>
                        <button class="layui-btn layui-btn-xs layui-btn-primary">最近30天
                        </button>
                        <button class="layui-btn layui-btn-xs layui-btn-primary">自定义时间
                        </button>--}}
                    </div>
                </div>
                <div class="layui-card-body">
                    <div id="cpu_monitor" style="width: 100%;height: 400px;"></div>
                </div>
            </div>
        </div>
    </div>

    <div class="layui-row layui-col-space10">
        <div class="layui-col-xs12 layui-col-md6">
            <div class="layui-card">
                <div class="layui-card-header">
                    <span>内存</span>
                    <div style="float: right;">
                        {{--<button class="layui-btn layui-btn-xs layui-btn-primary">昨天
                        </button>--}}
                        <button class="layui-btn layui-btn-xs">今天</button>
                        {{--<button class="layui-btn layui-btn-xs layui-btn-primary">
                            最近七天
                        </button>
                        <button class="layui-btn layui-btn-xs layui-btn-primary">
                            最近30天
                        </button>
                        <button class="layui-btn layui-btn-xs layui-btn-primary">
                            自定义时间
                        </button>--}}
                    </div>
                </div>
                <div class="layui-card-body">
                    <div id="memory_monitor" style="width: 100%;height: 400px;"></div>
                </div>
            </div>
        </div>
        <div class="layui-col-xs12 layui-col-md6">
            <div class="layui-card">
                <div class="layui-card-header">
                    <span>网络</span>
                    <div style="float: right;">
                        {{--<button class="layui-btn layui-btn-xs layui-btn-primary">昨天
                        </button>--}}
                        <button class="layui-btn layui-btn-xs">今天</button>
                        {{--<button class="layui-btn layui-btn-xs layui-btn-primary">最近七天
                        </button>
                        <button class="layui-btn layui-btn-xs layui-btn-primary">最近30天
                        </button>
                        <button class="layui-btn layui-btn-xs layui-btn-primary">自定义时间
                        </button>--}}
                    </div>
                </div>
                <div class="layui-card-body">
                    <div id="network_monitor" style="width: 100%;height: 400px;"></div>
                </div>
            </div>
        </div>
    </div>
</div>
<script>
    layui.use(['admin', 'view', 'form', 'echarts', 'element', 'carousel'], function () {
        var admin = layui.admin;
        var view = layui.view;
        var $ = layui.jquery;
        var form = layui.form;

        // 获取监控开关和保存天数
        admin.req({
            url: '/api/panel/monitor/getMonitorSwitchAndDays',
            type: 'get',
            dataType: 'json',
            success: function (res) {
                if (res.code === 0) {
                    if (res.data.monitor == 1) {
                        $('#monitor-switch').attr('checked', true);
                    } else {
                        $('#monitor-switch').attr('checked', false);
                    }
                    $('input[name="monitor-save-days"]').val(res.data.monitor_days);
                    // 移除禁用
                    $('input[name="monitor-save-days"]').removeAttr('disabled');
                    form.render();
                }
            }
        });
        // 监听switch开关：是否开启监控
        form.on('switch(monitor)', function (data) {
            admin.req({
                url: '/api/panel/monitor/setMonitorSwitch',
                type: 'post',
                dataType: 'json',
                data: {switch: data.elem.checked},
                success: function (res) {
                    if (res.code === 0) {
                        layer.msg('修改成功', {icon: 1});
                    } else {
                        layer.msg(res.msg, {icon: 2});
                    }
                }
            });
        });
        // 监听保存天数按钮
        $('#save_monitor_date').click(function () {
            var days = $('input[name="monitor-save-days"]').val();
            if (days == '') {
                layer.msg('请输入保存天数', {icon: 2});
                return false;
            }
            admin.req({
                url: '/api/panel/monitor/setMonitorSaveDays',
                type: 'post',
                dataType: 'json',
                data: {days: days},
                success: function (res) {
                    if (res.code === 0) {
                        layer.msg('修改成功', {icon: 1});
                    } else {
                        layer.msg(res.msg, {icon: 2});
                    }
                }
            });
        });

        // 监听清除监控数据按钮
        $('#clear_monitor_record').click(function () {
            layer.confirm('确定要清除监控数据吗？', function (index) {
                admin.req({
                    url: '/api/panel/monitor/clearMonitorData',
                    type: 'post',
                    dataType: 'json',
                    success: function (res) {
                        if (res.code === 0) {
                            layer.msg('清除成功', {icon: 1});
                            setTimeout(function () {
                                admin.render();
                            }, 1000);
                        } else {
                            layer.msg(res.msg, {icon: 2});
                        }
                    }
                });
                layer.close(index);
            });
        });

        // 获取监控数据
        admin.req({
            url: '/api/panel/monitor/getMonitorData',
            type: 'get',
            dataType: 'json',
            success: function (res) {
                if (res.code !== 0) {
                    layer.msg(res.msg, {icon: 2});
                    return false;
                }
                let loadChart = renderEcharts('load_monitor', '负载监控', undefined, res.data.times, [{
                    name: '负载',
                    type: 'line',
                    smooth: true,
                    itemStyle: {normal: {areaStyle: {type: 'default'}}},
                    data: res.data.uptime.uptime,
                    markPoint: {
                        data: [{type: 'max', name: '最大值'}, {type: 'min', name: '最小值'}]
                    },
                    markLine: {
                        data: [{type: 'average', name: '平均值'}]
                    }
                }], [{
                    type: 'value',
                }]);
                let cpuChart = renderEcharts('cpu_monitor', 'CPU监控', undefined, res.data.times, [{
                    name: '使用率',
                    type: 'line',
                    smooth: true,
                    itemStyle: {normal: {areaStyle: {type: 'default'}}},
                    data: res.data.cpu.use,
                    markPoint: {
                        data: [{type: 'max', name: '最大值'}, {type: 'min', name: '最小值'}]
                    },
                    markLine: {
                        data: [{type: 'average', name: '平均值'}]
                    }
                }], [{
                    name: '单位 %',
                    min: 0,
                    max: 100,
                    type: 'value',
                    axisLabel: {
                        formatter: '{value} %'
                    }
                }]);
                let memoryChart = renderEcharts('memory_monitor', '内存', {
                    x: 'left',
                    data: ["内存", "Swap"]
                }, res.data.times, [{
                    name: '内存',
                    type: 'line',
                    smooth: true,
                    itemStyle: {normal: {areaStyle: {type: 'default'}}},
                    data: res.data.memory.mem_use,
                    markPoint: {
                        data: [{type: 'max', name: '最大值'}, {type: 'min', name: '最小值'}]
                    },
                    markLine: {
                        data: [{type: 'average', name: '平均值'}]
                    }
                }, {
                    name: 'Swap',
                    type: 'line',
                    smooth: true,
                    itemStyle: {normal: {areaStyle: {type: 'default'}}},
                    data: res.data.memory.swap_use,
                    markPoint: {
                        data: [{type: 'max', name: '最大值'}, {type: 'min', name: '最小值'}]
                    },
                    markLine: {
                        data: [{type: 'average', name: '平均值'}]
                    }
                }], [{
                    name: '单位 MB',
                    min: 0,
                    max: res.data.mem_total,
                    type: 'value',
                    axisLabel: {
                        formatter: '{value} M'
                    }
                }]);
                let networkChart = renderEcharts('network_monitor', '网络', {
                    x: 'left',
                    data: ["出", "入"]
                }, res.data.times, [{
                    name: '出',
                    type: 'line',
                    smooth: true,
                    itemStyle: {normal: {areaStyle: {type: 'default'}}},
                    data: res.data.network.tx_now,
                    markPoint: {
                        data: [{type: 'max', name: '最大值'}, {type: 'min', name: '最小值'}]
                    },
                    markLine: {
                        data: [{type: 'average', name: '平均值'}]
                    }
                }, {
                    name: '入',
                    type: 'line',
                    smooth: true,
                    itemStyle: {normal: {areaStyle: {type: 'default'}}},
                    data: res.data.network.rx_now,
                    markPoint: {
                        data: [{type: 'max', name: '最大值'}, {type: 'min', name: '最小值'}]
                    },
                    markLine: {
                        data: [{type: 'average', name: '平均值'}]
                    }
                }], [{
                    name: '单位 Kb/s',
                    type: 'value',
                    axisLabel: {
                        formatter: '{value} Kb'
                    }
                }]);

                // 在窗口大小改变时，重置图表大小
                window.addEventListener("resize", function () {
                    loadChart.resize();
                    cpuChart.resize();
                    memoryChart.resize();
                    networkChart.resize();
                });
            }, error: function (xhr, status, error) {
                console.log('耗子Linux面板：ajax请求出错，错误' + error);
            }
        });
    });

    // 渲染图表
    function renderEcharts(element_id, title, legend = undefined, data_xAxis, series, yAxis = undefined) {
        var Chart = echarts.init(document.getElementById(element_id), layui.echartsTheme);
        var option = {
            title: {text: title, x: 'center', textStyle: {fontSize: 20}},
            tooltip: {trigger: 'axis'},
            legend: legend,
            xAxis: [{type: 'category', boundaryGap: false, data: data_xAxis}],
            yAxis: yAxis,
            dataZoom: {
                show: true,
                realtime: true,
                start: 0,
                end: 100
            },
            series: series
        };

        Chart.setOption(option);
        return Chart;
    }
</script>