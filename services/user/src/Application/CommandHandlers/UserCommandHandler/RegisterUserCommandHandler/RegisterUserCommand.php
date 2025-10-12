<?php

declare(strict_types=1);

namespace App\Application\CommandHandlers\UserCommandHandler\RegisterUserCommandHandler;

final readonly class RegisterUserCommand
{
    public function __construct(
        public string $username,
    ) {
    }
}
