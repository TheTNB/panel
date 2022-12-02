<?php
/**
 * 耗子Linux面板 - 安全控制器
 * @author 耗子
 */

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;

class SafesController extends Controller
{
    /**
     * 获取防火墙状态
     * @return JsonResponse
     */
    public function getFirewallStatus(): JsonResponse
    {
        $firewallStatus = trim(shell_exec("systemctl status firewalld | grep Active | awk '{print $3}'"));
        $res['code'] = 0;
        $res['msg'] = 'success';
        if ($firewallStatus == '(running)') {
            $res['data'] = 1;
        } else {
            $res['data'] = 0;
        }

        return response()->json($res);
    }

    /**
     * 设置防火墙状态
     * @param  Request  $request
     * @return JsonResponse
     */
    public function setFirewallStatus(Request $request): JsonResponse
    {
        $status = $request->input('status');
        if ($status) {
            shell_exec("systemctl enable firewalld");
            shell_exec("systemctl start firewalld");
        } else {
            shell_exec("systemctl stop firewalld");
            shell_exec("systemctl disable firewalld");
        }
        $res['code'] = 0;
        $res['msg'] = 'success';
        return response()->json($res);
    }

    /**
     * 获取SSH状态
     * @return JsonResponse
     */
    public function getSshStatus(): JsonResponse
    {
        $sshStatus = trim(shell_exec("systemctl status sshd | grep Active | awk '{print $3}'"));
        $res['code'] = 0;
        $res['msg'] = 'success';
        if ($sshStatus == '(running)') {
            $res['data'] = 1;
        } else {
            $res['data'] = 0;
        }

        return response()->json($res);
    }

    /**
     * 设置SSH状态
     * @param  Request  $request
     * @return JsonResponse
     */
    public function setSshStatus(Request $request): JsonResponse
    {
        $status = $request->input('status');
        if ($status) {
            shell_exec("systemctl enable sshd");
            shell_exec("systemctl start sshd");
        } else {
            shell_exec("systemctl stop sshd");
            shell_exec("systemctl disable sshd");
        }
        $res['code'] = 0;
        $res['msg'] = 'success';
        return response()->json($res);
    }

    /**
     * 获取SSH端口
     * @return JsonResponse
     */
    public function getSshPort(): JsonResponse
    {
        $sshPort = trim(shell_exec("cat /etc/ssh/sshd_config | grep 'Port ' | awk '{print $2}'"));
        $res['code'] = 0;
        $res['msg'] = 'success';
        $res['data'] = $sshPort;
        return response()->json($res);
    }

    /**
     * 设置SSH端口
     * @param  Request  $request
     * @return JsonResponse
     */
    public function setSshPort(Request $request): JsonResponse
    {
        $port = $request->input('port');
        $oldPort = trim(shell_exec("cat /etc/ssh/sshd_config | grep 'Port ' | awk '{print $2}'"));
        shell_exec("sed -i 's/#Port ".$oldPort."/Port ".$port."/g' /etc/ssh/sshd_config");
        shell_exec("sed -i 's/Port ".$oldPort."/Port ".$port."/g' /etc/ssh/sshd_config");
        // 判断ssh是否开启
        $sshStatus = trim(shell_exec("systemctl status sshd | grep Active | awk '{print $3}'"));
        if ($sshStatus == '(running)') {
            shell_exec("systemctl restart sshd");
        }
        $res['code'] = 0;
        $res['msg'] = 'success';
        return response()->json($res);
    }

    /**
     * 获取ping状态
     * @return JsonResponse
     */
    public function getPingStatus(): JsonResponse
    {
        $pingStatus = trim(shell_exec("firewall-cmd --query-rich-rule='rule protocol value=icmp drop' 2>&1"));
        $res['code'] = 0;
        $res['msg'] = 'success';
        if ($pingStatus == 'yes') {
            $res['data'] = 0;
        } else {
            $res['data'] = 1;
        }

        return response()->json($res);
    }

    /**
     * 设置ping状态
     * @param  Request  $request
     * @return JsonResponse
     */
    public function setPingStatus(Request $request): JsonResponse
    {
        $status = $request->input('status');
        if ($status) {
            shell_exec("firewall-cmd --permanent --remove-rich-rule='rule protocol value=icmp drop' 2>&1");
        } else {
            shell_exec("firewall-cmd --permanent --add-rich-rule='rule protocol value=icmp drop' 2>&1");
        }
        shell_exec("firewall-cmd --reload");
        $res['code'] = 0;
        $res['msg'] = 'success';
        return response()->json($res);
    }

    /**
     * 获取防火墙规则
     * @return JsonResponse
     */
    public function getFirewallRules(): JsonResponse
    {
        $firewallRules = trim(shell_exec("firewall-cmd --list-all 2>&1"));
        // 判断是否开启
        if (str_contains($firewallRules, 'not running')) {
            $res['code'] = 0;
            $res['msg'] = 'success';
            $res['data'] = [];
            return response()->json($res);
        }
        // 正则匹配出ports
        preg_match('/ports: (.*)/', $firewallRules, $matches);
        $rawPorts = $matches[1];
        // 22/tcp 80/tcp 443/tcp 8888/tcp 5432/tcp
        $ports = explode(' ', $rawPorts);
        // 对ports进行分割为port=>protocol形式
        $rules = [];
        foreach ($ports as $port) {
            $rule = explode('/', $port);
            $rules[] = [
                'port' => $rule[0],
                'protocol' => $rule[1],
            ];
        }

        $res['code'] = 0;
        $res['msg'] = 'success';
        $res['data'] = $rules;
        return response()->json($res);
    }

    /**
     * 添加防火墙规则
     * @param  Request  $request
     * @return JsonResponse
     */
    public function addFirewallRule(Request $request): JsonResponse
    {
        $port = $request->input('port');
        $protocol = $request->input('protocol');
        // 判断是否开启
        $firewallStatus = trim(shell_exec("firewall-cmd --state 2>&1"));
        if ($firewallStatus != 'running') {
            $res['code'] = 1;
            $res['msg'] = '防火墙未开启';
            return response()->json($res);
        }
        // 清空当前规则
        shell_exec("firewall-cmd --remove-port=".$port."/".$protocol." --permanent");
        // 添加新的防火墙规则
        shell_exec("firewall-cmd --add-port=".$port."/".$protocol." --permanent");
        // 重启防火墙
        shell_exec("firewall-cmd --reload");
        $res['code'] = 0;
        $res['msg'] = 'success';
        return response()->json($res);
    }

    /**
     * 删除防火墙规则
     * @param  Request  $request
     * @return JsonResponse
     */
    public function deleteFirewallRule(Request $request): JsonResponse
    {
        $port = $request->input('port');
        $protocol = $request->input('protocol');
        // 判断是否开启
        $firewallStatus = trim(shell_exec("firewall-cmd --state 2>&1"));
        if ($firewallStatus != 'running') {
            $res['code'] = 1;
            $res['msg'] = '防火墙未开启';
            return response()->json($res);
        }
        // 清空当前规则
        shell_exec("firewall-cmd --remove-port=".$port."/".$protocol." --permanent");
        // 重启防火墙
        shell_exec("firewall-cmd --reload");
        $res['code'] = 0;
        $res['msg'] = 'success';
        return response()->json($res);
    }
}
