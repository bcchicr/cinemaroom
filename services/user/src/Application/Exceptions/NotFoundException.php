<?php

declare(strict_types=1);

namespace App\Application\Exceptions;

use App\Infrastructure\Attributes\GrpcStatusProvider\GrpcStatus;

use const Grpc\STATUS_NOT_FOUND;

#[GrpcStatus(code: STATUS_NOT_FOUND)]
final class NotFoundException extends ApplicationException
{
    #[\Override]
    public function getConventionalCode(): ErrorCode
    {
        return ErrorCode::NotFound;
    }
}
