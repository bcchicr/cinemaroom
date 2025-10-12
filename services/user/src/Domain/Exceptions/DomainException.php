<?php

declare(strict_types=1);

namespace App\Domain\Exceptions;

abstract class DomainException extends \RuntimeException
{
    abstract public function getConventionalCode(): ErrorCode;
}
