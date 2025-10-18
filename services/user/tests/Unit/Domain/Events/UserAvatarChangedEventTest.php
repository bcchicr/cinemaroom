<?php

declare(strict_types=1);

namespace App\Tests\Unit\Domain\Events;

use App\Domain\Events\UserAvatarChanged;
use App\Domain\ValueObjects\Avatar;
use App\Domain\ValueObjects\UserId;
use PHPUnit\Framework\TestCase;
use Ramsey\Uuid\Uuid;

final class UserAvatarChangedEventTest extends TestCase
{
    public function testInstantiates(): void
    {
        $userId = new UserId(Uuid::uuid7());
        $oldAvatar = Avatar::null();
        $newAvatar = Avatar::fromString('/path/to/avatar.jpg');

        $event = new UserAvatarChanged($userId, $oldAvatar, $newAvatar);

        $this->assertTrue($userId->equals($event->id()));
        $this->assertTrue($oldAvatar->equals($event->oldAvatar()));
        $this->assertTrue($newAvatar->equals($event->newAvatar()));
    }
}
