<?php

declare(strict_types=1);

namespace App\Application\Resources;

use App\Domain\Aggregates\User\User;
use Ramsey\Uuid\UuidInterface;

class UserResource
{
    public function __construct(
        private readonly User $user,
    ) {}

    public function id(): UuidInterface
    {
        return $this->user->id()->value();
    }

    public function username(): string
    {
        return $this->user->username();
    }

    public function bio(): string
    {
        return $this->user->bio();
    }

    public function avatar(): ?AvatarResource
    {
        if ($this->user->avatar()->isNull()) {
            return null;
        }

        return new AvatarResource($this->user->avatar());
    }
}
