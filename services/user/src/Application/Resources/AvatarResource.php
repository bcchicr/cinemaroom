<?php

declare(strict_types=1);

namespace App\Application\Resources;

use App\Application\Exceptions\InvalidArgumentException;
use App\Domain\ValueObjects\Avatar;

class AvatarResource
{
    public function __construct(
        private readonly Avatar $avatar,
    ) {
        if ($avatar->isNull()) {
            throw new InvalidArgumentException('Avatar is null');
        }
    }

    public function path(): string
    {
        return $this->avatar->path();
    }
}
