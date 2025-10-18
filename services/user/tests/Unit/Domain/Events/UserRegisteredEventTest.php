<?php

declare(strict_types=1);

namespace App\Tests\Unit\Domain\Events;

use App\Domain\Events\UserRegistered;
use App\Domain\ValueObjects\UserId;
use PHPUnit\Framework\TestCase;
use Ramsey\Uuid\Uuid;

final class UserRegisteredEventTest extends TestCase
{
    public function testInstantiates(): void
    {
        $userId = new UserId(Uuid::uuid7());
        $username = 'test_user';

        $event = new UserRegistered($userId, $username);

        $this->assertTrue($event->id()->equals($userId));
        $this->assertEquals($username, $event->username());
    }
}
