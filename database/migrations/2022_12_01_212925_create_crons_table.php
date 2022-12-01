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
        Schema::create('crons', function (Blueprint $table) {
            $table->id();
            $table->string('name')->nullable()->comment('任务名称');
            $table->boolean('status')->comment('任务状态');
            $table->string('type')->comment('任务类型');
            $table->string('time')->comment('任务周期');
            $table->text('shell')->comment('任务脚本文件');
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
        Schema::dropIfExists('crons');
    }
};
