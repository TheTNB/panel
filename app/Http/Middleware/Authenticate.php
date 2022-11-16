<?php

namespace App\Http\Middleware;

use Closure;
use Illuminate\Auth\Middleware\Authenticate as Middleware;

use Illuminate\Http\JsonResponse;
use Illuminate\Http\Response;

class Authenticate extends Middleware
{
    /**
     * Handle an incoming request.
     */
    public function handle($request, Closure $next, ...$guards)
    {
       /* // 获取请求头中的token
        $token = $request->header('access_token') ?? $request->input('access_token');
        // 将token放入请求中
        $request->headers->set('Authorization', 'Bearer ' . $token);
        // 验证token
        $this->authenticate($request, $guards);*/
        // 验证通过

        return $next($request);
    }
    /**
     * Get the path the user should be redirected to when they are not authenticated.
     *
     * @param  \Illuminate\Http\Request  $request
     *
     * @return string|null
     */
    protected function redirectTo($request)
    {
        abort(response()->json([
            'code' => 1001,
            'message' => '登录状态失效'
        ]));
    }
}
