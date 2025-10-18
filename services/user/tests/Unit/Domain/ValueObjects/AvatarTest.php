<?php

declare(strict_types=1);

namespace App\Tests\Unit\Domain\ValueObjects;

use App\Domain\Exceptions\FailedInvariantException;
use App\Domain\ValueObjects\Avatar;
use PHPUnit\Framework\Attributes\DataProvider;
use PHPUnit\Framework\TestCase;

final class AvatarTest extends TestCase
{
    public static function empty_string_provider(): array
    {
        return [
            [''],
            [' '],
        ];
    }

    public static function valid_path_provider(): array
    {
        return [
            ['/path/to/avatar.jpg'],
            [null],
        ];
    }

    public function testCreatesWithNonEmptyString(): void
    {
        $path = '/path/to/avatar.jpg';
        $avatar = Avatar::fromString('/path/to/avatar.jpg');
        $this->assertFalse($avatar->isNull());
        $this->assertSame($path, $avatar->path());
    }

    #[DataProvider('empty_string_provider')]
    public function testRejectsEmptyString(string $path): void
    {
        $this->expectException(FailedInvariantException::class);
        Avatar::fromString($path);
    }

    public function testRejectsLongString(): void
    {
        $this->expectException(FailedInvariantException::class);
        Avatar::fromString(str_repeat('a', 256));
    }

    public function testCreateNullPath(): void
    {
        $avatar = Avatar::null();
        $this->assertTrue($avatar->isNull());
    }

    public function testCannotGetNullPath(): void
    {
        $this->expectException(FailedInvariantException::class);
        $avatar = Avatar::null();
        $avatar->path();
    }

    public function testEqualSamePaths(): void
    {
        $path = '/path/to/avatar.jpg';
        $avatar1 = Avatar::fromString($path);
        $avatar2 = Avatar::fromString($path);

        $this->assertEquals($avatar1->hash(), $avatar2->hash());
        $this->assertTrue($avatar1->equals($avatar2));
    }

    public function testNotEqualDifferentPaths(): void
    {
        $path1 = '/path/to/avatar1.jpg';
        $avatar1 = Avatar::fromString($path1);

        $path2 = '/path/to/avatar2.jpg';
        $avatar2 = Avatar::fromString($path2);

        $this->assertFalse($avatar1->equals($avatar2));
    }

    public function testEqualNullAndNull(): void
    {
        $avatar1 = Avatar::null();
        $avatar2 = Avatar::null();

        $this->assertTrue($avatar1->equals($avatar2));
    }

    public function testNotEqualNullAndNotNull(): void
    {
        $avatar1 = Avatar::null();
        $avatar2 = Avatar::fromString('/path/to/avatar.jpg');

        $this->assertFalse($avatar1->equals($avatar2));
    }

    public function testCorrectPathCastsToString(): void
    {
        $path = '/path/to/avatar.jpg';

        $avatar = Avatar::fromString($path);

        $this->assertEquals($path, $avatar->toString());
        $this->assertEquals($path, (string) $avatar);
    }

    public function testNullPathNotCastsToString(): void
    {
        $this->expectException(FailedInvariantException::class);
        $avatar = Avatar::null();
        $avatar->toString();
    }

    /**
     * @throws \JsonException
     */
    #[DataProvider('valid_path_provider')]
    public function testSerializesDeserializes(?string $path): void
    {
        $avatar = Avatar::create($path);

        $serialized = $avatar->toArray();
        $jsoned = json_encode($avatar, JSON_THROW_ON_ERROR);

        $this->assertEquals($path, $serialized['path']);
        $this->assertEquals(
            json_encode($serialized, JSON_THROW_ON_ERROR),
            $jsoned,
        );

        $unserialized = Avatar::create($serialized['path']);

        $jsonDecoded = json_decode($jsoned, true, 512, JSON_THROW_ON_ERROR);
        $unjsoned = Avatar::create($jsonDecoded['path']);

        $this->assertTrue($avatar->equals($unserialized));
        $this->assertTrue($avatar->equals($unjsoned));
    }
}
