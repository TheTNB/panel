<!doctype html>
<html lang="en">
<head>
    <!-- Required meta tags -->
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <!-- Bootstrap CSS -->
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css">

    <title>File manager and TinyMCE</title>
</head>
<body>

<div class="container">
    <form method="post">
        <textarea id="my-textarea"></textarea>
    </form>
</div>

<!-- TinyMCE -->
<script src='https://cloud.tinymce.com/stable/tinymce.min.js'></script>
<script>
  document.addEventListener('DOMContentLoaded', function() {
    tinymce.init({
      selector: '#my-textarea',
      plugins: [
        'advlist autolink lists link image charmap print preview hr anchor pagebreak',
        'searchreplace wordcount visualblocks visualchars code fullscreen',
        'insertdatetime media nonbreaking save table contextmenu directionality',
        'emoticons template paste textcolor colorpicker textpattern',
      ],
      toolbar: 'insertfile undo redo | styleselect | bold italic | alignleft aligncenter alignright alignjustify | bullist numlist outdent indent | link image media',
      relative_urls: false,
      file_browser_callback: function(field_name, url, type, win) {
        tinyMCE.activeEditor.windowManager.open({
          file: '/file-manager/tinymce',
          title: 'Laravel File Manager',
          width: window.innerWidth * 0.8,
          height: window.innerHeight * 0.8,
          resizable: 'yes',
          close_previous: 'no',
        }, {
          setUrl: function(url) {
            win.document.getElementById(field_name).value = url;
          },
        });
      },
    });
  });
</script>
</body>
</html>
