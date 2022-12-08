<?php
/**
 * 耗子Linux面板 - 帮助函数
 * @author 耗子
 */

/**
 * 获取网络统计信息
 * @return array
 */
function getNetInfo(): array
{
    $networkRaw = file_get_contents('/proc/net/dev');
    $networkArr = explode("\n", $networkRaw);
    foreach ($networkArr as $key => $val) {
        if ($key < 2) {
            continue;
        }
        $val = str_replace(':', ' ', trim($val));
        $val = preg_replace("/[ ]+/", " ", $val);
        $arr = explode(' ', $val);
        if (!empty($arr[0])) {
            $arr = array($arr[0], $arr[1], $arr[9]);
            $allRs[$arr[0].$key] = $arr;
        }
    }
    ksort($allRs);
    $tx = 0;
    $rx = 0;
    foreach ($allRs as $key => $val) {
        // 排除本地lo
        if (str_contains($key, 'lo')) {
            continue;
        }
        $tx += $val[2];
        $rx += $val[1];
    }
    $res['tx'] = $tx;
    $res['rx'] = $rx;
    return $res;
}

/**
 * 格式化bytes
 * @param $size
 * @return string
 */
function formatBytes($size): string
{
    $size = is_numeric($size) ? $size : 0;
    $units = array(' B', ' KB', ' MB', ' GB', ' TB');
    for ($i = 0; $size >= 1024 && $i < 4; $i++) {
        $size /= 1024;
    }
    return round($size, 2).$units[$i];
}

/**
 * 裁剪字符串
 * @param $begin
 * @param $end
 * @param $str
 * @return string
 */
function cut($begin, $end, $str): string
{
    $b = mb_strpos($str, $begin) + mb_strlen($begin);
    $e = mb_strpos($str, $end) - $b;
    return mb_substr($str, $b, $e);
}
