<?php

declare(strict_types=1);

namespace App\Application\QueryHandlers\UserQueryHandler\GetUserQueryHandler;

use Ramsey\Uuid\UuidInterface;

final readonly class GetUserQuery
{
    public function __construct(
        public UuidInterface $id,
    ) {
    }
}
