<?php

declare(strict_types=1);

namespace App\Infrastructure\Persistence\QueryHandlers\UserQueryHandler\GetUserQueryHandler;

use App\Application\Mappers\UserMapper;
use App\Application\QueryHandlers\UserQueryHandler\GetUserQueryHandler\GetUserQuery;
use App\Application\QueryHandlers\UserQueryHandler\GetUserQueryHandler\GetUserQueryHandler;
use App\Application\Resources\UserResource;
use App\Domain\Aggregates\User\User;
use App\Domain\ValueObjects\Avatar;
use App\Domain\ValueObjects\UserId;
use Doctrine\Bundle\DoctrineBundle\Repository\ServiceEntityRepository;
use Doctrine\ORM\EntityManagerInterface;
use Doctrine\Persistence\ManagerRegistry;

final class DoctrineGetUserQueryHandler implements GetUserQueryHandler
{
    public function __construct(
        private EntityManagerInterface $em,
    ) {}

    public function handle(GetUserQuery $query): UserResource
    {
        $qb = $this->em->createQueryBuilder('u');

        $qb->select(
            'u.id AS user_id',
            'u.username',
            'u.bio',
            'u.avatar.path AS avatar_path',
        )
            ->from(User::class, 'u')
            ->where('u.id = :id')
            ->setParameter('id', new UserId($query->id));

        $result = $qb->getQuery()->getOneOrNullResult(\Doctrine\ORM\Query::HYDRATE_ARRAY);

        if ($result === null) {
            throw new \RuntimeException('User not found');
        }

        $avatar = null;
        if ($result['avatar_path'] !== null) {
            $avatar = new Avatar(
                $result['avatar_path'],
            );
        }

        $user = User::hydrate(
            $result['user_id'],
            $result['username'],
            $result['bio'],
            $avatar,
        );

        return UserMapper::toResource($user);
    }
}
