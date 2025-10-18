<?php

declare(strict_types=1);

namespace App\Endpoints\Grpc\Mappers;

use App\Application\Resources\UserResource;
use GRPC\Services\Common\v1\Uuid;
use GRPC\Services\Users\v1\Avatar;
use GRPC\Services\Users\v1\User;

final readonly class UserMapper
{
    public function toGrpc(UserResource $user): User
    {
        $grpcUser = new User();
        $grpcUser->setId(
            (new Uuid())->setValue($user->id()->toString()),
        );
        $grpcUser->setUsername($user->username());
        $grpcUser->setBio($user->bio());

        if (null !== $user->avatar()) {
            $grpcUser->setAvatar(
                (new Avatar())->setPath($user->avatar()->path()),
            );
        }

        return $grpcUser;
    }
}
