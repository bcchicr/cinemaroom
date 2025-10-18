<?php

declare(strict_types=1);

namespace App\Tests\Unit\Application\CommandHandlers\Resources;

use App\Application\Resources\UserResource;
use App\Domain\Aggregates\User\User;
use App\Domain\ValueObjects\Avatar;
use App\Domain\ValueObjects\UserId;
use PHPUnit\Framework\TestCase;
use Ramsey\Uuid\Uuid;

final class UserResourceTest extends TestCase
{
    public function testInstantiatesWithNonNullAvatar(): void
    {
        $user = User::hydrate(
            new UserId(Uuid::uuid7()),
            'username',
            bio: 'bio',
            avatar: Avatar::fromString('/path/to/avatar.jpg'),
        );

        $resource = new UserResource($user);

        $this->assertTrue($user->id()->value()->equals($resource->id()));
        $this->assertEquals($user->username(), $resource->username());
        $this->assertEquals($user->avatar()->path(), $resource->avatar()->path());
        $this->assertEquals($user->bio(), $resource->bio());
    }

    public function testInstantiatesWithNullAvatar(): void
    {
        $user = User::hydrate(
            new UserId(Uuid::uuid7()),
            'username',
            bio: 'bio',
            avatar: Avatar::null(),
        );

        $resource = new UserResource($user);

        $this->assertTrue($user->id()->value()->equals($resource->id()));
        $this->assertEquals($user->username(), $resource->username());
        $this->assertNull($resource->avatar());
        $this->assertEquals($user->bio(), $resource->bio());
    }
}
