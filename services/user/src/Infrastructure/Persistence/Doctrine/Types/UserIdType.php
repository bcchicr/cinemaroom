<?php

declare(strict_types=1);

namespace App\Infrastructure\Persistence\Doctrine\Types;

use App\Domain\ValueObjects\UserId;
use Doctrine\DBAL\Platforms\AbstractPlatform;
use Doctrine\DBAL\Types\StringType;
use Ramsey\Uuid\Uuid;

final class UserIdType extends StringType
{
    public const string NAME = 'user_id';

    public function getName(): string
    {
        return self::NAME;
    }

    public function convertToPHPValue($value, $platform): ?UserId
    {
        return null === $value ? null : new UserId(Uuid::fromString($value));
    }

    public function convertToDatabaseValue($value, $platform): ?string
    {
        return $value instanceof UserId ? $value->toString() : null;
    }

    public function requiresSQLCommentHint(AbstractPlatform $platform): bool
    {
        return true;
    }
}
