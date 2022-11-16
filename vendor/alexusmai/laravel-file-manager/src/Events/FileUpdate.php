<?php

namespace Alexusmai\LaravelFileManager\Events;

use Illuminate\Http\Request;

class FileUpdate
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
     * @var \Illuminate\Http\UploadedFile
     */
    private $file;

    /**
     * FileUpdate constructor.
     *
     * @param Request $request
     */
    public function __construct(Request $request)
    {
        $this->disk = $request->input('disk');
        $this->path = $request->input('path');
        $this->file = $request->file('file');
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
        if ($this->path) {
            return $this->path.'/'.$this->file->getClientOriginalName();
        }

        return $this->file->getClientOriginalName();
    }
}
