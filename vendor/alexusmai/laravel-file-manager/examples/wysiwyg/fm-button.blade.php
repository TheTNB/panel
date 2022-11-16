<!doctype html>
<html lang="en">
<head>
    <!-- Required meta tags -->
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <!-- Bootstrap CSS -->
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css">

    <title>File Manager and standalone button</title>
</head>
<body>

<div class="container">
    <div class="form-row">
        <div class="form-group col-md-6">
            <label for="image_label">Image</label>
            <div class="input-group">
                <input type="text" id="image_label" class="form-control" name="image"
                       aria-label="Image" aria-describedby="button-image">
                <div class="input-group-append">
                    <button class="btn btn-outline-secondary" type="button" id="button-image">Select</button>
                </div>
            </div>
        </div>
    </div>
</div>

<script>
  document.addEventListener("DOMContentLoaded", function() {

    document.getElementById('button-image').addEventListener('click', (event) => {
      event.preventDefault();

      window.open('/file-manager/fm-button', 'fm', 'width=1400,height=800');
    });
  });

  // set file link
  function fmSetLink($url) {
    document.getElementById('image_label').value = $url;
  }
</script>
</body>
</html>
