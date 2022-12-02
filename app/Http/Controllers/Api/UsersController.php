<?php
/**
 * 耗子Linux面板 - 用户控制器
 * @author 耗子
 */

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use App\Models\Setting;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;
use Illuminate\Validation\ValidationException;

class UsersController extends Controller
{
    /**
     * 登录
     *
     * @param  Request  $request
     * @return JsonResponse
     */
    public function login(Request $request): JsonResponse
    {
        // 消毒数据
        try {
            $credentials = $this->validate($request, [
                'username' => 'required|max:255',
                'password' => 'required|max:255',
                'remember' => 'required|boolean'
            ]);
        } catch (ValidationException $e) {
            return response()->json([
                'code' => 1,
                'msg' => '参数错误：'.$e->getMessage(),
                'errors' => $e->errors()
            ], 422);
        }
        if (auth()->attempt(['username' => $credentials['username'], 'password' => $credentials['password']],
            $credentials['remember'])) {
            $user = auth()->user();
            // 多设备登录
            $multiLogin = Setting::query()->where('name', 'multi_login')->value('value');
            if ($multiLogin != 1) {
                $user->tokens()->delete();
            }
            $token = $user->createToken('token')->plainTextToken;
            return response()->json(['code' => 0, 'msg' => '登录成功', 'data' => ['access_token' => $token]]);
        } else {
            return response()->json(['code' => 1, 'msg' => '登录失败，用户名或密码错误']);
        }
    }

    public function getInfo(Request $request): JsonResponse
    {
        $user = $request->user();
        $res['code'] = 0;
        $res['msg'] = 'success';
        $res['data']['username'] = $user->username;
        return response()->json($res);
    }
}
