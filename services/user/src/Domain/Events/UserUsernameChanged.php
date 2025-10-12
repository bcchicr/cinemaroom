<?php

declare(strict_types=1);

namespace App\Domain\Events;

use App\Domain\ValueObjects\UserId;

final readonly class UserUsernameChanged extends DomainEvent
{
    public function __construct(
        private UserId $id,
        private string $oldUsername,
        private string $newUsername,
    ) {
        parent::__construct();
    }

    public function eventName(): string
    {
        return 'user.username_changed';
    }

    public function id(): UserId
    {
        return $this->id;
    }

    public function oldUsername(): string
    {
        return $this->oldUsername;
    }

    public function newUsername(): string
    {
        return $this->newUsername;
    }
}
