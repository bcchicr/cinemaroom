<?php

declare(strict_types=1);

namespace App\Domain\Exceptions;

use App\Infrastructure\Attributes\GrpcStatusProvider\GrpcStatus;
use const Grpc\STATUS_INVALID_ARGUMENT;

#[GrpcStatus(code: STATUS_INVALID_ARGUMENT)]
final class InvalidArgumentException extends DomainException
{
    public function getConventionalCode(): ErrorCode
    {
        return ErrorCode::InvalidArgument;
    }
}
