<?php

declare(strict_types=1);

namespace App\Domain\ValueObjects;

use App\Domain\Abstracts\ValueObject;
use App\Domain\Exceptions\FailedInvariantException;
use App\Domain\Exceptions\InvalidArgumentException;
use Doctrine\ORM\Mapping as ORM;

#[ORM\Embeddable()]
final class Avatar extends ValueObject
{
    #[ORM\Column(type: 'string', length: 255, nullable: true)]
    private ?string $path;

    private function __construct(
        ?string $path,
    ) {
        $this->setPath($path);
    }

    private function setPath(?string $path): void
    {
        if (null === $path) {
            $this->path = null;

            return;
        }

        $path = trim($path);
        if (empty($path)) {
            throw new InvalidArgumentException('Path cannot be empty');
        }

        if (mb_strlen($path) > 255) {
            throw new InvalidArgumentException('Path cannot be longer than 255 characters');
        }

        $this->path = $path;
    }

    public static function create(?string $path): self
    {
        if (null === $path) {
            return self::null();
        }

        return self::fromString($path);
    }

    public static function null(): self
    {
        return new self(null);
    }

    public static function fromString(string $path): self
    {
        return new self($path);
    }

    public function toString(): string
    {
        if ($this->isNull()) {
            throw new FailedInvariantException('Path cannot be null');
        }

        return $this->path;
    }

    public function isNull(): bool
    {
        return null === $this->path;
    }

    public function equals(ValueObject $other): bool
    {
        return $other instanceof self
            && $this->path === $other->path;
    }

    public function path(): string
    {
        if ($this->isNull()) {
            throw new FailedInvariantException('Path cannot be null');
        }

        return $this->path;
    }

    public function toArray(): array
    {
        return [
            'path' => $this->path,
        ];
    }
}
