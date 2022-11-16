<!DOCTYPE html>
<html lang="{{ app()->getLocale() }}">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <!-- CSRF Token -->
    <meta name="csrf-token" content="{{ csrf_token() }}">

    <title>{{ config('app.name', 'File Manager') }}</title>

    <!-- Styles -->
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css">
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap-icons@1.8.1/font/bootstrap-icons.css">
    <link rel="stylesheet" href="{{ asset('vendor/file-manager/css/file-manager.css') }}">
</head>
<body>
<div class="container-fluid">
    <div class="row">
        <div class="col-md-12" id="fm-main-block">
            <div id="fm"></div>
        </div>
    </div>
</div>

<!-- File manager -->
<script src="{{ asset('vendor/file-manager/js/file-manager.js') }}"></script>
<script>
  document.addEventListener('DOMContentLoaded', function() {
    // set fm height
    document.getElementById('fm-main-block').setAttribute('style', 'height:' + window.innerHeight + 'px');

    // Add callback to file manager
    fm.$store.commit('fm/setFileCallBack', function(fileUrl) {
      window.opener.fmSetLink(fileUrl);
      window.close();
    });
  });
</script>
</body>
</html>

