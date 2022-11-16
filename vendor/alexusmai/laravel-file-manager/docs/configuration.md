# Configuration

Open configuration file - config/file-manager.php

- fill the disk list from config/filesystem.php (select the desired drive names)
- set cache
- select file manager windows configuration

**!!! Be sure to add your middleware to restrict access to the application !!!**

**Don't forget to configure your php and Nginx**

```
// PHP
upload_max_filesize,
post_max_size

// Nginx
client_max_body_size
```

### You can set default disk and default path

You have two variants for how to do it:

1. Add this params to the config file (config/file-manager.php)

```php
/**
 * Default disk for left manager
 * null - auto select the first disk in the disk list
 */
'leftDisk'  => 'public',

/**
 * Default disk for right manager
 * null - auto select the first disk in the disk list
 */
'rightDisk' => null,

/**
 * Default path for left manager
 * null - root directory
 */
'leftPath'  => 'directory/sub-directory',

/**
 * Default path for right manager
 * null - root directory
 */
'rightPath' => null,
```

2 Or you can add this params in URL

```
http://site.name/?leftDisk=public

http://site.name/?leftDisk=public&rightDisk=images

http://site.name/?leftDisk=public&leftPath=directory/sub-directory

http://site.name/?leftDisk=public&leftPath=directory2&rightDisk=images&rightPath=cars/vw/golf
// %2F - /, %20 - space
http://site.name/?leftDisk=public&leftPath=directory2&rightDisk=images&rightPath=cars%2Fvw%2Fgolf
```

leftDisk and leftPath is default for the file manager windows configuration - 1,2

**You can't add a disk that does not exist in the diskList array !**

**! Params in URL have more weight than params in config file. It means that URL params can overwrite your config params. !**

## Disk settings example

- config/filesystems.php

```php
// Filesystem Disks
'disks' => [
    // images folder in public path
    'images' => [
        'driver' => 'local',
        'root' => public_path('images'),
        'url' => env('APP_URL').'/images',
    ],

    // public folder in storage/app/public
    'public' => [
        'driver' => 'local',
        'root' => storage_path('app/public'),
        'url' => env('APP_URL').'/storage', // https://laravel.com/docs/5.7/filesystem#file-urls
        'visibility' => 'public',
    ],

    // ftp
    'dd-wrt' => [
        'driver'   => 'ftp',
        'host'     => 'ftp.dd-wrt.com',
        'username' => 'anonymous',
        'passive'  => true,
        'timeout'  => 30,
    ],
],
```

- config/file-manager.php

```php
// You need to enter the disks you want to use in the file manager
'diskList'  => ['images', 'public'],
```

## Normalization of uploaded file names

If you expect to work with files that may have filenames that are not considered *normal* `f.e.: "DCIM_2021  - čšč& (1).jpg.jpg"` so basically any time you give the option to upload files to users, you can set `slugifyNames` to `true` and have the names ran through `Str::slug()` before saving it so you file will look something like `dcim-2021-csc-1.jpg` 

## Dynamic configuration

You can create your own configuration, for example for different users or their roles.

Create new class - example - TestConfigRepository

```php
namespace App\Http;

use Alexusmai\LaravelFileManager\Services\ConfigService\ConfigRepository;

class TestConfigRepository implements ConfigRepository
{
    // implement all methods from interface
    
    /**
     * Get disk list
     *
     * ['public', 'local', 's3']
     *
     * @return array
     */
    public function getDiskList(): array
    {
        if (\Auth::id() === 1) {
            return [
                ['public', 'local', 's3'],
            ];
        }
        
        return ['public'];
    }
    
    ...
}
```

For example see [src/Services/ConfigService/DefaultConfigRepository.php](https://github.com/alexusmai/laravel-file-manager/blob/master/src/Services/ConfigService/DefaultConfigRepository.php)

## What's next

[ACL](./acl.md)

[Integration](./integration.md)
