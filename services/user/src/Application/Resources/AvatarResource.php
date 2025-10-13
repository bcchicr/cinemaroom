<?php

declare(strict_types=1);

namespace App\Application\Resources;

final readonly class AvatarResource
{
    public function __construct(
        public string $avatarPath,
    ) {
    }
}
