<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;

class Cron extends Model
{
    use HasFactory;

    // 白名单
    protected $fillable = [
        'name',
        'status',
        'type',
        'time',
        'shell',
        'created_at',
        'updated_at',
    ];
}
