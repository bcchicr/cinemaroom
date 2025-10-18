<?php

declare(strict_types=1);

namespace App\Tests\Unit\Endpoints\Grpc\Mappers;

use App\Application\Resources\AvatarResource;
use App\Application\Resources\UserResource;
use App\Endpoints\Grpc\Mappers\UserMapper;
use PHPUnit\Framework\TestCase;
use Ramsey\Uuid\Uuid;

final class UserMapperTest extends TestCase
{
    public function testMapUserResourceWithAvatarToGrpcUser(): void
    {
        $avatar = $this->createMock(AvatarResource::class);
        $avatar->method('path')->willReturn('/path/to/avatar.jpg');

        $user = $this->createMock(UserResource::class);
        $user->method('id')->willReturn(Uuid::uuid7());
        $user->method('username')->willReturn('test_user');
        $user->method('bio')->willReturn('test_bio');
        $user->method('avatar')->willReturn($avatar);

        $mapper = new UserMapper();

        $userGrpc = $mapper->toGrpc($user);

        $this->assertEquals($userGrpc->getId()->getValue(), $user->id()->toString());
        $this->assertEquals($userGrpc->getUsername(), $user->username());
        $this->assertEquals($userGrpc->getBio(), $user->bio());
        $this->assertTrue($userGrpc->hasAvatar());
        $this->assertEquals($userGrpc->getAvatar()->getPath(), $user->avatar()->path());
    }

    public function testMapUserResourceWithoutAvatarToGrpcUser(): void
    {
        $user = $this->createMock(UserResource::class);
        $user->method('id')->willReturn(Uuid::uuid7());
        $user->method('username')->willReturn('test_user');
        $user->method('bio')->willReturn('test_bio');
        $user->method('avatar')->willReturn(null);

        $mapper = new UserMapper();

        $userGrpc = $mapper->toGrpc($user);

        $this->assertEquals($userGrpc->getId()->getValue(), $user->id()->toString());
        $this->assertEquals($userGrpc->getUsername(), $user->username());
        $this->assertEquals($userGrpc->getBio(), $user->bio());
        $this->assertFalse($userGrpc->hasAvatar());
    }
}
