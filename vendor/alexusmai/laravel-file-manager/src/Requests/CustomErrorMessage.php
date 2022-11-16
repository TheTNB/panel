<?php

namespace Alexusmai\LaravelFileManager\Requests;

use Illuminate\Contracts\Validation\Validator;
use Illuminate\Http\Exceptions\HttpResponseException;

trait CustomErrorMessage
{
    /**
     * Validation error response
     *
     * @param  Validator  $validator
     */
    protected function failedValidation(Validator $validator)
    {
        $message = (method_exists($this, 'message'))
            ? $this->container->call([$this, 'message'])
            : 'The given data was invalid.';

        throw new HttpResponseException(response()->json([
            'errors'  => $validator->errors(),
            'message' => $message,
        ], 422));
    }
}
