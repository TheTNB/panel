<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;

class Setting extends Model
{
    use HasFactory;

    // 白名单
    protected $fillable = ['name', 'value', 'created_at', 'updated_at'];
}
