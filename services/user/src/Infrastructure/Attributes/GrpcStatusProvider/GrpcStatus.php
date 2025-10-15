<?php

declare(strict_types=1);

namespace App\Infrastructure\Attributes\GrpcStatusProvider;

#[\Attribute(\Attribute::TARGET_CLASS)]
final readonly class GrpcStatus
{
    public function __construct(
        public int $code,
    ) {
    }
}
