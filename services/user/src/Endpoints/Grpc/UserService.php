<?php

declare(strict_types=1);

namespace App\Endpoints\Grpc;

use App\Application\CommandHandlers\UserCommandHandler\EditUserCommandHandler\EditUserCommand;
use App\Application\CommandHandlers\UserCommandHandler\EditUserCommandHandler\EditUserCommandHandler;
use App\Application\CommandHandlers\UserCommandHandler\RegisterUserCommandHandler\RegisterUserCommand;
use App\Application\CommandHandlers\UserCommandHandler\RegisterUserCommandHandler\RegisterUserCommandHandler;
use App\Application\QueryHandlers\UserQueryHandler\GetUserQueryHandler\GetUserQuery;
use App\Application\QueryHandlers\UserQueryHandler\GetUserQueryHandler\GetUserQueryHandler;
use App\Endpoints\Grpc\Mappers\UserMapper;
use GRPC\Services\Users\v1\EditRequest;
use GRPC\Services\Users\v1\EditResponse;
use GRPC\Services\Users\v1\GetRequest;
use GRPC\Services\Users\v1\GetResponse;
use GRPC\Services\Users\v1\RegisterRequest;
use GRPC\Services\Users\v1\RegisterResponse;
use GRPC\Services\Users\v1\UserServiceInterface;
use Ramsey\Uuid\Uuid;
use Spiral\RoadRunner\GRPC;

final readonly class UserService implements UserServiceInterface
{
    public function __construct(
        private RegisterUserCommandHandler $registerUserCommandHandler,
        private EditUserCommandHandler $editUserCommandHandler,
        private GetUserQueryHandler $getUserQueryHandler,
        private UserMapper $userMapper,
    ) {
    }

    #[\Override]
    public function Register(GRPC\ContextInterface $ctx, RegisterRequest $in): RegisterResponse
    {
        $command = new RegisterUserCommand(
            userId: Uuid::fromString($in->getId()?->getValue() ?? ''),
            username: $in->getUsername(),
        );
        $user = $this->registerUserCommandHandler->handle($command);

        $response = new RegisterResponse();
        $response->setUser($this->userMapper->toGrpc($user));

        return $response;
    }

    #[\Override]
    public function Get(GRPC\ContextInterface $ctx, GetRequest $in): GetResponse
    {
        $id = $in->getId()?->getValue() ?? '';
        $query = new GetUserQuery(Uuid::fromString($id));
        $user = $this->getUserQueryHandler->handle($query);

        $response = new GetResponse();
        $response->setUser($this->userMapper->toGrpc($user));

        return $response;
    }

    #[\Override]
    public function Edit(GRPC\ContextInterface $ctx, EditRequest $in): EditResponse
    {
        $username = null;
        if ($in->hasUsername()) {
            $username = $in->getUsername();
        }

        $bio = null;
        if ($in->hasBio()) {
            $bio = $in->getBio();
        }

        $avatarPath = null;
        if ($in->hasAvatar()) {
            $avatarPath = $in->getAvatar()?->getPath();
        }

        $command = new EditUserCommand(
            id: Uuid::fromString($in->getId()?->getValue() ?? ''),
            username: $username,
            bio: $bio,
            avatarPath: $avatarPath,
        );
        $user = $this->editUserCommandHandler->handle($command);

        $response = new EditResponse();
        $response->setUser($this->userMapper->toGrpc($user));

        return $response;
    }
}
