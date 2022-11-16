<?php

namespace Alexusmai\LaravelFileManager\Events;

use Illuminate\Http\Request;

class ZipCreated
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
    private $name;

    /**
     * @var array|string|null
     */
    private $elements;

    /**
     * ZipCreated constructor.
     *
     * @param  Request  $request
     */
    public function __construct(Request $request)
    {
        $this->disk = $request->input('disk');
        $this->path = $request->input('path');
        $this->name = $request->input('name');
        $this->elements = $request->input('elements');
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
    public function name()
    {
        return $this->name;
    }

    /**
     * @return array|string|null
     */
    public function elements()
    {
        return $this->elements;
    }
}
