<?php

declare(strict_types=1);

namespace App\Application\Resources;

use Ramsey\Uuid\UuidInterface;

final readonly class UserResource
{
    public function __construct(
        public UuidInterface $id,
        public string $username,
        public ?string $bio,
        public ?string $avatarPath,
    ) {
    }
}
