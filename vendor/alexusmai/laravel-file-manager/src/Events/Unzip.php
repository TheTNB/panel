<?php

namespace Alexusmai\LaravelFileManager\Events;

use Illuminate\Http\Request;

class Unzip
{
    /**
     * @var string
     */
    private $disk;

    /**
     * @var string
     */
    private $path;

    /**
     * @var string
     */
    private $folder;

    /**
     * Unzip constructor.
     *
     * @param  Request  $request
     */
    public function __construct(Request $request)
    {
        $this->disk = $request->input('disk');
        $this->path = $request->input('path');
        $this->folder = $request->input('folder');
    }

    /**
     * @return string
     */
    public function disk()
    {
        return $this->disk;
    }

    /**
     * @return string
     */
    public function path()
    {
        return $this->path;
    }

    /**
     * @return string
     */
    public function folder()
    {
        return $this->folder;
    }
}
