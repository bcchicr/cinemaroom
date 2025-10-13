<?php

declare(strict_types=1);

namespace App\Domain\Exceptions;

final class FailedInvariantException extends DomainException
{
    public function getConventionalCode(): ErrorCode
    {
        return ErrorCode::FailedInvariant;
    }
}
