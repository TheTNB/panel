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
        Schema::create('websites', function (Blueprint $table) {
            $table->id();
            $table->string('name')->index()->comment('网站名称');
            $table->boolean('status')->comment('是否运行');
            $table->string('path')->comment('网站路径');
            $table->integer('php')->comment('PHP版本');
            $table->boolean('ssl')->comment('是否启用SSL');
            $table->dateTime('ssl_date')->comment('SSL到期时间')->nullable();
            $table->string('note')->comment('备注')->nullable();
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
        Schema::dropIfExists('websites');
    }
};
