<?php

declare(strict_types=1);

namespace App\Application\Mappers;

use App\Application\Resources\AvatarResource;
use App\Domain\ValueObjects\Avatar;

final readonly class AvatarMapper
{
    public static function toResource(Avatar $avatar): AvatarResource
    {
        return new AvatarResource(
            avatarPath: $avatar->path(),
        );
    }
}
