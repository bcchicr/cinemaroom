<?php

declare(strict_types=1);

namespace App\Tests\Unit\Domain\Events;

use App\Domain\Events\UserUsernameChanged;
use App\Domain\ValueObjects\UserId;
use PHPUnit\Framework\TestCase;
use Ramsey\Uuid\Uuid;

final class UserUsernameChangedEventTest extends TestCase
{
    public function testInstantiationAndGetters(): void
    {
        $userId = new UserId(Uuid::uuid4());
        $oldUsername = 'old_username';
        $newUsername = 'new_username';

        $event = new UserUsernameChanged($userId, $oldUsername, $newUsername);

        $this->assertTrue($event->id()->equals($userId));
        $this->assertSame($oldUsername, $event->oldUsername());
        $this->assertSame($newUsername, $event->newUsername());
    }
}
