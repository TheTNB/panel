<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;

class Task extends Model
{
    use HasFactory;

    // 白名单
    protected $fillable = [
        'job_id',
        'name',
        'shell',
        'log',
        'status',
        'created_at',
        'updated_at',
    ];
}
