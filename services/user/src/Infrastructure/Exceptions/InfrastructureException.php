<?php

declare(strict_types=1);

namespace App\Infrastructure\Exceptions;

use App\Infrastructure\Attributes\GrpcStatusProvider\GrpcStatus;

use const Grpc\STATUS_INTERNAL;

abstract class InfrastructureException extends \RuntimeException
{
    abstract public function getConventionalCode(): ErrorCode;
}
