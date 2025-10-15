<?php

declare(strict_types=1);

namespace App\Application\CommandHandlers\UserCommandHandler\EditUserCommandHandler;

use App\Application\Exceptions\NotFoundException;
use App\Application\Mappers\UserMapper;
use App\Application\Resources\UserResource;
use App\Domain\Events\EventDispatcher;
use App\Domain\Factories\AvatarFactory;
use App\Domain\Repositories\UserRepository;
use App\Domain\ValueObjects\UserId;

final readonly class EditUserCommandHandler
{
    public function __construct(
        private UserRepository $userRepository,
        private EventDispatcher $eventDispatcher,
        private AvatarFactory $avatarFactory,
    ) {
    }

    public function handle(EditUserCommand $command): UserResource
    {
        $user = $this->userRepository->findById(new UserId($command->id));
        if (null === $user) {
            throw new NotFoundException('User not found');
        }

        if (null !== $command->username) {
            $user->changeUsername($command->username);
        }

        if (null !== $command->bio) {
            $user->changeBio($command->bio);
        }

        if (null !== $command->avatarPath) {
            $newAvatar = $this->avatarFactory->create($command->avatarPath);
            $user->changeAvatar($newAvatar);
        }

        $this->userRepository->save($user);

        $events = $user->pullPendingEvents();
        foreach ($events as $event) {
            $this->eventDispatcher->dispatch($event);
        }

        return UserMapper::toResource($user);
    }
}
