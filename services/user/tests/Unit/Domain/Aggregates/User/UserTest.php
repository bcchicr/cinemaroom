<?php

declare(strict_types=1);

namespace App\Tests\Unit\Domain\Aggregates\User;

use App\Domain\Aggregates\User\User;
use App\Domain\Events\UserAvatarChanged;
use App\Domain\Events\UserBioChanged;
use App\Domain\Events\UserRegistered;
use App\Domain\Events\UserUsernameChanged;
use App\Domain\Exceptions\InvalidArgumentException;
use App\Domain\ValueObjects\Avatar;
use App\Domain\ValueObjects\UserId;
use PHPUnit\Framework\Attributes\DataProvider;
use PHPUnit\Framework\TestCase;
use Ramsey\Uuid\Uuid;

final class UserTest extends TestCase
{
    private User $user;

    public function testRegister(): void
    {
        $userId = new UserId(Uuid::uuid7());
        $username = 'test_user';

        $user = User::register($userId, $username);

        $this->assertTrue($user->id()->equals($userId));
        $this->assertEquals($username, $user->username());
        $this->assertEquals('', $user->bio());
        $this->assertTrue($user->avatar()->isNull());

        $events = $user->pullPendingEvents();
        $this->assertCount(1, $events);

        $event = $events[0];
        $this->assertInstanceOf(UserRegistered::class, $event);
    }

    public function testCannotRegisterWithLongUsername(): void
    {
        $this->expectException(InvalidArgumentException::class);
        User::register(new UserId(Uuid::uuid7()), str_repeat('a', 256));
    }

    #[DataProvider('empty_username_provider')]
    public function testCannotRegisterWithEmptyUsername(string $username): void
    {
        $this->expectException(InvalidArgumentException::class);
        User::register(new UserId(Uuid::uuid7()), $username);
    }

    public static function empty_username_provider(): array
    {
        return [
            [''],
            [' '],
            ["\t"],
            ["\n"],
            ["\r"],
        ];
    }

    public function testChangeUsername(): void
    {
        $newUsername = 'new_username';
        $this->user->changeUsername($newUsername);

        $this->assertEquals($newUsername, $this->user->username());

        $events = $this->user->pullPendingEvents();
        $this->assertCount(1, $events);

        $event = $events[0];
        $this->assertInstanceOf(UserUsernameChanged::class, $event);
    }

    public function testCannotChangeLongUsername(): void
    {
        $this->expectException(InvalidArgumentException::class);
        $this->user->changeUsername(str_repeat('a', 256));
    }

    public function testChangeBio(): void
    {
        $newBio = 'new_bio';
        $this->user->changeBio($newBio);

        $this->assertEquals($newBio, $this->user->bio());

        $events = $this->user->pullPendingEvents();
        $this->assertCount(1, $events);

        $event = $events[0];
        $this->assertInstanceOf(UserBioChanged::class, $event);
    }

    public function testChangeAvatar(): void
    {
        $newAvatar = Avatar::fromString('/path/to/avatar.jpg');
        $this->user->changeAvatar($newAvatar);

        $this->assertTrue($newAvatar->equals($this->user->avatar()));

        $events = $this->user->pullPendingEvents();
        $this->assertCount(1, $events);

        $event = $events[0];
        $this->assertInstanceOf(UserAvatarChanged::class, $event);
    }

    public function testCannotChangeAvatarToNull(): void
    {
        $this->expectException(InvalidArgumentException::class);
        $newAvatar = Avatar::null();
        $this->user->changeAvatar($newAvatar);
    }

    protected function setUp(): void
    {
        parent::setUp();

        $userId = new UserId(Uuid::uuid7());
        $username = 'test_user';
        $this->user = User::hydrate(
            id: $userId,
            username: $username,
            bio: '',
            avatar: Avatar::null(),
        );
    }
}
