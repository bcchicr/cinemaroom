<?php

declare(strict_types=1);

namespace App\Domain\Factories;

use App\Domain\ValueObjects\Avatar;

final readonly class AvatarFactory
{
    public function create(
        string $path,
    ) {
        return new Avatar($path);
    }
}
