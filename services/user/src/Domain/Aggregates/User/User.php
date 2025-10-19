<?php

declare(strict_types=1);

namespace App\Domain\Aggregates\User;

use App\Domain\Abstracts\AggregateRoot;
use App\Domain\Events\UserAvatarChanged;
use App\Domain\Events\UserBioChanged;
use App\Domain\Events\UserRegistered;
use App\Domain\Events\UserUsernameChanged;
use App\Domain\Exceptions\FailedInvariantException;
use App\Domain\Exceptions\InvalidArgumentException;
use App\Domain\ValueObjects\Avatar;
use App\Domain\ValueObjects\UserId;
use Doctrine\ORM\Mapping as ORM;
use Doctrine\ORM\Mapping\Entity;

#[Entity]
#[ORM\Table(name: 'users')]
class User extends AggregateRoot
{
    #[ORM\Id()]
    #[ORM\Column(type: 'user_id', length: 36, unique: true)]
    protected UserId $id;
    #[ORM\Column(type: 'string', length: 255)]
    private string $username;
    #[ORM\Column(type: 'text')]
    private string $bio;
    #[ORM\Embedded(class: Avatar::class, columnPrefix: 'avatar_')]
    private Avatar $avatar;

    private function __construct(
        UserId $id,
        string $username,
        Avatar $avatar,
        string $bio,
    ) {
        $this->setId($id);
        $this->setUsername($username);
        $this->setBio($bio);
        $this->setAvatar($avatar);
    }

    public function setId(UserId $id): void
    {
        $this->id = $id;
    }

    public function setUsername(string $username): void
    {
        $username = trim($username);
        if (empty($username)) {
            throw new InvalidArgumentException('Username cannot be empty');
        }

        if (mb_strlen($username) > 255) {
            throw new InvalidArgumentException('Username cannot be longer than 255 characters');
        }
        $this->username = $username;
    }

    public function setBio(string $bio): void
    {
        $this->bio = $bio;
    }

    public function setAvatar(Avatar $avatar): void
    {
        $this->avatar = $avatar;
    }

    public static function hydrate(
        UserId $id,
        string $username,
        string $bio,
        Avatar $avatar,
    ): self {
        return new self(
            id: $id,
            username: $username,
            avatar: $avatar,
            bio: $bio,
        );
    }

    public static function register(
        UserId $id,
        string $username,
    ): self {
        $user = new self(
            id: $id,
            username: $username,
            avatar: Avatar::null(),
            bio: '',
        );
        $user->record(new UserRegistered(
            id: $id,
            username: $username,
        ));

        return $user;
    }

    #[\Override]
    public function id(): UserId
    {
        return $this->id;
    }

    public function username(): string
    {
        return $this->username;
    }

    public function bio(): string
    {
        return $this->bio;
    }

    public function avatar(): Avatar
    {
        return $this->avatar;
    }

    public function changeUsername(string $newUsername): void
    {
        $oldUsername = $this->username;
        $this->setUsername($newUsername);
        $this->record(new UserUsernameChanged(
            id: $this->id,
            oldUsername: $oldUsername,
            newUsername: $newUsername,
        ));
    }

    public function changeBio(string $newBio): void
    {
        $oldBio = $this->bio;
        $this->setBio($newBio);
        $this->record(new UserBioChanged(
            id: $this->id,
            oldBio: $oldBio,
            newBio: $newBio,
        ));
    }

    public function changeAvatar(Avatar $newAvatar): void
    {
        $oldAvatar = $this->avatar;

        if ($newAvatar->isNull()) {
            throw new InvalidArgumentException('Cannot change avatar to null');
        }

        $this->setAvatar($newAvatar);
        $this->record(new UserAvatarChanged(
            id: $this->id,
            oldAvatar: $oldAvatar,
            newAvatar: $newAvatar,
        ));
    }
}
