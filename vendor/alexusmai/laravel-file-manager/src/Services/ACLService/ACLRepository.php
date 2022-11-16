<?php

namespace Alexusmai\LaravelFileManager\Services\ACLService;

/**
 * Interface ACLRepository
 *
 * @package Alexusmai\LaravelFileManager\Services\ACLService
 */
interface ACLRepository
{
    /**
     * Get user ID
     *
     * @return mixed
     */
    public function getUserID();

    /**
     * Get ACL rules list for user
     *
     * You need to return an array, like this:
     *
     *  0 => [
     *      "disk" => "public"
     *      "path" => "music"
     *      "access" => 0
     *  ],
     *  1 => [
     *      "disk" => "public"
     *      "path" => "images"
     *      "access" => 1
     *  ]
     *
     * OR [] - if no results for selected user
     *
     * @return array
     */
    public function getRules(): array;
}
