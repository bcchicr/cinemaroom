<?php

declare(strict_types=1);

namespace App\Tests\Unit\Application\CommandHandlers\UserCommandHandler;

use App\Application\CommandHandlers\UserCommandHandler\RegisterUserCommandHandler\RegisterUserCommand;
use App\Application\CommandHandlers\UserCommandHandler\RegisterUserCommandHandler\RegisterUserCommandHandler;
use App\Domain\Aggregates\User\User;
use App\Domain\Events\EventDispatcher;
use App\Domain\Repositories\UserRepository;
use App\Domain\ValueObjects\UserId;
use PHPUnit\Framework\TestCase;
use Ramsey\Uuid\Uuid;

final class RegisterUserCommandHandlerTest extends TestCase
{
    public function testInvoke(): void
    {
        $userRepository = $this->createMock(UserRepository::class);
        $eventDispatcher = $this->createMock(EventDispatcher::class);

        $username = 'test_user';
        $id = new UserId(Uuid::uuid7());
        $command = new RegisterUserCommand($id->value(), $username);

        $userRepository->expects($this->once())
            ->method('save')
            ->with($this->callback(function (User $user) use ($id, $username) {
                return $id->equals($user->id())
                    && $username === $user->username();
            }));

        $eventDispatcher->expects($this->once())
            ->method('dispatch');

        $handler = new RegisterUserCommandHandler($userRepository, $eventDispatcher);
        $result = $handler->handle($command);

        $this->assertTrue($id->value()->equals($result->id()));
        $this->assertEquals($username, $result->username());
    }
}
