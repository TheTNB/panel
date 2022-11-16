<!DOCTYPE html>
<html>
<head>
    <meta name="csrf-token" content="{{ csrf_token() }}">
</head>
<body>
<link rel="stylesheet" href="https://cdnjs.cdn.wepublish.cn/bootstrap/5.1.3/css/bootstrap.min.css">
<link rel="stylesheet" href="https://cdnjs.cdn.wepublish.cn/bootstrap-icons/1.8.1/font/bootstrap-icons.min.css">
<div id="fm" style="height: 800px;"></div>
<link rel="stylesheet" href="{{ asset('../../vendor/file-manager/css/file-manager.css') }}">
<script src="{{ asset('../../vendor/file-manager/js/file-manager.js') }}"></script>
</body>
</html>
