<?php

declare(strict_types=1);

namespace App\Infrastructure\Persistence\QueryHandlers\UserQueryHandler\GetUserQueryHandler;

use App\Application\Exceptions\NotFoundException;
use App\Application\QueryHandlers\UserQueryHandler\GetUserQueryHandler\GetUserQuery;
use App\Application\QueryHandlers\UserQueryHandler\GetUserQueryHandler\GetUserQueryHandler;
use App\Application\Resources\UserResource;
use App\Domain\Aggregates\User\User;
use App\Domain\ValueObjects\UserId;
use Doctrine\Bundle\DoctrineBundle\Repository\ServiceEntityRepository;
use Doctrine\Persistence\ManagerRegistry;

final class DoctrineGetUserQueryHandler extends ServiceEntityRepository implements GetUserQueryHandler
{
    public function __construct(ManagerRegistry $registry)
    {
        parent::__construct($registry, User::class);
    }

    public function handle(GetUserQuery $query): UserResource
    {
        $user = $this->find(new UserId($query->id));

        if (null === $user) {
            throw new NotFoundException('User not found.');
        }

        return new UserResource($user);
    }
}
