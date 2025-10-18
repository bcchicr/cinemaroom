<?php

declare(strict_types=1);

namespace App\Domain\Exceptions;

final class InvalidArgumentException extends DomainException
{
    public function getConventionalCode(): ErrorCode
    {
        return ErrorCode::InvalidArgument;
    }
}
