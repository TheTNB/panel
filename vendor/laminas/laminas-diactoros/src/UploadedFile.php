<?php

declare(strict_types=1);

namespace Laminas\Diactoros;

use Psr\Http\Message\StreamInterface;
use Psr\Http\Message\UploadedFileInterface;

use function dirname;
use function fclose;
use function file_exists;
use function fopen;
use function fwrite;
use function is_dir;
use function is_resource;
use function is_string;
use function is_writable;
use function move_uploaded_file;
use function str_starts_with;
use function unlink;

use const PHP_SAPI;
use const UPLOAD_ERR_CANT_WRITE;
use const UPLOAD_ERR_EXTENSION;
use const UPLOAD_ERR_FORM_SIZE;
use const UPLOAD_ERR_INI_SIZE;
use const UPLOAD_ERR_NO_FILE;
use const UPLOAD_ERR_NO_TMP_DIR;
use const UPLOAD_ERR_OK;
use const UPLOAD_ERR_PARTIAL;

class UploadedFile implements UploadedFileInterface
{
    public const ERROR_MESSAGES = [
        UPLOAD_ERR_OK         => 'There is no error, the file uploaded with success',
        UPLOAD_ERR_INI_SIZE   => 'The uploaded file exceeds the upload_max_filesize directive in php.ini',
        UPLOAD_ERR_FORM_SIZE  => 'The uploaded file exceeds the MAX_FILE_SIZE directive that was '
            . 'specified in the HTML form',
        UPLOAD_ERR_PARTIAL    => 'The uploaded file was only partially uploaded',
        UPLOAD_ERR_NO_FILE    => 'No file was uploaded',
        UPLOAD_ERR_NO_TMP_DIR => 'Missing a temporary folder',
        UPLOAD_ERR_CANT_WRITE => 'Failed to write file to disk',
        UPLOAD_ERR_EXTENSION  => 'A PHP extension stopped the file upload.',
    ];

    private int $error;

    private ?string $file = null;

    private bool $moved = false;

    /** @var null|StreamInterface */
    private $stream;

    /**
     * @param string|resource|StreamInterface $streamOrFile
     * @throws Exception\InvalidArgumentException
     */
    public function __construct(
        $streamOrFile,
        private int $size,
        int $errorStatus,
        private ?string $clientFilename = null,
        private ?string $clientMediaType = null
    ) {
        if ($errorStatus === UPLOAD_ERR_OK) {
            if (is_string($streamOrFile)) {
                $this->file = $streamOrFile;
            }
            if (is_resource($streamOrFile)) {
                $this->stream = new Stream($streamOrFile);
            }

            if (! $this->file && ! $this->stream) {
                if (! $streamOrFile instanceof StreamInterface) {
                    throw new Exception\InvalidArgumentException('Invalid stream or file provided for UploadedFile');
                }
                $this->stream = $streamOrFile;
            }
        }

        if (0 > $errorStatus || 8 < $errorStatus) {
            throw new Exception\InvalidArgumentException(
                'Invalid error status for UploadedFile; must be an UPLOAD_ERR_* constant'
            );
        }
        $this->error = $errorStatus;
    }

    /**
     * {@inheritdoc}
     *
     * @throws Exception\UploadedFileAlreadyMovedException If the upload was not successful.
     */
    public function getStream(): StreamInterface
    {
        if ($this->error !== UPLOAD_ERR_OK) {
            throw Exception\UploadedFileErrorException::dueToStreamUploadError(
                self::ERROR_MESSAGES[$this->error]
            );
        }

        if ($this->moved) {
            throw new Exception\UploadedFileAlreadyMovedException();
        }

        if ($this->stream instanceof StreamInterface) {
            return $this->stream;
        }

        $this->stream = new Stream($this->file);
        return $this->stream;
    }

    /**
     * {@inheritdoc}
     *
     * @see http://php.net/is_uploaded_file
     * @see http://php.net/move_uploaded_file
     *
     * @param string $targetPath Path to which to move the uploaded file.
     * @throws Exception\UploadedFileErrorException If the upload was not successful.
     * @throws Exception\InvalidArgumentException If the $path specified is invalid.
     * @throws Exception\UploadedFileErrorException On any error during the
     *     move operation, or on the second or subsequent call to the method.
     */
    public function moveTo($targetPath): void
    {
        if ($this->moved) {
            throw new Exception\UploadedFileAlreadyMovedException('Cannot move file; already moved!');
        }

        if ($this->error !== UPLOAD_ERR_OK) {
            throw Exception\UploadedFileErrorException::dueToStreamUploadError(
                self::ERROR_MESSAGES[$this->error]
            );
        }

        if (! is_string($targetPath) || empty($targetPath)) {
            throw new Exception\InvalidArgumentException(
                'Invalid path provided for move operation; must be a non-empty string'
            );
        }

        $targetDirectory = dirname($targetPath);
        if (! is_dir($targetDirectory) || ! is_writable($targetDirectory)) {
            throw Exception\UploadedFileErrorException::dueToUnwritableTarget($targetDirectory);
        }

        $sapi = PHP_SAPI;
        switch (true) {
            case empty($sapi) || str_starts_with($sapi, 'cli') || str_starts_with($sapi, 'phpdbg') || ! $this->file:
                // Non-SAPI environment, or no filename present
                $this->writeFile($targetPath);

                if ($this->stream instanceof StreamInterface) {
                    $this->stream->close();
                }
                if (is_string($this->file) && file_exists($this->file)) {
                    unlink($this->file);
                }
                break;
            default:
                // SAPI environment, with file present
                if (false === move_uploaded_file($this->file, $targetPath)) {
                    throw Exception\UploadedFileErrorException::forUnmovableFile();
                }
                break;
        }

        $this->moved = true;
    }

    /**
     * {@inheritdoc}
     *
     * @return int|null The file size in bytes or null if unknown.
     */
    public function getSize(): ?int
    {
        return $this->size;
    }

    /**
     * {@inheritdoc}
     *
     * @see http://php.net/manual/en/features.file-upload.errors.php
     *
     * @return int One of PHP's UPLOAD_ERR_XXX constants.
     */
    public function getError(): int
    {
        return $this->error;
    }

    /**
     * {@inheritdoc}
     *
     * @return string|null The filename sent by the client or null if none
     *     was provided.
     */
    public function getClientFilename(): ?string
    {
        return $this->clientFilename;
    }

    /**
     * {@inheritdoc}
     */
    public function getClientMediaType(): ?string
    {
        return $this->clientMediaType;
    }

    /**
     * Write internal stream to given path
     */
    private function writeFile(string $path): void
    {
        $handle = fopen($path, 'wb+');
        if (false === $handle) {
            throw Exception\UploadedFileErrorException::dueToUnwritablePath();
        }

        $stream = $this->getStream();
        $stream->rewind();
        while (! $stream->eof()) {
            fwrite($handle, $stream->read(4096));
        }

        fclose($handle);
    }
}
