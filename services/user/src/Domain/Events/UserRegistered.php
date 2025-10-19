<?php

declare(strict_types=1);

namespace App\Domain\Events;

use App\Domain\ValueObjects\UserId;

final readonly class UserRegistered extends DomainEvent
{
    public function __construct(
        private UserId $id,
        private string $username,
    ) {
        parent::__construct();
    }

    #[\Override]
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
}
