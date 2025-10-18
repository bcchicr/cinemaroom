<?php

declare(strict_types=1);

namespace App\Application\Exceptions;

final class InvalidArgumentException extends ApplicationException
{
    public function getConventionalCode(): ErrorCode
    {
        return ErrorCode::InvalidArgument;
    }
}
