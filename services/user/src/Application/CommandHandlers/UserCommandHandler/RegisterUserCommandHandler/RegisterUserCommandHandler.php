<?php

declare(strict_types=1);

namespace App\Application\CommandHandlers\UserCommandHandler\RegisterUserCommandHandler;

use App\Application\Resources\UserResource;
use App\Domain\Aggregates\User\User;
use App\Domain\Events\EventDispatcher;
use App\Domain\Repositories\UserRepository;

final readonly class RegisterUserCommandHandler
{
    public function __construct(
        private UserRepository $userRepository,
        private EventDispatcher $eventDispatcher,
    ) {
    }

    public function handle(RegisterUserCommand $command): UserResource
    {
        $userId = $this->userRepository->nextIdentity();
        $user = User::register(
            id: $userId,
            username: $command->username,
        );

        $this->userRepository->save($user);

        $events = $user->pullPendingEvents();
        foreach ($events as $event) {
            $this->eventDispatcher->dispatch($event);
        }

        return new UserResource($user);
    }
}
