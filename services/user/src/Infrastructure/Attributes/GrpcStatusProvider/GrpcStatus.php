<?php

declare(strict_types=1);

namespace App\Infrastructure\Attributes\GrpcStatusProvider;

use Google\Rpc\Code;

#[\Attribute(\Attribute::TARGET_CLASS)]
final readonly class GrpcStatus
{
    public function __construct(
        public int $code,
    ) {}
}
