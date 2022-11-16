<?php
/**
 * 耗子Linux面板 - 数据库控制器
 * @author 耗子
 */
namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use Illuminate\Http\Request;

class DatabasesController extends Controller
{
    public function get_mysql_list()
    {
        $password    = Db::table('setting')->where('name', 'mysql_root_password')->value('value');
        $db_raw      = trim(shell_exec("/www/server/mysql/bin/mysql -u root -p".$password." -e \"SHOW DATABASES;\""));
        $res['code'] = 0;
        $res['msg']  = 'success';
        $db          = explode("\n", $db_raw);
        foreach ($db as $key => $value) {
            if ($value == 'Database' || $value == 'information_schema' || $value == 'mysql' || $value == 'performance_schema' || $value == 'sys') {
                unset($db[$key]);
            } else {
                $res['data'][] = $value;
            }
        }

        return response()->json($res);
    }
}
