<?php

declare(strict_types=1);

namespace App\Domain\Events;

use App\Domain\ValueObjects\UserId;

final readonly class UserBioChanged extends DomainEvent
{
    public function __construct(
        private UserId $id,
        private ?string $oldBio,
        private string $newBio,
    ) {
        parent::__construct();
    }

    public function eventName(): string
    {
        return 'user.bio_changed';
    }

    public function id(): UserId
    {
        return $this->id;
    }

    public function oldBio(): string
    {
        return $this->oldBio;
    }

    public function newBio(): string
    {
        return $this->newBio;
    }
}
