<?php

declare(strict_types=1);

namespace App\Application\QueryHandlers\UserQueryHandler\GetUserQueryHandler;

use App\Application\Resources\UserResource;

interface GetUserQueryHandler
{
    public function handle(GetUserQuery $query): UserResource;
}
