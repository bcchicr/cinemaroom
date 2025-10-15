<?php

declare(strict_types=1);

namespace App\Domain\Exceptions;

use App\Infrastructure\Attributes\GrpcStatusProvider\GrpcStatus;

use const Grpc\STATUS_FAILED_PRECONDITION;

#[GrpcStatus(code: STATUS_FAILED_PRECONDITION)]
final class FailedInvariantException extends DomainException
{
    public function getConventionalCode(): ErrorCode
    {
        return ErrorCode::FailedInvariant;
    }
}
