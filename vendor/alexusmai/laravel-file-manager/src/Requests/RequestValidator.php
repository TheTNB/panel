<?php

namespace Alexusmai\LaravelFileManager\Requests;

use Alexusmai\LaravelFileManager\Services\ConfigService\ConfigRepository;
use Illuminate\Foundation\Http\FormRequest;
use Illuminate\Support\Facades\Storage;

class RequestValidator extends FormRequest
{
    use CustomErrorMessage;

    /**
     * Determine if the user is authorized to make this request.
     *
     * @return bool
     */
    public function authorize(): bool
    {
        return true;
    }

    /**
     * Get the validation rules that apply to the request.
     *
     * @return array
     */
    public function rules(): array
    {
        $config = resolve(ConfigRepository::class);

        return [
            'disk' => [
                'sometimes',
                'string',
                function ($attribute, $value, $fail) use($config) {
                    if (!in_array($value, $config->getDiskList()) ||
                        !array_key_exists($value, config('filesystems.disks'))
                    ) {
                        return $fail('diskNotFound');
                    }
                },
            ],
            'path' => [
                'sometimes',
                'string',
                'nullable',
                function ($attribute, $value, $fail) {
                    if ($value && !Storage::disk($this->input('disk'))->exists($value)
                    ) {
                        return $fail('pathNotFound');
                    }
                },
            ],
        ];
    }

    /**
     * Not found message
     *
     * @return string
     */
    public function message()
    {
        return 'notFound';
    }
}
