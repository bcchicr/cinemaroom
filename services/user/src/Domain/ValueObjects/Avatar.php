<?php

declare(strict_types=1);

namespace App\Domain\ValueObjects;

use App\Domain\Abstracts\ValueObject;
use Doctrine\ORM\Mapping as ORM;

#[ORM\Embeddable()]
final class Avatar extends ValueObject
{
    #[ORM\Column(type: 'string', length: 255, nullable: true)]
    private ?string $path;

    public function __construct(
        string $path,
    ) {
        $this->path = $path;
    }

    public function __toString()
    {
        return $this->path ?? '';
    }

    public function equals(ValueObject $other): bool
    {
        return $this->path === $other->path;
    }

    public function path(): ?string
    {
        return $this->path;
    }

    public function toString(): string
    {
        return $this->path ?? '';
    }

    public function hash(): string
    {
        return md5($this->path ?? '');
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
