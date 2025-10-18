<?php

declare(strict_types=1);

namespace App\Tests\Unit\Domain\ValueObjects;

use App\Domain\ValueObjects\UserId;
use PHPUnit\Framework\TestCase;
use Ramsey\Uuid\Uuid;

final class UserIdTest extends TestCase
{
    public function testCreatesWithValidUuid(): void
    {
        $uuid = Uuid::uuid7();
        $userId = new UserId($uuid);
        $this->assertTrue($userId->value()->equals($uuid));
    }

    public function testEqualWithSameUuid(): void
    {
        $uuid = Uuid::uuid7();
        $userId1 = new UserId($uuid);
        $userId2 = new UserId($uuid);

        $this->assertEquals($userId1->hash(), $userId2->hash());
        $this->assertTrue($userId1->equals($userId2));
    }

    public function testNotEqualWithDifferentUuids(): void
    {
        $uuid1 = Uuid::uuid7();
        $userId1 = new UserId($uuid1);
        $uuid2 = Uuid::uuid7();
        $userId2 = new UserId($uuid2);

        $this->assertFalse($userId1->equals($userId2));
    }

    public function testCastsToString(): void
    {
        $uuid = Uuid::uuid7();
        $userId = new UserId($uuid);

        $this->assertSame($uuid->toString(), $userId->toString());
        $this->assertSame($uuid->toString(), (string) $userId);
    }

    /**
     * @throws \JsonException
     */
    public function testSerializesDeserializes()
    {
        $uuid = Uuid::uuid7();

        $userId = new UserId($uuid);

        $serialized = $userId->toArray();
        $jsoned = json_encode($userId, JSON_THROW_ON_ERROR);

        $this->assertEquals($uuid->toString(), $serialized['value']);
        $this->assertEquals(
            json_encode($serialized, JSON_THROW_ON_ERROR),
            $jsoned,
        );

        $unserialized = new UserId(Uuid::fromString($serialized['value']));

        $jsonDecoded = json_decode($jsoned, true, 512, JSON_THROW_ON_ERROR);
        $unjsoned = new UserId(Uuid::fromString($jsonDecoded['value']));

        $this->assertTrue($userId->equals($unserialized));
        $this->assertTrue($userId->equals($unjsoned));
    }
}
