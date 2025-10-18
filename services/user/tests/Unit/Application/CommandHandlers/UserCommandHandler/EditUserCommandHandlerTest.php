<?php

declare(strict_types=1);

namespace App\Tests\Unit\Application\CommandHandlers\UserCommandHandler;

use App\Application\CommandHandlers\UserCommandHandler\EditUserCommandHandler\EditUserCommand;
use App\Application\CommandHandlers\UserCommandHandler\EditUserCommandHandler\EditUserCommandHandler;
use App\Application\Exceptions\NotFoundException;
use App\Domain\Aggregates\User\User;
use App\Domain\Events\EventDispatcher;
use App\Domain\Repositories\UserRepository;
use App\Domain\ValueObjects\Avatar;
use App\Domain\ValueObjects\UserId;
use PHPUnit\Framework\TestCase;
use Ramsey\Uuid\Uuid;

final class EditUserCommandHandlerTest extends TestCase
{
    public function testHandle(): void
    {
        $userRepository = $this->createMock(UserRepository::class);
        $eventDispatcher = $this->createMock(EventDispatcher::class);

        $id = Uuid::uuid7();

        $newUsername = 'new_username';
        $newBio = 'new_bio';
        $newAvatarPath = '/new/avatar.jpg';
        $command = new EditUserCommand(
            $id,
            $newUsername,
            $newBio,
            $newAvatarPath,
        );

        $oldUsername = 'old_username';
        $user = User::register(new UserId($id), $oldUsername);
        $user->pullPendingEvents();

        $userRepository->expects($this->once())
            ->method('findById')
            ->with($this->callback(fn (UserId $userId) => $userId->value()->equals($id)))
            ->willReturn($user);

        $userRepository->expects($this->once())
            ->method('save')
            ->with($this->callback(function (User $user) use ($id, $newUsername, $newBio, $newAvatarPath) {
                return $user->id()->value()->equals($id)
                    && $user->username() === $newUsername
                    && $user->bio() === $newBio
                    && $user->avatar()->path() === $newAvatarPath;
            }));

        $eventDispatcher->expects($this->exactly(3))
            ->method('dispatch');

        $handler = new EditUserCommandHandler($userRepository, $eventDispatcher);
        $result = $handler->handle($command);

        $this->assertTrue($id->equals($result->id()));
        $this->assertEquals($newUsername, $result->username());
        $this->assertEquals($newBio, $result->bio());
        $this->assertEquals($newAvatarPath, $result->avatar()->path());
    }

    public function testHandleNullValues(): void
    {
        $userRepository = $this->createMock(UserRepository::class);
        $eventDispatcher = $this->createMock(EventDispatcher::class);

        $id = Uuid::uuid7();

        $command = new EditUserCommand(
            $id,
            null,
            null,
            null,
        );

        $oldUsername = 'old_username';
        $oldBio = 'old_bio';
        $oldAvatarPath = '/old/avatar.jpg';
        $user = User::hydrate(
            id: new UserId($id),
            username: $oldUsername,
            bio: $oldBio,
            avatar: Avatar::fromString($oldAvatarPath),
        );
        $user->pullPendingEvents();

        $userRepository->expects($this->once())
            ->method('findById')
            ->with($this->callback(fn (UserId $userId) => $userId->value()->equals($id)))
            ->willReturn($user);

        $userRepository->expects($this->once())
            ->method('save')
            ->with($this->callback(function (User $user) use ($id, $oldUsername, $oldBio, $oldAvatarPath) {
                return $user->id()->value()->equals($id)
                    && $user->username() === $oldUsername
                    && $user->bio() === $oldBio
                    && $user->avatar()->path() === $oldAvatarPath;
            }));

        $eventDispatcher->expects($this->exactly(0))
            ->method('dispatch');

        $handler = new EditUserCommandHandler($userRepository, $eventDispatcher);
        $result = $handler->handle($command);

        $this->assertTrue($id->equals($result->id()));
        $this->assertEquals($oldUsername, $result->username());
        $this->assertEquals($oldBio, $result->bio());
        $this->assertEquals($oldAvatarPath, $result->avatar()->path());
    }

    public function testHandleNotFound(): void
    {
        $userRepository = $this->createMock(UserRepository::class);
        $eventDispatcher = $this->createMock(EventDispatcher::class);

        $id = Uuid::uuid4();
        $command = new EditUserCommand(
            $id,
            'new_username',
            'new_bio',
            '/new/avatar.jpg',
        );

        $userRepository->expects($this->once())
            ->method('findById')
            ->with($this->callback(fn (UserId $userId) => $userId->value()->equals($id)))
            ->willReturn(null);

        $this->expectException(NotFoundException::class);

        $handler = new EditUserCommandHandler($userRepository, $eventDispatcher);
        $handler->handle($command);
    }
}
