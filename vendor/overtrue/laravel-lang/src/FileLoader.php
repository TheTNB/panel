<?php

namespace Overtrue\LaravelLang;

use Illuminate\Support\Str;
use Illuminate\Filesystem\Filesystem;
use Illuminate\Translation\FileLoader as LaravelTranslationFileLoader;
use RuntimeException;

class FileLoader extends LaravelTranslationFileLoader
{
    /**
     * @var array
     */
    protected $paths;

    /**
     * @var string[]
     */
    protected $customJsonPaths = [];

    /**
     * Create a new file loader instance.
     *
     * @param \Illuminate\Filesystem\Filesystem $files
     * @param string $path
     * @param array $paths
     */
    public function __construct(Filesystem $files, string $path, array $paths = [])
    {
        $this->paths = $paths;

        parent::__construct($files, $path);
    }

    /**
     * Load the messages for the given locale.
     *
     * @param string $locale
     * @param string $group
     * @param string $namespace
     *
     * @return array
     */
    public function load($locale, $group, $namespace = null)
    {
        $defaults = [];

        $locale = str_replace('-', '_', $locale);

        foreach ($this->paths as $path) {
            $defaults = array_replace_recursive($defaults, $this->loadPath($path, $locale, $group));
        }

        return array_replace_recursive($defaults, parent::load($locale, $group, $namespace));
    }

    /**
     * Fall back to base locale (i.e. de) if a countries specific locale (i.e. de-CH) is not available.
     *
     * @param string $path
     * @param string $locale
     * @param string $group
     *
     * @return array
     */
    protected function loadPath($path, $locale, $group): array
    {
        $result = parent::loadPath($path, $locale, $group);

        if (empty($result) && Str::contains($locale, '-')) {
            return parent::loadPath($path, strstr($locale, '-', true), $group);
        }

        return $result;
    }

    /**
     * Add a new JSON path to the loader.
     *
     * @param  string  $path
     * @return void
     */
    public function addJsonPath($path)
    {
        $this->customJsonPaths[] = $path;
        parent::addJsonPath($path);
    }

    /**
     * Load a locale from the given JSON file path.
     *
     * @param  string  $locale
     * @return array
     *
     * @throws RuntimeException
     */
    protected function loadJsonPaths($locale)
    {
        return collect(array_merge($this->jsonPaths, [$this->path]))
            ->reduce(function ($output, $path) use ($locale) {
                if (in_array($path, $this->customJsonPaths)) {
                    $locale = "{$locale}/{$locale}";
                }

                if ($this->files->exists($full = "{$path}/{$locale}.json")) {
                    $decoded = json_decode($this->files->get($full), true);

                    if (is_null($decoded) || json_last_error() !== JSON_ERROR_NONE) {
                        throw new RuntimeException("Translation file [{$full}] contains an invalid JSON structure.");
                    }

                    $output = array_merge($output, $decoded);
                }

                return $output;
            }, []);
    }
}
