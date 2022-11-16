<?php

namespace Alexusmai\LaravelFileManager\Events;

use Illuminate\Http\Request;

class Deleted
{
    /**
     * @var string
     */
    private $disk;

    /**
     * @var array
     */
    private $items;

    /**
     * Deleted constructor.
     *
     * @param Request $request
     */
    public function __construct($disk, $items)
    {
        $this->disk = $disk;
        $this->items = $items;
    }

    /**
     * @return string
     */
    public function disk()
    {
        return $this->disk;
    }

    /**
     * @return array
     */
    public function items()
    {
        return $this->items;
    }
}
