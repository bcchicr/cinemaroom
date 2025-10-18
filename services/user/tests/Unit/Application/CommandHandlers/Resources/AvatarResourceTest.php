<?php

declare(strict_types=1);

namespace App\Tests\Unit\Application\CommandHandlers\Resources;

use App\Application\Exceptions\InvalidArgumentException;
use App\Application\Resources\AvatarResource;
use App\Domain\ValueObjects\Avatar;
use PHPUnit\Framework\TestCase;

final class AvatarResourceTest extends TestCase
{
    public function testInstantiatesFromString(): void
    {
        $avatar = Avatar::fromString('/path/to/avatar.jpg');

        $resource = new AvatarResource($avatar);

        $this->assertEquals($avatar->path(), $resource->path());
    }

    public function testNotInstantiatesNull(): void
    {
        $this->expectException(InvalidArgumentException::class);
        $avatar = Avatar::null();

        $resource = new AvatarResource($avatar);
    }
}
