<?php
/**
 * 耗子Linux面板 - 用户控制器
 * @author 耗子
 */
namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use Illuminate\Http\Request;
use Illuminate\Validation\ValidationException;

class UsersController extends Controller
{
    /**
     * 登录
     *
     * @param  Request  $request
     * @return \Illuminate\Http\JsonResponse
     */
    public function login(Request $request)
    {
        /*$user = User::create([
            'id' => '',
            'username' => 'haozi',
            'password' => Hash::make('haozi'),
        ]);
        return response()->json([
            'code' => 200,
            'message' => '注册成功',
            'data' => $user,
        ]);*/
        // 消毒数据
        try {
            $credentials = $this->validate($request, [
                'username' => 'required|max:255',
                'password' => 'required|max:255',
                'remember' => 'required|boolean'
            ]);
        } catch (ValidationException $e) {
            return response()->json([
                'message' => '参数错误',
                'errors' => $e->errors()
            ], 422);
        }
        if (auth()->attempt(['username' => $credentials['username'], 'password' => $credentials['password']], $credentials['remember'])) {
            $user = auth()->user();
            $user->tokens()->delete();
            $token = $user->createToken('token')->plainTextToken;
            return response()->json(['code' => 0, 'msg' => '登录成功', 'data' => ['access_token' => $token]]);
        } else {
            return response()->json(['code' => 1, 'msg' => '登录失败，用户名或密码错误']);
        }
    }
    public function getInfo(Request $request)
    {
        $user = $request->user();
        $res['code'] = 0;
        $res['msg'] = 'success';
        $res['data']['username'] = 'haozi';
        return response()->json($res);
    }
}
