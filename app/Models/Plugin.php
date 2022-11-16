<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;

class Plugin extends Model
{
    use HasFactory;

    // 白名单
    protected $fillable = ['slug', 'name', 'version', 'show', 'created_at', 'updated_at'];
}
