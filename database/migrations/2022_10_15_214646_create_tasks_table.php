<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::create('tasks', function (Blueprint $table) {
            $table->id();
            $table->integer('job_id')->comment('任务ID');
            $table->string('name')->comment('任务名');
            $table->string('status')->default('waiting')->comment('任务状态');
            $table->string('shell')->nullable()->comment('任务脚本');
            $table->string('log')->nullable()->comment('任务日志目录');
            $table->timestamps();
        });
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('tasks');
    }
};
