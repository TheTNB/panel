<!doctype html>
<html lang="en">
<head>
    <!-- Required meta tags -->
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">

    <!-- Bootstrap CSS -->
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css">

    <title>File manager and CKeditor</title>
</head>
<body>

<div class="container">
    <textarea name="editor"></textarea>
</div>

<script src="https://cdn.ckeditor.com/4.11.2/standard/ckeditor.js"></script>
<script>
  CKEDITOR.replace( 'editor', {filebrowserImageBrowseUrl: '/file-manager/ckeditor'});
</script>
</body>
</html>
