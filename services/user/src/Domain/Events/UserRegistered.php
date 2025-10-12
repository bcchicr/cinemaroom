<?php

declare(strict_types=1);

namespace App\Domain\Events;

use App\Domain\ValueObjects\Avatar;
use App\Domain\ValueObjects\UserId;

final readonly class UserRegistered extends DomainEvent
{
    public function __construct(
        private UserId $id,
        private string $username,
        private ?string $bio,
        private ?Avatar $avatar,
    ) {
        parent::__construct();
    }

    public function eventName(): string
    {
        return 'user.registered';
    }

    public function id(): UserId
    {
        return $this->id;
    }

    public function username(): string
    {
        return $this->username;
    }

    public function bio(): ?string
    {
        return $this->bio;
    }

    public function avatar(): ?Avatar
    {
        return $this->avatar;
    }
}
