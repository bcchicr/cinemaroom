<?php

declare(strict_types=1);

namespace App\Application\CommandHandlers\UserCommandHandler\EditUserCommandHandler;

use Ramsey\Uuid\UuidInterface;

final readonly class EditUserCommand
{
    public function __construct(
        public UuidInterface $id,
        public ?string $username,
        public ?string $bio,
        public ?string $avatarPath,
    ) {
    }
}
