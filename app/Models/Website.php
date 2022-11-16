<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;

class Website extends Model
{
    use HasFactory;

    // 白名单
    protected $fillable = ['name', 'status', 'path', 'php', 'ssl', 'ssl_date', 'note', 'created_at', 'updated_at'];
}
