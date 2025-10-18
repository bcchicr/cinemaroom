<?php

declare(strict_types=1);

namespace App\Tests\Unit\Domain\Events;

use App\Domain\Events\UserBioChanged;
use App\Domain\ValueObjects\UserId;
use PHPUnit\Framework\TestCase;
use Ramsey\Uuid\Uuid;

final class UserBioChangedEventTest extends TestCase
{
    public function testInstantiationAndGetters(): void
    {
        $userId = new UserId(Uuid::uuid4());
        $oldBio = 'old bio';
        $newBio = 'new bio';

        $event = new UserBioChanged($userId, $oldBio, $newBio);

        $this->assertTrue($event->id()->equals($userId));
        $this->assertEquals($oldBio, $event->oldBio());
        $this->assertEquals($newBio, $event->newBio());
    }
}
