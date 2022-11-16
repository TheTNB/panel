# Events

### BeforeInitialization

> Alexusmai\LaravelFileManager\Events\BeforeInitialization

Example:

```php
\Event::listen('Alexusmai\LaravelFileManager\Events\BeforeInitialization',
    function ($event) {
        
    }
);
```

### DiskSelected

> Alexusmai\LaravelFileManager\Events\DiskSelected

Example:

```php
\Event::listen('Alexusmai\LaravelFileManager\Events\DiskSelected',
    function ($event) {
        \Log::info('DiskSelected:', [$event->disk()]);
    }
);
```

### FilesUploading

> Alexusmai\LaravelFileManager\Events\FilesUploading

```php
\Event::listen('Alexusmai\LaravelFileManager\Events\FilesUploading',
    function ($event) {
        \Log::info('FilesUploading:', [
            $event->disk(),
            $event->path(),
            $event->files(),
            $event->overwrite(),
        ]);
    }
);
```

### FilesUploaded

> Alexusmai\LaravelFileManager\Events\FilesUploaded

```php
\Event::listen('Alexusmai\LaravelFileManager\Events\FilesUploaded',
    function ($event) {
        \Log::info('FilesUploaded:', [
            $event->disk(),
            $event->path(),
            $event->files(),
            $event->overwrite(),
        ]);
    }
);
```

### Deleting

> Alexusmai\LaravelFileManager\Events\Deleting

```php
\Event::listen('Alexusmai\LaravelFileManager\Events\Deleting',
    function ($event) {
        \Log::info('Deleting:', [
            $event->disk(),
            $event->items(),
        ]);
    }
);
```

### Deleted

> Alexusmai\LaravelFileManager\Events\Deleted

```php
\Event::listen('Alexusmai\LaravelFileManager\Events\Deleted',
    function ($event) {
        \Log::info('Deleted:', [
            $event->disk(),
            $event->items(),
        ]);
    }
);
```

### Paste

> Alexusmai\LaravelFileManager\Events\Paste

```php
\Event::listen('Alexusmai\LaravelFileManager\Events\Paste',
    function ($event) {
        \Log::info('Paste:', [
            $event->disk(),
            $event->path(),
            $event->clipboard(),
        ]);
    }
);
```

### Rename

> Alexusmai\LaravelFileManager\Events\Rename

```php
\Event::listen('Alexusmai\LaravelFileManager\Events\Rename',
    function ($event) {
        \Log::info('Rename:', [
            $event->disk(),
            $event->newName(),
            $event->oldName(),
            $event->type(), // 'file' or 'dir'
        ]);
    }
);
```

### Download

> Alexusmai\LaravelFileManager\Events\Download

```php
\Event::listen('Alexusmai\LaravelFileManager\Events\Download',
    function ($event) {
        \Log::info('Download:', [
            $event->disk(),
            $event->path(),
        ]);
    }
);
```

*When using a text editor, the file you are editing is also downloaded! And this event is triggered!*

### DirectoryCreating

> Alexusmai\LaravelFileManager\Events\DirectoryCreating

```php
\Event::listen('Alexusmai\LaravelFileManager\Events\DirectoryCreating',
    function ($event) {
        \Log::info('DirectoryCreating:', [
            $event->disk(),
            $event->path(),
            $event->name(),
        ]);
    }
);
```

### DirectoryCreated

> Alexusmai\LaravelFileManager\Events\DirectoryCreated

```php
\Event::listen('Alexusmai\LaravelFileManager\Events\DirectoryCreated',
    function ($event) {
        \Log::info('DirectoryCreated:', [
            $event->disk(),
            $event->path(),
            $event->name(),
        ]);
    }
);
```

### FileCreating

> Alexusmai\LaravelFileManager\Events\FileCreating

```php
\Event::listen('Alexusmai\LaravelFileManager\Events\FileCreating',
    function ($event) {
        \Log::info('FileCreating:', [
            $event->disk(),
            $event->path(),
            $event->name(),
        ]);
    }
);
```

### FileCreated

> Alexusmai\LaravelFileManager\Events\FileCreated

```php
\Event::listen('Alexusmai\LaravelFileManager\Events\FileCreated',
    function ($event) {
        \Log::info('FileCreated:', [
            $event->disk(),
            $event->path(),
            $event->name(),
        ]);
    }
);
```

### FileUpdate

> Alexusmai\LaravelFileManager\Events\FileUpdate

```php
\Event::listen('Alexusmai\LaravelFileManager\Events\FileUpdate',
    function ($event) {
        \Log::info('FileUpdate:', [
            $event->disk(),
            $event->path(),
        ]);
    }
);
```

### Zip

> Alexusmai\LaravelFileManager\Events\Zip

```php
\Event::listen('Alexusmai\LaravelFileManager\Events\Zip',
    function ($event) {
        \Log::info('Zip:', [
            $event->disk(),
            $event->path(),
            $event->name(),
            $event->elements(),
        ]);
    }
);
```

### ZipCreated

> Alexusmai\LaravelFileManager\Events\ZipCreated

```php
\Event::listen('Alexusmai\LaravelFileManager\Events\ZipCreated',
    function ($event) {
        \Log::info('ZipCreated:', [
            $event->disk(),
            $event->path(),
            $event->name(),
            $event->elements(),
        ]);
    }
);
```

### ZipFailed

> Alexusmai\LaravelFileManager\Events\ZipCreated

```php
\Event::listen('Alexusmai\LaravelFileManager\Events\ZipFailed',
    function ($event) {
        \Log::info('ZipFailed:', [
            $event->disk(),
            $event->path(),
            $event->name(),
            $event->elements(),
        ]);
    }
);
```

### Unzip

> Alexusmai\LaravelFileManager\Events\Unzip

```php
\Event::listen('Alexusmai\LaravelFileManager\Events\Unzip',
    function ($event) {
        \Log::info('Unzip:', [
            $event->disk(),
            $event->path(),
            $event->folder(),
        ]);
    }
);
```

### UnzipCreated

> Alexusmai\LaravelFileManager\Events\UnzipCreated

```php
\Event::listen('Alexusmai\LaravelFileManager\Events\UnzipCreated',
    function ($event) {
        \Log::info('UnzipCreated:', [
            $event->disk(),
            $event->path(),
            $event->folder(),
        ]);
    }
);
```

### UnzipFailed

> Alexusmai\LaravelFileManager\Events\UnzipFailed

```php
\Event::listen('Alexusmai\LaravelFileManager\Events\UnzipFailed',
    function ($event) {
        \Log::info('UnzipFailed:', [
            $event->disk(),
            $event->path(),
            $event->folder(),
        ]);
    }
);
```
