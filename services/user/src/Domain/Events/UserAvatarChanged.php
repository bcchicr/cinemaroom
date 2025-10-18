<?php

declare(strict_types=1);

namespace App\Domain\Events;

use App\Domain\ValueObjects\Avatar;
use App\Domain\ValueObjects\UserId;

final readonly class UserAvatarChanged extends DomainEvent
{
    public function __construct(
        private UserId $id,
        private Avatar $oldAvatar,
        private Avatar $newAvatar,
    ) {
        parent::__construct();
    }

    public function eventName(): string
    {
        return 'user.avatar_changed';
    }

    public function id(): UserId
    {
        return $this->id;
    }

    public function oldAvatar(): Avatar
    {
        return $this->oldAvatar;
    }

    public function newAvatar(): Avatar
    {
        return $this->newAvatar;
    }
}
