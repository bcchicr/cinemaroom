<?php

declare(strict_types=1);

namespace App\Application\CommandHandlers\UserCommandHandler\RegisterUserCommandHandler;

use Ramsey\Uuid\UuidInterface;

final readonly class RegisterUserCommand
{
    public function __construct(
        public UuidInterface $userId,
        public string $username,
    ) {}
}
