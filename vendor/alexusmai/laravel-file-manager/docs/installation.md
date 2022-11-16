# Installation

1. Install package - using composer

```
composer require alexusmai/laravel-file-manager
```

For Laravel 5 - 8 use v2.5.4

```
composer require alexusmai/laravel-file-manager "2.5.4"
```

2. If you use Laravel 5.4, then add service provider to config/app.php (for the Laravel 5.5 and higher skip this step):

```php
Alexusmai\LaravelFileManager\FileManagerServiceProvider::class,
```

3. Publish configuration file

```bash
php artisan vendor:publish --tag=fm-config
```

4. You can install npm package directly and use it in your vue application - more information about it -
   [vue-laravel-file-manager](https://github.com/alexusmai/vue-laravel-file-manager)

> OR

Publish compiled and minimized js and css files

```
php artisan vendor:publish --tag=fm-assets
```

Open the view file where you want to place the application block, and add:

* add a csrf token to head block if you did not do it before

```
<!-- CSRF Token -->
<meta name="csrf-token" content="{{ csrf_token() }}">
```

* For version 3 and higher - the frontend package uses **Bootstrap 5** and **Bootstrap Icons** styles, if you already use it,
  then you do not need to connect any styles. Otherwise, add -

```
<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css">
<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap-icons@1.8.1/font/bootstrap-icons.css">
```

* For old versions - the frontend package uses **Bootstrap 4** and **Font Awesome 5** styles, if you already use it,
  then you do not need to connect any styles. Otherwise, add -

```
<link rel="stylesheet" href="https://use.fontawesome.com/releases/v5.7.0/css/all.css">
<link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.1.3/css/bootstrap.min.css">
```

* add file manager styles

```
<link rel="stylesheet" href="{{ asset('vendor/file-manager/css/file-manager.css') }}">
```

* add file manager js

```
<script src="{{ asset('vendor/file-manager/js/file-manager.js') }}"></script>
```

* For version 3 and higher - add div for application (set application height!)

```
<div id="fm" style="height: 600px;"></div>
```    

* For old versions - add div for application (set application height!)

```
<div style="height: 600px;">
    <div id="fm"></div>
</div>
```

## What's next

[Configuration](./configuration.md)
