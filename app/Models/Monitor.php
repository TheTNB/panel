<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;

class Monitor extends Model
{
    use HasFactory;

    // 白名单
    protected $fillable = [
        'info',
        'created_at',
        'updated_at',
    ];
}
