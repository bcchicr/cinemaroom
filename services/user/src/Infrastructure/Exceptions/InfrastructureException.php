<?php

declare(strict_types=1);

namespace App\Infrastructure\Exceptions;

abstract class InfrastructureException extends \RuntimeException
{
    abstract public function getConventionalCode(): ErrorCode;
}
