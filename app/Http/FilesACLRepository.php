<?php

namespace App\Http;

use Alexusmai\LaravelFileManager\Services\ACLService\ACLRepository;

class FilesACLRepository implements ACLRepository
{
    /**
     * Get user ID
     *
     * @return mixed
     */
    public function getUserID()
    {
        return auth('sanctum')->id();
    }

    /**
     * Get ACL rules list for user
     *
     * @return array
     */
    public function getRules(): array
    {
        if (auth('sanctum')->check()) {
            return [
                ['disk' => 'www', 'path' => '*', 'access' => 2],
            ];
        } else {
            return [
                ['disk' => 'www', 'path' => '*', 'access' => 0],
            ];
        }
    }
}
