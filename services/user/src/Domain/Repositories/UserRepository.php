<?php

declare(strict_types=1);

namespace App\Domain\Repositories;

use App\Domain\Aggregates\User\User;
use App\Domain\ValueObjects\UserId;

interface UserRepository
{
    public function findById(UserId $id): ?User;

    public function save(User $user): void;
}
