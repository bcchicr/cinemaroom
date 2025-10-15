<?php

declare(strict_types=1);

namespace App\Application\Exceptions;

abstract class ApplicationException extends \RuntimeException
{
    abstract public function getConventionalCode(): ErrorCode;
}
