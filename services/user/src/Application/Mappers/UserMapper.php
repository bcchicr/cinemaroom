<?php

declare(strict_types=1);

namespace App\Application\Mappers;

use App\Application\Resources\UserResource;
use App\Domain\Aggregates\User\User;

final readonly class UserMapper
{
    public static function toResource(User $user): UserResource
    {

        return new UserResource(
            id: $user->id()->value(),
            username: $user->username(),
            bio: $user->bio(),
            avatarPath: $user->avatar()?->path(),
        );
    }
}
