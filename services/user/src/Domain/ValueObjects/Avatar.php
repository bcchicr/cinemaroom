<?php

declare(strict_types=1);

namespace App\Domain\ValueObjects;

use App\Domain\Abstracts\ValueObject;
use App\Domain\Exceptions\FailedInvariantException;
use Doctrine\ORM\Mapping as ORM;

#[ORM\Embeddable()]
final class Avatar extends ValueObject
{
    #[ORM\Column(type: 'string', length: 255, nullable: true)]
    private ?string $path;

    public function __construct(
        ?string $path,
    ) {
        $this->setPath($path);
    }

    private function setPath(?string $path): void
    {
        if (null !== $path && empty(trim($path))) {
            throw new FailedInvariantException('Path cannot be empty');
        }
        $this->path = $path;
    }

    public function __toString(): string
    {
        return $this->toString();
    }

    public function toString(): string
    {
        return $this->path ?? '';
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

    public function isNull(): bool
    {
        return null === $this->path;
    }

    public function hash(): string
    {
        return md5($this->toString());
    }

    public function jsonSerialize(): array
    {
        return $this->toArray();
    }

    public function toArray(): array
    {
        return [
            'path' => $this->path,
        ];
    }
}
