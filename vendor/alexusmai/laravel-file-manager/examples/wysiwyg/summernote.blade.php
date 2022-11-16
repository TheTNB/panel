<!doctype html>
<html lang="en">
<head>
    <!-- Required meta tags -->
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <!-- Bootstrap CSS -->
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css">
    <!-- SummerNote -->
    <link href="https://cdnjs.cloudflare.com/ajax/libs/summernote/0.8.11/summernote-bs4.css" rel="stylesheet">

    <title>File Manager and SummerNote</title>
</head>
<body>

<div class="container">
    <div id="summernote"></div>
</div>

<!-- jQuery first, then Popper.js, then Bootstrap JS -->
<script src="https://code.jquery.com/jquery-3.3.1.slim.min.js"></script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.14.3/umd/popper.min.js"></script>
<script src="https://stackpath.bootstrapcdn.com/bootstrap/4.1.3/js/bootstrap.min.js"></script>
<!-- SummerNote js -->
<script src="https://cdnjs.cloudflare.com/ajax/libs/summernote/0.8.11/summernote-bs4.js"></script>

<script>
  $(document).ready(function(){
    // File manager button (image icon)
    const FMButton = function(context) {
      const ui = $.summernote.ui;
      const button = ui.button({
        contents: '<i class="note-icon-picture"></i> ',
        tooltip: 'File Manager',
        click: function() {
          window.open('/file-manager/summernote', 'fm', 'width=1400,height=800');
        }
      });
      return button.render();
    };

    $('#summernote').summernote({
      toolbar: [
        // [groupName, [list of button]]
        ['style', ['bold', 'italic', 'underline', 'clear']],
        ['font', ['strikethrough', 'superscript', 'subscript']],
        ['fontsize', ['fontsize']],
        ['color', ['color']],
        ['para', ['ul', 'ol', 'paragraph']],
        ['height', ['height']],
        ['fm-button', ['fm']],
      ],
      buttons: {
        fm: FMButton
      }
    });
  });

  // set file link
  function fmSetLink(url) {
    $('#summernote').summernote('insertImage', url);
  }
</script>
</body>
</html>
