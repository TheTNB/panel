/**
 * console  
 */

layui.define(function(exports){
  
  /*
    下面通过 layui.use 分段加载不同的模块，实现不同区域的同时渲染，从而保证视图的快速呈现
  */
  
  
  //区块轮播切换
  layui.use(['admin', 'carousel'], function(){
    var $ = layui.$
    ,admin = layui.admin
    ,carousel = layui.carousel
    ,element = layui.element
    ,device = layui.device();

    //轮播切换
    $('.layadmin-carousel').each(function(){
      var othis = $(this);
      carousel.render({
        elem: this
        ,width: '100%'
        ,arrow: 'none'
        ,interval: othis.data('interval')
        ,autoplay: othis.data('autoplay') === true
        ,trigger: (device.ios || device.android) ? 'click' : 'hover'
        ,anim: othis.data('anim')
      });
    });
    
    element.render('progress');
    
  });

  //数据概览
  layui.use(['admin', 'carousel', 'echarts'], function(){
    var $ = layui.$
    ,admin = layui.admin
    ,carousel = layui.carousel
    ,echarts = layui.echarts;
    
    var echartsApp = [], options = [
      //今日流量趋势
      {
        title: {
          text: '今日流量趋势',
          x: 'center',
          textStyle: {
            fontSize: 14
          }
        },
        tooltip : {
          trigger: 'axis'
        },
        legend: {
          data:['','']
        },
        xAxis : [{
          type : 'category',
          boundaryGap : false,
          data: ['06:00','06:30','07:00','07:30','08:00','08:30','09:00','09:30','10:00','11:30','12:00','12:30','13:00','13:30','14:00','14:30','15:00','15:30','16:00','16:30','17:00','17:30','18:00','18:30','19:00','19:30','20:00','20:30','21:00','21:30','22:00','22:30','23:00','23:30']
        }],
        yAxis : [{
          type : 'value'
        }],
        series : [{
          name:'PV',
          type:'line',
          smooth:true,
          itemStyle: {normal: {areaStyle: {type: 'default'}}},
          data: [111,222,333,444,555,777,3333,33333,55555,88888,33333,3333,7777,11888,28888,38888,58888,42222,39999,28888,17777,9777,6555,5555,3333,2222,3111,6999,5888,2777,1777,999,888,777]
        },{
          name:'UV',
          type:'line',
          smooth:true,
          itemStyle: {normal: {areaStyle: {type: 'default'}}},
          data: [11,22,33,44,55,66,333,3333,5555,12312,3333,333,777,1188,2777,3888,7777,4222,3999,2888,1777,966,655,555,333,222,311,699,588,277,166,99,88,77]
        }]
      },
      
      /*
      //访客浏览器分布
      { 
        title : {
          text: '访客浏览器分布',
          x: 'center',
          textStyle: {
            fontSize: 14
          }
        },
        tooltip : {
          trigger: 'item',
          formatter: "{a} <br/>{b} : {c} ({d}%)"
        },
        legend: {
          orient : 'vertical',
          x : 'left',
          data:['Chrome','Firefox','IE 8.0','Safari','其它浏览器']
        },
        series : [{
          name:'访问来源',
          type:'pie',
          radius : '55%',
          center: ['50%', '50%'],
          data:[
            {value:9052, name:'Chrome'},
            {value:1610, name:'Firefox'},
            {value:3200, name:'IE 8.0'},
            {value:535, name:'Safari'},
            {value:1700, name:'其它浏览器'}
          ]
        }]
      },
      */
      
      //新增的用户量
      {
        title: {
          text: '最近一周新增的用户量',
          x: 'center',
          textStyle: {
            fontSize: 14
          }
        },
        tooltip : { //提示框
          trigger: 'axis',
          formatter: "{b}<br>新增用户：{c}"
        },
        xAxis : [{ //X轴
          type : 'category',
          data : ['11-07', '11-08', '11-09', '11-10', '11-11', '11-12', '11-13']
        }],
        yAxis : [{  //Y轴
          type : 'value'
        }],
        series : [{ //内容
          type: 'line',
          data:[200, 300, 400, 610, 150, 270, 380],
        }]
      }
    ]
    ,elemDataView = $('#LAY-index-dataview').children('div')
    ,renderDataView = function(index){
      echartsApp[index] = echarts.init(elemDataView[index], layui.echartsTheme);
      echartsApp[index].setOption(options[index]);
      //window.onresize = echartsApp[index].resize;
      admin.resize(function(){
        echartsApp[index].resize();
      });
    };
    
    
    //没找到DOM，终止执行
    if(!elemDataView[0]) return;
    
    
    
    renderDataView(0);
    
    //触发数据概览轮播
    var carouselIndex = 0;
    carousel.on('change(LAY-index-dataview)', function(obj){
      renderDataView(carouselIndex = obj.index);
    });
    
    //触发侧边伸缩
    layui.admin.on('side', function(){
      setTimeout(function(){
        renderDataView(carouselIndex);
      }, 300);
    });
    
    //触发路由
    layui.admin.on('hash(tab)', function(){
      layui.router().path.join('') || renderDataView(carouselIndex);
    });
  });
  
  //地图
  layui.use(['carousel', 'echarts'], function(){
    var $ = layui.$
    ,carousel = layui.carousel
    ,echarts = layui.echarts;
    
    var echartsApp = [], options = [
      {
        title : {
          text: '访客地区分布',
          subtext: '不完全统计'
        },
        tooltip : {
          trigger: 'item'
        },
        dataRange: {
          orient: 'horizontal',
          min: 0,
          max: 60000,
          text:['高','低'],
          splitNumber:0
        },
        series : [
          { 
            name: '访客地区分布',
            type: 'map',
            mapType: 'china',    
            selectedMode : 'multiple',
            itemStyle:{
                normal:{label:{show:true}},
                emphasis:{label:{show:true}}
            },
            data:[
                {name:'西藏', value:60},
                {name:'青海', value:167},
                {name:'宁夏', value:210},
                {name:'海南', value:252},
                {name:'甘肃', value:502},
                {name:'贵州', value:570},
                {name:'新疆', value:661},
                {name:'云南', value:8890},
                {name:'重庆', value:10010},
                {name:'吉林', value:5056},
                {name:'山西', value:2123},
                {name:'天津', value:9130},
                {name:'江西', value:10170},
                {name:'广西', value:6172},
                {name:'陕西', value:9251},
                {name:'黑龙江', value:5125},
                {name:'内蒙古', value:1435},
                {name:'安徽', value:9530},
                {name:'北京', value:51919},
                {name:'福建', value:3756},
                {name:'上海', value:59190},
                {name:'湖北', value:37109},
                {name:'湖南', value:8966},
                {name:'四川', value:31020},
                {name:'辽宁', value:7222},
                {name:'河北', value:3451},
                {name:'河南', value:9693},
                {name:'浙江', value:62310},
                {name:'山东', value:39231},
                {name:'江苏', value:35911},
                {name:'广东', value:55891}
            ]
          }
        ]
      }
    ]
    ,elemDataView = $('#LAY-index-pagethree-home').children('div')
    ,renderDataView = function(index){
      echartsApp[index] = echarts.init(elemDataView[index], layui.echartsTheme);
      echartsApp[index].setOption(options[index]);
      window.onresize = echartsApp[index].resize;
    }; 
    //没找到DOM，终止执行
    if(!elemDataView[0]) return;
 
    renderDataView(0);  
  });


  //table
  layui.use('table', function(){
    var $ = layui.$
    ,table = layui.table;
    
    //今日热搜
    table.render({
      elem: '#LAY-index-topSearch'
      ,url: './res/json/console/top-search.js' //模拟接口
      ,page: true
      ,cols: [[
        {type: 'numbers', fixed: 'left'}
        ,{field: 'keywords', title: '关键词', minWidth: 300, templet: '<div><a href="https://www.baidu.com/s?wd={{ d.keywords }}" target="_blank" class="layui-table-link">{{ d.keywords }}</div>'}
        ,{field: 'frequency', title: '搜索次数', minWidth: 120, sort: true}
        ,{field: 'userNums', title: '用户数', sort: true}
      ]]
      ,skin: 'line'
    });
    
    //今日热贴
    table.render({
      elem: '#LAY-index-topCard'
      ,url: './res/json/console/top-card.js' //模拟接口
      ,page: true
      ,cellMinWidth: 120
      ,cols: [[
        {type: 'numbers', fixed: 'left'}
        ,{field: 'title', title: '标题', minWidth: 300, templet: '<div><a href="{{ d.href }}" target="_blank" class="layui-table-link">{{ d.title }}</div>'}
        ,{field: 'username', title: '发帖者'}
        ,{field: 'channel', title: '类别'}
        ,{field: 'crt', title: '点击率', sort: true}
      ]]
      ,skin: 'line'
    });
    
    //项目进展
    table.render({
      elem: '#LAY-home-homepage-console'
      ,url: './res/json/console/prograss.js' //模拟接口
      ,cols: [[
        {type: 'checkbox', fixed: 'left'}
        ,{field: 'prograss', title: '任务'}
        ,{field: 'time', title: '所需时间'}
        ,{field: 'complete', title: '完成情况'
          ,templet: function(d){
            if(d.complete == '已完成'){
              return '<del style="color: #5FB878;">'+ d.complete +'</del>'
            }else if(d.complete == '进行中'){
              return '<span style="color: #FFB800;">'+ d.complete +'</span>'
            }else{
              return '<span style="color: #FF5722;">'+ d.complete +'</span>'
            }
          }
        }
      ]]
      ,skin: 'line'
    });
  });
  
  exports('console', {})
});