<?php

namespace App\Http\Middleware;

use Closure;

class AddAccessToken
{
    /**
     * Handle an incoming request.
     */
    public function handle($request, Closure $next, ...$guards)
    {
        // 获取请求头中的token
        $token = $request->header('access_token') ?? $request->input('access_token');
        // 将token放入请求中
        $request->headers->set('Authorization', 'Bearer '.$token);

        return $next($request);
    }
}
