<?php

declare(strict_types=1);

namespace App\Application\QueryHandlers\UserQueryHandler\GetUserQueryHandler;

use App\Application\Exceptions\NotFoundException;
use App\Application\Resources\UserResource;

interface GetUserQueryHandler
{
    /**
     * @throws NotFoundException
     */
    public function handle(GetUserQuery $query): UserResource;
}
